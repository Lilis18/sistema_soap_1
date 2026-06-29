require 'webrick'
require 'humanize'

# Configurar humanize para español
Humanize.configure do |config|
  config.default_locale = :es
end

server = WEBrick::HTTPServer.new(Port: 8000)

server.mount_proc '/' do |req, res|
  n = req.query['n'] || "0"
  
  # Conversión nativa
  resultado = n.to_i.humanize
  
  res.body = resultado
end

trap('INT') { server.shutdown }
server.start