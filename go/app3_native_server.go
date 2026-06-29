package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log" // Añadido para ver errores en tu terminal
	"net/http"
	"regexp"
	"strconv"
)

const WSDL_URL = "https://www.dataaccess.com/webservicesserver/NumberConversion.wso"
const TRANSLATE_URL = "https://libretranslate.de/translate"

func makeSoapEnvelope(n string) string {
	return `<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>` + n + `</ubiNum></NumberToWords></soap:Body></soap:Envelope>`
}

func callNumberToWords(n string) (string, error) {
	// Añadimos timeout para que no se quede bloqueado
	client := &http.Client{}
	req, _ := http.NewRequest("POST", WSDL_URL, bytes.NewBufferString(makeSoapEnvelope(n)))
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords")
	
	resp, err := client.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()
	
	b, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`<[^>]*NumberToWordsResult>([^<]+)</[^>]*NumberToWordsResult>`)
	m := re.FindStringSubmatch(string(b))
	if len(m) >= 2 { return m[1], nil }
	return "Error SOAP", fmt.Errorf("soap error")
}

func translateToSpanish(text string) (string, error) {
	payload := map[string]string{"q": text, "source": "en", "target": "es", "format": "text"}
	p, _ := json.Marshal(payload)
	
	// Añadimos User-Agent para evitar bloqueos
	req, _ := http.NewRequest("POST", TRANSLATE_URL, bytes.NewBuffer(p))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0") 
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return fallbackTranslate(text), nil }
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	var js map[string]interface{}
	if err := json.Unmarshal(body, &js); err == nil {
		if v, ok := js["translatedText"]; ok { return fmt.Sprint(v), nil }
	}
	return fallbackTranslate(text), nil
}

func fallbackTranslate(text string) string {
	m := map[string]string{"zero": "cero", "one": "uno", "two": "dos", "three": "tres", "four": "cuatro", "five": "cinco", "six": "seis", "seven": "siete", "eight": "ocho", "nine": "nueve", "ten": "diez"}
	words := regexp.MustCompile(`\w+`).FindAllString(text, -1)
	for i, w := range words {
		if v, ok := m[w]; ok { words[i] = v }
	}
	return fmt.Sprint(words) // Simplificado
}

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("n")
	if q == "" { q = "0" }
	
	log.Printf("Petición recibida: n=%s", q) // Esto te dirá en consola si llega la petición
	
	eng, err := callNumberToWords(q)
	if err != nil {
		http.Error(w, "Error SOAP: "+err.Error(), 500)
		return
	}
	
	esp, _ := translateToSpanish(eng)
	
	// Si la traducción parece fallida (contiene inglés), usamos nativo
	if regexp.MustCompile(`\b(hundred|thousand|one|two)\b`).MatchString(esp) {
		n, _ := strconv.Atoi(q)
		fmt.Fprint(w, belowThousand(n)) // Usamos el nativo
	} else {
		fmt.Fprint(w, esp)
	}
}

// ... (Aquí mantén tus funciones belowThousand y numberToSpanish) ...
func belowThousand(num int) string {
    units := []string{"cero", "uno", "dos", "tres", "cuatro", "cinco", "seis", "siete", "ocho", "nueve", "diez", "once", "doce", "trece", "catorce", "quince", "dieciséis", "diecisiete", "dieciocho", "diecinueve"}
    tens := []string{"", "", "veinte", "treinta", "cuarenta", "cincuenta", "sesenta", "setenta", "ochenta", "noventa"}
    if num < 20 { return units[num] }
    if num < 100 { 
        t := num / 10; u := num % 10
        if t == 2 && u > 0 { return "veinti" + units[u] }
        if u > 0 { return tens[t] + " y " + units[u] }
        return tens[t]
    }
    hundreds := []string{"", "ciento", "doscientos", "trescientos", "cuatrocientos", "quinientos", "seiscientos", "setecientos", "ochocientos", "novecientos"}
    if num == 100 { return "cien" }
    h := num / 100; rem := num % 100
    if rem > 0 { return hundreds[h] + " " + belowThousand(rem) }
    return hundreds[h]
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor activo en http://localhost:8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}