require 'webrick'
require 'net/http'
require 'uri'

server = WEBrick::HTTPServer.new(Port: 8000)

server.mount_proc '/' do |req, res|
  n = req.query['n'] || '0'
  uri = URI("https://www.dataaccess.com/webservicesserver/NumberConversion.wso")
  
  http = Net::HTTP.new(uri.host, uri.port)
  http.use_ssl = true
  http.open_timeout = 5 # Timeout real de 5 segundos
  
  begin
    request = Net::HTTP::Post.new(uri.path)
    request['Content-Type'] = 'text/xml; charset=utf-8'
    request['SOAPAction'] = 'http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords'
    request.body = "<soap:Envelope xmlns:soap='http://schemas.xmlsoap.org/soap/envelope/'><soap:Body><NumberToWords xmlns='http://www.dataaccess.com/webservicesserver/'><ubiNum>#{n}</ubiNum></NumberToWords></soap:Body></soap:Envelope>"
    
    response = http.request(request)
    match = response.body.match(/<[^>]+:NumberToWordsResult>(.*?)<\/[^>]+:NumberToWordsResult>/)
    res.body = match ? match[1] : "Error parsing XML"
  rescue => e
    res.status = 500
    res.body = "Error de conexión real: #{e.message}"
  end
end

trap('INT') { server.shutdown }
server.start