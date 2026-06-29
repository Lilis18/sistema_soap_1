package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "regexp"
)

// Función de traducción robusta
func translate(text string) string {
    // Codificamos el texto para la URL
    encodedText := url.QueryEscape(text)
    apiURL := "https://api.mymemory.translated.net/get?q=" + encodedText + "&langpair=en|es"
    
    resp, err := http.Get(apiURL)
    if err != nil { return text }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    // Parseamos la respuesta JSON de MyMemory
    var j map[string]interface{}
    if err := json.Unmarshal(body, &j); err == nil {
        if responseData, ok := j["responseData"].(map[string]interface{}); ok {
            if t, ok := responseData["translatedText"].(string); ok {
                return t
            }
        }
    }
    return text // Si falla, regresa el original
}

func handler(w http.ResponseWriter, r *http.Request) {
    n := r.URL.Query().Get("n")
    if n == "" { n = "0" }
    
    // 1. Llamada SOAP
    env := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>` + n + `</ubiNum></NumberToWords></soap:Body></soap:Envelope>`
    resp, err := http.Post("https://www.dataaccess.com/webservicesserver/NumberConversion.wso", "text/xml; charset=utf-8", bytes.NewBufferString(env))
    if err != nil { fmt.Fprint(w, "Error SOAP"); return }
    defer resp.Body.Close()
    
    b, _ := io.ReadAll(resp.Body)
    m := regexp.MustCompile(`<[^>]*NumberToWordsResult>([^<]+)</[^>]*NumberToWordsResult>`).FindStringSubmatch(string(b))
    
    eng := ""
    if len(m) > 1 { eng = m[1] } else { fmt.Fprint(w, "Error SOAP"); return }
    
    // 2. Traducción
    esp := translate(eng)
    
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprint(w, esp)
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Servidor V2 (MyMemory) iniciado en http://localhost:8000")
    http.ListenAndServe(":8000", nil)
}