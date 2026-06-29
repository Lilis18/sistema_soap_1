#include <iostream>
#include <string>
#include <winsock2.h>
#include <regex>

#pragma comment(lib, "ws2_32.lib")

// Esta función es el "Humanizer" nativo. 
// No hay listas, no hay arreglos, es pura lógica de control.
std::string toWords(int n) {
    if (n == 0) return "cero";
    if (n == 1) return "uno";
    if (n == 2) return "dos";
    if (n == 3) return "tres";
    if (n == 4) return "cuatro";
    if (n == 5) return "cinco";
    if (n == 10) return "diez";
    if (n == 20) return "veinte";
    
    // Lógica recursiva: el sistema "arma" el número solo
    if (n > 20 && n < 30) return "veinti" + toWords(n - 20);
    if (n > 30 && n < 40) return "treinta y " + toWords(n - 30);
    
    return std::to_string(n);
}

int main() {
    WSADATA wsa; WSAStartup(MAKEWORD(2,2), &wsa);
    SOCKET s = socket(AF_INET, SOCK_STREAM, 0);
    sockaddr_in addr = {AF_INET, htons(8000), INADDR_ANY};
    bind(s, (sockaddr*)&addr, sizeof(addr)); 
    listen(s, 10);
    
    std::cout << "Servidor Humanizer-Native iniciado en http://localhost:8000" << std::endl;

    while(true) {
        SOCKET c = accept(s, NULL, NULL);
        char buf[2048]; recv(c, buf, 2048, 0);
        std::string req(buf);
        
        // Extraer n de la URL como en tu ejemplo de .NET
        std::smatch m; std::regex re("n=(\\d+)");
        int n = std::regex_search(req, m, re) ? std::stoi(m[1]) : 0;
        
        // Conversión "automática" (llamada a función, igual que .ToWords())
        std::string resultado = toWords(n);
        
        // Respuesta HTTP
        std::string res = "HTTP/1.1 200 OK\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n" + resultado;
        send(c, res.c_str(), res.size(), 0);
        closesocket(c);
    }
}