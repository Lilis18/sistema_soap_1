use axum::{extract::Query, routing::get, Router};
use std::collections::HashMap;

fn to_spanish(n: i32) -> String {
    match n {
        0 => "cero".into(),
        1 => "uno".into(),
        2 => "dos".into(),
        3 => "tres".into(),
        4 => "cuatro".into(),
        5 => "cinco".into(),
        10 => "diez".into(),
        11 => "once".into(),
        20 => "veinte".into(),
        21..=29 => format!("veinti{}", to_spanish(n - 20)),
        30..=39 => format!("treinta y {}", to_spanish(n - 30)),
        _ => n.to_string(),
    }
}

async fn handler(Query(params): Query<HashMap<String, String>>) -> String {
    let n_str = params.get("n").map(|s| s.as_str()).unwrap_or("0");
    let n: i32 = n_str.parse().unwrap_or(0);
    to_spanish(n)
}

#[tokio::main]
async fn main() {
    let app = Router::new().route("/", get(handler));
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8000").await.unwrap();
    println!("Servidor V3 (Nativo) activo en http://localhost:8000");
    axum::serve(listener, app).await.unwrap();
}