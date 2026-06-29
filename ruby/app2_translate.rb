require 'webrick'
require 'net/http'
require 'json'

server = WEBrick::HTTPServer.new(Port: 8000)

server.mount_proc '/' do |req, res|
  n = req.query['n'] || '0'
  
  # Paso 1: Llamada SOAP Real
  uri = URI("https://www.dataaccess.com/webservicesserver/NumberConversion.wso")
  http = Net::HTTP.new(uri.host, uri.port)
  http.use_ssl = true
  req_soap = Net::HTTP::Post.new(uri.path)
  req_soap['Content-Type'] = 'text/xml'
  req_soap.body = "<soap:Envelope xmlns:soap='http://schemas.xmlsoap.org/soap/envelope/'><soap:Body><NumberToWords xmlns='http://www.dataaccess.com/webservicesserver/'><ubiNum>#{n}</ubiNum></NumberToWords></soap:Body></soap:Envelope>"
  
  soap_resp = http.request(req_soap)
  eng = soap_resp.body.match(/>([^<]+)<\/.*Result>/)[1]
  
  # Paso 2: Llamada API Traducción Real
  trans_url = URI("https://api.mymemory.translated.net/get?q=#{eng}&langpair=en|es")
  trans_resp = Net::HTTP.get(trans_url)
  esp = JSON.parse(trans_resp)["responseData"]["translatedText"]
  
  res.body = esp
end

server.start