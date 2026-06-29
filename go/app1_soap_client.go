package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "regexp"
)

func handler(w http.ResponseWriter, r *http.Request) {
    n := r.URL.Query().Get("n")
    if n == "" { n = "0" }

    envelope := `<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>` + n + `</ubiNum></NumberToWords></soap:Body></soap:Envelope>`
    
    resp, err := http.Post("https://www.dataaccess.com/webservicesserver/NumberConversion.wso", "text/xml; charset=utf-8", bytes.NewBufferString(envelope))
    if err != nil { http.Error(w, err.Error(), 500); return }
    defer resp.Body.Close()
    
    b, _ := io.ReadAll(resp.Body)
    re := regexp.MustCompile(`<[^>]*NumberToWordsResult>([^<]+)</[^>]*NumberToWordsResult>`)
    m := re.FindStringSubmatch(string(b))
    
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    if len(m) > 1 { fmt.Fprint(w, m[1]) } else { fmt.Fprint(w, "Error") }
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Servidor V1 iniciado en http://localhost:8000")
    http.ListenAndServe(":8000", nil)
}v