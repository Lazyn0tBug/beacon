use axum::{
    async_trait,
    body::{Body, Bytes},
    extract::{FromRequest, Path, Query, Request},
    http::{HeaderName, StatusCode},
    middleware::{self, Next},
    response::{Html, IntoResponse, Json, Response},
    routing::{get, post},
    Router,
};
use http_body_util::BodyExt;
use hyper::{header, HeaderMap};
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::{
    collections::HashMap,
    net::{Ipv4Addr, Ipv6Addr, SocketAddr},
    process::id,
    time::Duration,
};
use tokio::net::TcpListener;
use tokio::signal;
use tokio::time::sleep;
use tower::{Layer, ServiceBuilder};
use tower_http::{
    compression::CompressionLayer,
    cors::{Any, CorsLayer},
    request_id::{MakeRequestUuid, SetRequestIdLayer},
    timeout::TimeoutLayer,
    trace::{self, TraceLayer},
};
use tracing::Level;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() {
    // initialize tracing
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "example_consume_body_in_extractor_or_middleware=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    let trace_layer = TraceLayer::new_for_http()
        .make_span_with(trace::DefaultMakeSpan::new().level(Level::INFO))
        .on_request(trace::DefaultOnRequest::new().level(Level::INFO))
        .on_response(trace::DefaultOnResponse::new().level(Level::INFO));

    let port = 5789;

    // let localhost_v6 = SocketAddr::new(Ipv6Addr::LOCALHOST.into(), port);
    // let listener_v6 = TcpListener::bind(&localhost_v6).await.unwrap();
    let id_routes = Router::new().route("/:id", get(user_details));

    // build our application with a route
    let app = Router::new()
        // `GET /` goes to `root`
        .route("/", get(root))
        .nest("/api", id_routes)
        // `POST /users` goes to `create_user`
        .route("/users", post(create_user))
        .route("/users/:name", get(get_user))
        .route("/query", post(query_params))
        .route("/headers", get(parse_headers))
        .layer(
            ServiceBuilder::new()
                .layer(trace_layer)
                .layer(CompressionLayer::new().zstd(true).gzip(true))
                .layer(TimeoutLayer::new(Duration::new(0, 200000)))
                .layer(SetRequestIdLayer::new(
                    HeaderName::from_static("x-request-id"),
                    MakeRequestUuid,
                ))
                .layer(
                    CorsLayer::new()
                        .allow_methods(axum::http::Method::GET)
                        .allow_origin(Any),
                ),
        )
        .fallback(api_fallback);

    // run our app with hyper
    let localhost_v4 = SocketAddr::new(Ipv4Addr::LOCALHOST.into(), port);
    let listener_v4 = TcpListener::bind(&localhost_v4).await.unwrap();

    tracing::debug!("listening on {}", listener_v4.local_addr().unwrap());
    axum::serve(listener_v4, app)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

// basic handler that responds with a static string
async fn root() -> Html<&'static str> {
    Html("<h1>Hello, World!</h1>")
}

async fn user_details(Path(id): Path<u32>) -> impl IntoResponse {
    Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "text/html")
        .body(Body::from(format!("<h1>User detail for {}</h1>", id)))
        .unwrap()
        .into_response()
}

async fn query_params(query: Query<HashMap<String, String>>) -> impl IntoResponse {
    let info = query.0;
    format!("{info:?}")
}

async fn parse_headers(headers: HeaderMap) -> String {
    format!("{headers:?}")
}

async fn get_user(Path(name): Path<String>) -> Json<User> {
    let user = User {
        id: 1228,
        username: name,
    };
    Json(user)
}

async fn create_user(
    // this argument tells axum to parse the request body
    // as JSON into a `CreateUser` type
    payload: Option<Json<CreateUser>>,
    // Json(payload): Json<CreateUser>,
) -> impl IntoResponse {
    if let Some(payload) = payload {
        // insert your application logic here
        let user = User {
            id: 1337,
            username: payload.username.clone(),
        };

        // Response::builder()
        //     .status(StatusCode::OK)
        //     .header(header::CONTENT_TYPE, "application/json")
        //     .body(Json(user))
        //     .unwrap()
        // this will be converted into a JSON response
        // with a status code of `201 Created`
        (StatusCode::CREATED, Json(user))
    } else {
        let user = User {
            id: 0,
            username: String::from("not found"),
        };
        // (StatusCode::NO_CONTENT, Json(json!({"id":0， "username": String::from("not found")})))
        (StatusCode::NOT_FOUND, Json(user))
    }
}

// the input to our `create_user` handler
#[derive(Deserialize)]
struct CreateUser {
    username: String,
}

// the output to our `create_user` handler
#[derive(Serialize)]
struct User {
    id: u64,
    username: String,
}

// middleware that shows how to consume the request body upfront
async fn print_request_body(request: Request, next: Next) -> Result<impl IntoResponse, Response> {
    let request = buffer_request_body(request).await?;

    Ok(next.run(request).await)
}

// the trick is to take the request apart, buffer the body, do what you need to do, then put
// the request back together
async fn buffer_request_body(request: Request) -> Result<Request, Response> {
    let (parts, body) = request.into_parts();

    // this wont work if the body is an long running stream
    let bytes = body
        .collect()
        .await
        .map_err(|err| (StatusCode::INTERNAL_SERVER_ERROR, err.to_string()).into_response())?
        .to_bytes();

    do_thing_with_request_body(bytes.clone());

    Ok(Request::from_parts(parts, Body::from(bytes)))
}

fn do_thing_with_request_body(bytes: Bytes) {
    tracing::debug!(body = ?bytes);
}

async fn api_fallback() -> impl IntoResponse {
    (
        StatusCode::NOT_FOUND,
        Json(serde_json::json!({ "status": "Not Found" })),
    )
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }
}
