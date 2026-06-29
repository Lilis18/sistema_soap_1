#include <iostream>
#include <string>
#include <winsock2.h>
#include <wininet.h>
#include <regex>
#include <map>

#pragma comment(lib, "ws2_32.lib")
#pragma comment(lib, "wininet.lib")

// Función Nativa de Red (isPost=true para POST, false para GET)
std::string requestHttp(bool isPost, std::string host, std::string path, std::string body = "", std::string headers = "") {
    HINTERNET hOpen = InternetOpen("App", INTERNET_OPEN_TYPE_DIRECT, NULL, NULL, 0);
    HINTERNET hConnect = InternetConnect(hOpen, host.c_str(), INTERNET_DEFAULT_HTTPS_PORT, NULL, NULL, INTERNET_SERVICE_HTTP, 0, 0);
    
    // Si es MyMemory (traducción), el puerto suele ser 80 (HTTP) o 443 (HTTPS)
    // Forzamos HTTPS para ambos, es más seguro
    HINTERNET hRequest = HttpOpenRequest(hConnect, isPost ? "POST" : "GET", path.c_str(), NULL, NULL, NULL, INTERNET_FLAG_SECURE, 0);

    // Enviar headers y cuerpo
    HttpSendRequest(hRequest, headers.c_str(), -1, (LPVOID)body.c_str(), body.length());

    char buf[4096]; DWORD read; std::string resp;
    while (InternetReadFile(hRequest, buf, sizeof(buf)-1, &read) && read > 0) { buf[read] = 0; resp += buf; }
    
    InternetCloseHandle(hRequest); InternetCloseHandle(hConnect); InternetCloseHandle(hOpen);
    return resp;
}

std::string callSoap(std::string n) {
    std::string env = "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><NumberToWords xmlns=\"http://www.dataaccess.com/webservicesserver/\"><ubiNum>"+n+"</ubiNum></NumberToWords></soap:Body></soap:Envelope>";
    std::string hdrs = "Content-Type: text/xml; charset=utf-8\r\nSOAPAction: http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords";
    
    // Aquí pasamos 'true' porque es POST
    std::string resp = requestHttp(true, "www.dataaccess.com", "/webservicesserver/NumberConversion.wso", env, hdrs);
    
    std::regex re("<[^>]+:NumberToWordsResult>(.*?)</[^>]+:NumberToWordsResult>");
    std::smatch m;
    if (std::regex_search(resp, m, re)) return m[1].str();
    
    return "Error";
}

std::string translateNative(std::string eng) {
    if (eng == "Error") return "Error";
    
    std::string path = "/get?q=" + eng + "&langpair=en|es";
    // Aquí pasamos 'false' porque es GET
    std::string resp = requestHttp(false, "api.mymemory.translated.net", path);
    
    std::regex re("\"translatedText\":\"([^\"]+)\"");
    std::smatch m;
    if (std::regex_search(resp, m, re)) return m[1].str();
    
    return eng; 
}

int main() {
    WSADATA wsa; WSAStartup(MAKEWORD(2,2), &wsa);
    SOCKET s = socket(AF_INET, SOCK_STREAM, 0);
    sockaddr_in addr = {AF_INET, htons(8000), INADDR_ANY};
    bind(s, (sockaddr*)&addr, sizeof(addr)); listen(s, 10);
    
    std::cout << "Servidor corriendo en http://localhost:8000" << std::endl;

    while(true) {
        SOCKET c = accept(s, NULL, NULL);
        char buf[2048]; recv(c, buf, 2048, 0);
        std::string req(buf);
        
        std::smatch m; std::regex re("n=(\\d+)");
        std::string n = std::regex_search(req, m, re) ? m[1].str() : "0";
        
        std::string eng = callSoap(n);
        std::string esp = translateNative(eng);
        
        // Respuesta HTTP robusta
        std::string res = "HTTP/1.1 200 OK\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n" + esp;
        send(c, res.c_str(), res.size(), 0);
        closesocket(c);
    }
}