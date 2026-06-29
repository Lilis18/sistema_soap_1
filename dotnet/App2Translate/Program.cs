using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using System.Text;
using System.Text.Json;
using System.Text.RegularExpressions;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();
var client = new HttpClient();

app.MapGet("/", async (HttpContext ctx) => {
    string n = ctx.Request.Query["n"].FirstOrDefault() ?? "0";
    
    // 1. SOAP
    string env = $"<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\"><soap:Body><NumberToWords xmlns=\"http://www.dataaccess.com/webservicesserver/\"><ubiNum>{n}</ubiNum></NumberToWords></soap:Body></soap:Envelope>";
    var req = new HttpRequestMessage(HttpMethod.Post, "https://www.dataaccess.com/webservicesserver/NumberConversion.wso") {
        Content = new StringContent(env, Encoding.UTF8, "text/xml")
    };
    req.Headers.Add("SOAPAction", "http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords");
    var resSoap = await client.SendAsync(req);
    var eng = Regex.Match(await resSoap.Content.ReadAsStringAsync(), "<[^>]*NumberToWordsResult>([^<]+)</[^>]*NumberToWordsResult>").Groups[1].Value;

    // 2. Traducir (MyMemory API)
    var transUrl = $"https://api.mymemory.translated.net/get?q={Uri.EscapeDataString(eng)}&langpair=en|es";
    var resTrans = await client.GetAsync(transUrl);
    var json = await resTrans.Content.ReadAsStringAsync();
    var doc = JsonDocument.Parse(json);
    string esp = doc.RootElement.GetProperty("responseData").GetProperty("translatedText").GetString() ?? eng;

    await ctx.Response.WriteAsync(esp);
});

app.Run();