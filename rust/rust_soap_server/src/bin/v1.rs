use axum::{extract::Query, routing::get, Router};
use std::collections::HashMap;

async fn handler(Query(params): Query<HashMap<String, String>>) -> String {
    let n = params.get("n").map(|s| s.as_str()).unwrap_or("0");
    let soap_xml = format!(r#"<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><NumberToWords xmlns="http://www.dataaccess.com/webservicesserver/"><ubiNum>{}</ubiNum></NumberToWords></soap:Body></soap:Envelope>"#, n);

    let client = reqwest::Client::new();
    let resp = client.post("https://www.dataaccess.com/webservicesserver/NumberConversion.wso")
        .header("Content-Type", "text/xml; charset=utf-8")
        .header("SOAPAction", "http://www.dataaccess.com/webservicesserver/NumberConversion.wso/NumberToWords")
        .body(soap_xml).send().await.unwrap().text().await.unwrap();

    let re = regex::Regex::new(r"<[^>]+:NumberToWordsResult>(.*?)</[^>]+:NumberToWordsResult>").unwrap();
    re.captures(&resp).map(|c| c[1].to_string()).unwrap_or("Error".into())
}

#[tokio::main]
async fn main() {
    let app = Router::new().route("/", get(handler));
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8000").await.unwrap();
    println!("Servidor V1 (SOAP) activo en http://localhost:8000");
    axum::serve(listener, app).await.unwrap();
} 