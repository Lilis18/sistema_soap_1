using Humanizer; // Requiere: dotnet add package Humanizer

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/", (HttpContext ctx) => {
    string nStr = ctx.Request.Query["n"].FirstOrDefault() ?? "0";
    if (int.TryParse(nStr, out int n)) {
        // La librería Humanizer se encarga de todo sin que tú escribas arreglos
        // ToWords() convierte el número a letras automáticamente
        string resultado = n.ToWords(new System.Globalization.CultureInfo("es"));
        return Results.Text(resultado, "text/plain; charset=utf-8");
    }
    return Results.BadRequest("Número inválido");
});

app.Run();