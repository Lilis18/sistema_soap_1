#include <iostream>
#include <string>
#include <winsock2.h>
#include <wininet.h>
#pragma comment(lib, "ws2_32.lib")
#pragma comment(lib, "wininet.lib")

std::string callSoap(std::string n) {
    HINTERNET hOpen = InternetOpen("App1", INTERNET_OPEN_TYPE_DIRECT, NULL, NULL, 0);
    HINTERNET hConnect = InternetConnect(hOpen, "www.dataaccess.com", INTERNET_DEFAULT_HTTPS_PORT, NULL, NULL, INTERNET_SERVICE_HTTP, 0, 0);
    HINTERNET hRequest = HttpOpenRequest(hConnect, "POST", "/webservicesserver/NumberConversion.wso", NULL, NULL, NULL, INTERNET_FLAG_SECURE, 0);

    std::string env = "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><NumberToWords xmlns=\"http://www.dataaccess.com/webservicesserver/\"><ubiNum>"+n+"</ubiNum></NumberToWords></soap:Body></soap:Envelope>";
    const char* hdrs = "Content-Type: text/xml; charset=utf-8\r\nSOAPAction: http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords";
    HttpSendRequest(hRequest, hdrs, -1, (LPVOID)env.c_str(), env.length());

    char buf[4096]; DWORD read; std::string resp;
    while (InternetReadFile(hRequest, buf, sizeof(buf)-1, &read) && read > 0) { buf[read] = 0; resp += buf; }
    
    size_t s = resp.find("<NumberToWordsResult>") + 21;
    size_t e = resp.find("</NumberToWordsResult>");
    InternetCloseHandle(hRequest); InternetCloseHandle(hConnect); InternetCloseHandle(hOpen);
    return resp.substr(s, e - s);
}

int main() {
    WSADATA wsa; WSAStartup(MAKEWORD(2,2), &wsa);
    SOCKET s = socket(AF_INET, SOCK_STREAM, 0);
    sockaddr_in addr = {AF_INET, htons(8000), INADDR_ANY};
    bind(s, (sockaddr*)&addr, sizeof(addr)); listen(s, 10);
    while(true) {
        SOCKET c = accept(s, NULL, NULL);
        char buf[2048]; recv(c, buf, 2048, 0);
        std::string req(buf);
        size_t pos = req.find("n=");
        std::string n = (pos != std::string::npos) ? req.substr(pos + 2, req.find(" ", pos) - (pos + 2)) : "0";
        std::string res = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n" + callSoap(n);
        send(c, res.c_str(), res.size(), 0); closesocket(c);
    }
}