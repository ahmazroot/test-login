#[macro_use] extern crate rocket;

use rocket::serde::{json::Json, Deserialize, Serialize};
use rocket::State;
use rocket::http::Status;
use rocket::fs::TempFile;
use rocket::form::Form;
use sqlx::sqlite::SqlitePool;
use bcrypt::{hash,verify,DEFAULT_COST};
use rocket_cors::{AllowedOrigins, CorsOptions};
use std::path::Path;

// Struktur untuk request login
#[derive(Deserialize)]
#[serde(crate = "rocket::serde")]
struct LoginRequest {
    username: String,
    password: String,
}

// Struktur untuk response login
#[derive(Serialize)]
#[serde(crate = "rocket::serde")]
struct LoginResponse {
    success: bool,
    message: String,
}

// Struktur untuk response register
#[derive(Serialize)]
#[serde(crate = "rocket::serde")]
struct RegisterResponse {
    success: bool,
    message: String,
}

#[derive(FromForm)]
struct RegisterForm<'r> {
    username: &'r str,
    password: &'r str,
    profilePhoto: TempFile<'r>,
    idPhoto: TempFile<'r>,
}

// Endpoint register
#[post("/register", data = "<form>")]
async fn register(
    pool: &State<SqlitePool>,
    mut form: Form<RegisterForm<'_>>,
) -> Result<Json<RegisterResponse>, Status> {


    let profile_filename = format!("profile_{}.jpg", form.username);
    let profile_path = Path::new("../database/uploads/profile").join(&profile_filename);
    form.profilePhoto.persist_to(&profile_path).await
        .map_err(|_| Status::InternalServerError)?;


    let id_filename = format!("id_{}.jpg", form.username);
    let id_path = Path::new("../database/uploads/id").join(&id_filename);
    form.idPhoto.persist_to(&id_path).await
        .map_err(|_| Status::InternalServerError)?;


    let hashed_password = hash(form.password, DEFAULT_COST)
        .map_err(|_| Status::InternalServerError)?;


    sqlx::query(
        "INSERT INTO users (username, password, profile_photo_path, id_photo_path) VALUES (?, ?, ?, ?)",
    )
    .bind(form.username)
    .bind(&hashed_password)
    .bind(format!("uploads/profile/{}", profile_filename))
    .bind(format!("uploads/id/{}", id_filename))
    .execute(&**pool)
    .await
    .map_err(|_| Status::InternalServerError)?;

    Ok(Json(RegisterResponse {
        success: true,
        message: "User registered successfully".to_string(),
    }))
}


#[post("/login", data = "<login_request>")]
async fn login(pool: &State<SqlitePool>, login_request: Json<LoginRequest>) -> Result<Json<LoginResponse>, Status> {
    let user = sqlx::query_as::<_, (String,)>(
        "SELECT password FROM users WHERE username = ?",
    )
    .bind(&login_request.username)
    .fetch_optional(&**pool)
    .await
    .map_err(|_| Status::InternalServerError)?;

    match user {
        Some((hashed_password,)) => {
            if verify(&login_request.password, &hashed_password).unwrap_or(false) {
                Ok(Json(LoginResponse {
                    success: true,
                    message: "Login successful".to_string(),
                }))
            } else {
                Ok(Json(LoginResponse {
                    success: false,
                    message: "Invalid credentials".to_string(),
                }))
            }
        }
        None => Ok(Json(LoginResponse {
            success: false,
            message: "User not found".to_string(),
        })),
    }
}

// Launch Rocket
#[launch]
async fn rocket() -> _ {
    let pool = SqlitePool::connect("sqlite:../database/login.db")
        .await
        .expect("Failed to connect to SQLite");

    let cors = CorsOptions::default()
        .allowed_origins(AllowedOrigins::all())
        .to_cors()
        .expect("Failed to create CORS configuration");

    rocket::build()
        .mount("/api", routes![login, register])
        .manage(pool)
        .attach(cors)
}