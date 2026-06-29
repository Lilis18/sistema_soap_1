import com.sun.net.httpserver.*;
import java.io.*;
import java.net.*;
import java.util.regex.*;
import java.nio.charset.StandardCharsets;

public class App2Server {
    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress(8000), 0);
        server.createContext("/", exchange -> {
            try {
                // 1. Obtener parámetro
                String query = exchange.getRequestURI().getQuery();
                String n = (query != null && query.startsWith("n=")) ? query.substring(2) : "0";
                
                // 2. Obtener resultado SOAP
                String eng = callSoap(n);
                
                // 3. Traducir usando MyMemory (API abierta y amigable)
                String esp = translateMyMemory(eng);
                
                // 4. Enviar respuesta limpia
                exchange.getResponseHeaders().set("Content-Type", "text/plain; charset=utf-8");
                byte[] bytes = esp.getBytes(StandardCharsets.UTF_8);
                exchange.sendResponseHeaders(200, bytes.length);
                try (OutputStream os = exchange.getResponseBody()) { os.write(bytes); }
                
            } catch (Exception e) {
                String errorMsg = "Error: " + e.getMessage();
                exchange.getResponseHeaders().set("Content-Type", "text/plain; charset=utf-8");
                exchange.sendResponseHeaders(500, errorMsg.length());
                try (OutputStream os = exchange.getResponseBody()) { os.write(errorMsg.getBytes()); }
            }
        });
        server.start();
        System.out.println("Servidor iniciado en http://localhost:8000");
    }

    // SOAP igual que siempre
    private static String callSoap(String n) throws Exception {
        String url = "https://www.dataaccess.com/webservicesserver/NumberConversion.wso";
        String env = "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><NumberToWords xmlns=\"http://www.dataaccess.com/webservicesserver/\"><ubiNum>"+n+"</ubiNum></NumberToWords></soap:Body></soap:Envelope>";
        HttpURLConnection con = (HttpURLConnection) new URL(url).openConnection();
        con.setDoOutput(true); con.setRequestMethod("POST");
        con.setRequestProperty("Content-Type", "text/xml; charset=utf-8");
        con.setRequestProperty("SOAPAction", "http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords");
        try(OutputStream os = con.getOutputStream()){ os.write(env.getBytes("UTF-8")); }
        BufferedReader br = new BufferedReader(new InputStreamReader(con.getInputStream()));
        StringBuilder sb = new StringBuilder(); String line; while((line=br.readLine())!=null) sb.append(line);
        Matcher m = Pattern.compile("<[^>]*NumberToWordsResult>([^<]+)</[^>]*NumberToWordsResult>").matcher(sb.toString());
        return m.find() ? m.group(1).trim() : "Error";
    }

    // MÉTODO NUEVO: API de MyMemory
    private static String translateMyMemory(String text) throws Exception {
        // Codificar el texto para la URL
        String encodedText = URLEncoder.encode(text, StandardCharsets.UTF_8.toString());
        // URL de MyMemory (no requiere API Key para uso básico)
        String urlStr = "https://api.mymemory.translated.net/get?q=" + encodedText + "&langpair=en|es";
        
        URL url = new URL(urlStr);
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod("GET");
        conn.setRequestProperty("User-Agent", "Mozilla/5.0");

        BufferedReader br = new BufferedReader(new InputStreamReader(conn.getInputStream(), "UTF-8"));
        StringBuilder res = new StringBuilder();
        String line;
        while ((line = br.readLine()) != null) res.append(line);
        
        // Extraer texto traducido del JSON de MyMemory
        // Formato: {"responseData":{"translatedText":"diez"},...}
        Matcher m = Pattern.compile("\"translatedText\"\\s*:\\s*\"([^\"]+)\"").matcher(res.toString());
        
        if (m.find()) {
            return m.group(1); 
        }
        return text; // Fallback al original si falla
    }
}