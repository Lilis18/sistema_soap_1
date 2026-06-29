use axum::{extract::Query, routing::get, Router};
use std::collections::HashMap;

async fn get_soap(n: &str) -> String {
    let xml = format!(r#"<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>{}</ubiNum></NumberToWords></soap:Body></soap:Envelope>"#, n);
    let client = reqwest::Client::new();
    let resp = client.post("https://www.dataaccess.com/webservicesserver/NumberConversion.wso")
        .header("Content-Type", "text/xml; charset=utf-8")
        .body(xml).send().await.unwrap().text().await.unwrap();
    regex::Regex::new(r"<[^>]+:NumberToWordsResult>(.*?)</[^>]+:NumberToWordsResult>").unwrap()
        .captures(&resp).map(|c| c[1].to_string()).unwrap_or("Error".into())
}

async fn translate(eng: &str) -> String {
    let url = format!("https://api.mymemory.translated.net/get?q={}&langpair=en|es", eng);
    let resp = reqwest::get(url).await.unwrap().text().await.unwrap();
    regex::Regex::new(r#""translatedText":"([^"]+)""#).unwrap()
        .captures(&resp).map(|c| c[1].to_string()).unwrap_or(eng.into())
}

async fn handler(Query(params): Query<HashMap<String, String>>) -> String {
    let n = params.get("n").map(|s| s.as_str()).unwrap_or("0");
    let eng = get_soap(n).await;
    translate(&eng).await
}

#[tokio::main]
async fn main() {
    let app = Router::new().route("/", get(handler));
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8000").await.unwrap();
    println!("Servidor V2 (SOAP+Trad) activo en http://localhost:8000");
    axum::serve(listener, app).await.unwrap();
}