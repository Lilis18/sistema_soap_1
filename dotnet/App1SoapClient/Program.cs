using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using System.Text;
using System.Text.Json;
using System.Text.RegularExpressions;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();
var httpClient = new HttpClient();

app.MapGet("/", async (HttpContext ctx) => {
    string n = ctx.Request.Query["n"].FirstOrDefault() ?? "0";
    string envelope = $"<?xml version=\"1.0\" encoding=\"utf-8\"?><soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><NumberToWords xmlns=\"http://www.dataaccess.com/webservicesserver/\"><ubiNum>{n}</ubiNum></NumberToWords></soap:Body></soap:Envelope>";

    var req = new HttpRequestMessage(HttpMethod.Post, "https://www.dataaccess.com/webservicesserver/NumberConversion.wso");
    req.Content = new StringContent(envelope, Encoding.UTF8, "text/xml");
    req.Headers.Add("SOAPAction", "http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords");

    var resp = await httpClient.SendAsync(req);
    var text = await resp.Content.ReadAsStringAsync();
    var m = Regex.Match(text, "<[^>]*NumberToWordsResult>([^<]+)</[^>]*NumberToWordsResult>");
    
    await ctx.Response.WriteAsync(m.Success ? m.Groups[1].Value.Trim() : "Error SOAP");
});

app.Run();