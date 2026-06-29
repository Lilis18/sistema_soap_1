import com.sun.net.httpserver.*;
import com.ibm.icu.text.RuleBasedNumberFormat; // Librería estándar de traducción
import java.net.InetSocketAddress;
import java.util.Locale;
import java.io.*;

public class App3NativeServer {
    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress(8000), 0);
        
        server.createContext("/", exchange -> {
            String q = exchange.getRequestURI().getQuery();
            long n = 0; 
            if(q != null && q.startsWith("n=")) { 
                try { n = Long.parseLong(q.substring(2)); } catch(Exception ignored) {} 
            }

            // Uso de la lógica similar a PHP (NumberFormatter)
            // Esto convierte números a palabras en español automáticamente
            RuleBasedNumberFormat formatter = new RuleBasedNumberFormat(new Locale("es"), RuleBasedNumberFormat.SPELLOUT);
            String resp = formatter.format(n);
            
            exchange.getResponseHeaders().set("Content-Type", "text/plain; charset=utf-8");
            byte[] bytes = resp.getBytes("UTF-8");
            exchange.sendResponseHeaders(200, bytes.length);
            try(OutputStream os = exchange.getResponseBody()){ os.write(bytes); }
        });
        
        server.start();
        System.out.println("Java Native Server (ICU) iniciado en http://localhost:8000");
    }
}