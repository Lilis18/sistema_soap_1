import http from 'http';
import writtenNumber from 'written-number';

// Servidor ligero
const server = http.createServer((req, res) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const n = parseInt(url.searchParams.get('n') || '0', 10);

  // Conversión nativa utilizando la librería (equivalente a NumberFormatter::SPELLOUT)
  // 'lang: es' configura el idioma a español
  const result = writtenNumber(n, { lang: 'es' });

  res.writeHead(200, { 'Content-Type': 'text/plain; charset=utf-8' });
  res.end(result);
});

server.listen(8000, () => {
  console.log('Servidor nativo iniciado en http://localhost:8000');
});