import http from 'http';
import { translate } from '@vitalets/google-translate-api'; // Librería estable
import fetch from 'node-fetch';

const WSDL_URL = 'https://www.dataaccess.com/webservicesserver/NumberConversion.wso';

async function callNumberToWords(n) {
  const body = `<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>${n}</ubiNum></NumberToWords></soap:Body></soap:Envelope>`;
  const res = await fetch(WSDL_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'text/xml; charset=utf-8', 'SOAPAction': 'http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords' },
    body: body
  });
  const text = await res.text();
  const m = text.match(/<[^>]*NumberToWordsResult>([^<]+)<\/[^>]*NumberToWordsResult>/);
  return m ? m[1].trim() : "Error SOAP";
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const n = url.searchParams.get('n') || '0';

  try {
    const eng = await callNumberToWords(n);
    // Traducción robusta
    const result = await translate(eng, { from: 'en', to: 'es' });
    
    res.writeHead(200, { 'Content-Type': 'text/plain; charset=utf-8' });
    res.end(result.text);
  } catch (e) {
    res.writeHead(500);
    res.end('Error: ' + e.message);
  }
});

server.listen(8000, () => console.log('Servidor V2 (Traductor) listo en http://localhost:8000/?n=10'));