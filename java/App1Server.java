import com.sun.net.httpserver.*;
import java.io.*;
import java.net.*;
import java.util.regex.*;

public class App1Server {
    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress(8000), 0);
        server.createContext("/", exchange -> {
            String query = exchange.getRequestURI().getQuery();
            String n = (query != null && query.startsWith("n=")) ? query.substring(2) : "0";
            
            try {
                // Intentamos realizar la petición SOAP
                String response = callSoap(n);
                exchange.getResponseHeaders().set("Content-Type", "text/plain; charset=utf-8");
                exchange.sendResponseHeaders(200, response.getBytes("UTF-8").length);
                try (OutputStream os = exchange.getResponseBody()) {
                    os.write(response.getBytes("UTF-8"));
                }
            } catch (Exception e) {
                // Si hay un error, captúralo aquí para que el servidor siga funcionando
                String errorMsg = "Error en el servidor: " + e.getMessage();
                exchange.getResponseHeaders().set("Content-Type", "text/plain; charset=utf-8");
                exchange.sendResponseHeaders(500, errorMsg.getBytes("UTF-8").length);
                try (OutputStream os = exchange.getResponseBody()) {
                    os.write(errorMsg.getBytes("UTF-8"));
                }
            }
        });
        server.start();
        System.out.println("Servidor V1 iniciado en http://localhost:8000");
    }

    private static String callSoap(String n) throws Exception {
        String url = "https://www.dataaccess.com/webservicesserver/NumberConversion.wso";
        String envelope = "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><NumberToWords xmlns=\"http://www.dataaccess.com/webservicesserver/\"><ubiNum>"+n+"</ubiNum></NumberToWords></soap:Body></soap:Envelope>";
        HttpURLConnection con = (HttpURLConnection)new URL(url).openConnection();
        con.setDoOutput(true); con.setRequestMethod("POST");
        con.setRequestProperty("Content-Type", "text/xml; charset=utf-8");
        con.setRequestProperty("SOAPAction", "http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords");
        try(OutputStream os = con.getOutputStream()){ os.write(envelope.getBytes("UTF-8")); }
        BufferedReader br = new BufferedReader(new InputStreamReader(con.getInputStream()));
        StringBuilder sb = new StringBuilder(); String line;
        while((line=br.readLine())!=null) sb.append(line);
        Matcher m = Pattern.compile("<[^>]*NumberToWordsResult>([^<]+)</[^>]*NumberToWordsResult>").matcher(sb.toString());
        return m.find() ? m.group(1) : "Error en SOAP";
    }
}