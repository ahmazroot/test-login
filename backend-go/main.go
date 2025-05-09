package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}


func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "../database/login.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit for file uploads
	})

	// CORS middleware
	app.Use(cors.New())

	// Serve static files
	app.Static("/uploads", "../database/uploads")

	app.Post("/api/register", func(c *fiber.Ctx) error {
		// Parse multipart form
		form, err := c.MultipartForm()
		if err != nil {
                return c.Status(fiber.StatusBadRequest).JSON(RegisterResponse{
                    Success: false,
                    Message: "Invalid form data",
                })
            }

            // Validate required fields
            if len(form.Value["username"]) == 0 || len(form.Value["password"]) == 0 {
                return c.Status(fiber.StatusBadRequest).JSON(RegisterResponse{
                    Success: false,
                    Message: "Username and password are required",
                })
            }

		// Get form fields
		username := form.Value["username"][0]
		password := form.Value["password"][0]

		// Handle profile photo
		var profilePhotoPath string
		if profileFiles := form.File["profilePhoto"]; len(profileFiles) > 0 {
			profileFile := profileFiles[0]
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(profileFile.Filename))
			profilePhotoPath = fmt.Sprintf("uploads/profile/%s", filename)
			if err := c.SaveFile(profileFile, fmt.Sprintf("../database/%s", profilePhotoPath)); err != nil {
				return c.JSON(RegisterResponse{
					Success: false,
					Message: "Error saving profile photo",
				})
			}
		}

		// Handle ID photo
		var idPhotoPath string
		if idFiles := form.File["idPhoto"]; len(idFiles) > 0 {
			idFile := idFiles[0]
			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(idFile.Filename))
			idPhotoPath = fmt.Sprintf("uploads/id/%s", filename)
			if err := c.SaveFile(idFile, fmt.Sprintf("../database/%s", idPhotoPath)); err != nil {
				return c.JSON(RegisterResponse{
					Success: false,
					Message: "Error saving ID photo",
				})
			}
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(RegisterResponse{
				Success: false,
				Message: "Error processing password",
			})
		}

		// Insert user into database
		_, err = db.Exec(
			"INSERT INTO users (username, password, profile_photo_path, id_photo_path) VALUES (?, ?, ?, ?)",
			username,
			string(hashedPassword),
			profilePhotoPath,
			idPhotoPath,
		)
		if err != nil {
			return c.JSON(RegisterResponse{
				Success: false,
				Message: "Error creating user",
			})
		}

		return c.JSON(RegisterResponse{
			Success: true,
			Message: "User registered successfully",
		})
	})

	app.Post("/api/login", func(c *fiber.Ctx) error {
		var request LoginRequest
		if err := c.BodyParser(&request); err != nil {
			return c.JSON(LoginResponse{
				Success: false,
				Message: "Invalid request",
			})
		}

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", request.Username).Scan(&storedPassword)
		if err == sql.ErrNoRows {
			return c.JSON(LoginResponse{
				Success: false,
				Message: "User not found",
			})
		} else if err != nil {
			return c.JSON(LoginResponse{
				Success: false,
				Message: "Database error",
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(request.Password))
		if err != nil {
			return c.JSON(LoginResponse{
				Success: false,
				Message: "Invalid credentials",
			})
		}

		return c.JSON(LoginResponse{
			Success: true,
			Message: "Login successful",
		})
	})

	log.Println("Go Fiber server running on http://localhost:7002")
	log.Fatal(app.Listen(":7002"))
}