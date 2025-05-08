use axum::{
    routing::post,
    http::StatusCode,
    Json, Router,
    extract::State,
};
use serde::{Deserialize, Serialize};
use sqlx::sqlite::SqlitePool;
use tower_http::cors::{CorsLayer, Any};
use bcrypt::{hash, DEFAULT_COST, verify};
use std::path::Path;
use tokio::fs;

#[derive(Deserialize)]
struct LoginRequest {
    username: String,
    password: String,
}

#[derive(Serialize)]
struct LoginResponse {
    success: bool,
    message: String,
}

#[derive(Serialize)]
struct RegisterResponse {
    success: bool,
    message: String,
}

async fn login_handler(
    State(pool): State<SqlitePool>,
    Json(payload): Json<LoginRequest>,
) -> Result<Json<LoginResponse>, (StatusCode, Json<LoginResponse>)> {
    let user = sqlx::query_as::<_, (String,)>(
        "SELECT password FROM users WHERE username = ?",
    )
    .bind(&payload.username)
    .fetch_optional(&pool)
    .await
    .map_err(|_| {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(LoginResponse {
                success: false,
                message: "Database error".to_string(),
            }),
        )
    })?;

    match user {
        Some((hashed_password,)) => {
            if verify(&payload.password, &hashed_password).unwrap_or(false) {
                Ok(Json(LoginResponse {
                    success: true,
                    message: "Login successful".to_string(),
                }))
            } else {
                Err((
                    StatusCode::UNAUTHORIZED,
                    Json(LoginResponse {
                        success: false,
                        message: "Invalid credentials".to_string(),
                    }),
                ))
            }
        }
        None => Err((
            StatusCode::UNAUTHORIZED,
            Json(LoginResponse {
                success: false,
                message: "User not found".to_string(),
            }),
        )),
    }
}

async fn register_handler(
    State(pool): State<SqlitePool>,
    mut multipart: axum::extract::Multipart,
) -> Result<Json<RegisterResponse>, (StatusCode, Json<RegisterResponse>)> {
    let mut username = String::new();
    let mut password = String::new();
    let mut profile_photo_path = None;
    let mut id_photo_path = None;

    while let Some(field) = multipart.next_field().await.unwrap() {
        let name = field.name().unwrap().to_string();
        match name.as_str() {
            "username" => username = field.text().await.unwrap(),
            "password" => password = field.text().await.unwrap(),
            "profilePhoto" => {
                let data = field.bytes().await.unwrap();
                let file_name = format!("profile_{}.jpg", username);
                let path = Path::new("../database/uploads/profile").join(&file_name);
                fs::write(&path, &data).await.unwrap();
                profile_photo_path = Some(format!("uploads/profile/{}", file_name));
            },
            "idPhoto" => {
                let data = field.bytes().await.unwrap();
                let file_name = format!("id_{}.jpg", username);
                let path = Path::new("../database/uploads/id").join(&file_name);
                fs::write(&path, &data).await.unwrap();
                id_photo_path = Some(format!("uploads/id/{}", file_name));
            },
            _ => {}
        }
    }

    let hashed_password = hash(&password, DEFAULT_COST)
        .map_err(|_| (StatusCode::INTERNAL_SERVER_ERROR, Json(RegisterResponse {
            success: false,
            message: "Failed to hash password".to_string(),
        })))?;

    sqlx::query(
        "INSERT INTO users (username, password, profile_photo_path, id_photo_path) VALUES (?, ?, ?, ?)",
    )
    .bind(&username)
    .bind(&hashed_password)
    .bind(&profile_photo_path)
    .bind(&id_photo_path)
    .execute(&pool)
    .await
    .map_err(|_| (StatusCode::INTERNAL_SERVER_ERROR, Json(RegisterResponse {
        success: false,
        message: "Failed to create user/ user duplicate".to_string(),
    })))?;

    Ok(Json(RegisterResponse {
        success: true,
        message: "User registered successfully".to_string(),
    }))
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let pool = SqlitePool::connect("sqlite:../database/login.db").await?;

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    let app = Router::new()
        .route("/api/login", post(login_handler))
        .route("/api/register", post(register_handler))
        .layer(cors)
        .with_state(pool);

    println!("Axum server running on http://localhost:3001");
    let listener = tokio::net::TcpListener::bind("0.0.0.0:7001").await?;
    axum::serve(listener, app).await?;

    Ok(())
}