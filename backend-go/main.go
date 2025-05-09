package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

type LoginRequest struct {
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

func loginHandler(c *fiber.Ctx, db *sql.DB) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(LoginResponse{
			Success: false,
			Message: "Invalid request",
		})
	}

	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
				Success: false,
				Message: "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(LoginResponse{
			Success: false,
			Message: "Database error",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	return c.JSON(LoginResponse{
		Success: true,
		Message: "Login successful",
	})
}

func registerHandler(c *fiber.Ctx, db *sql.DB) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(RegisterResponse{
			Success: false,
			Message: "Invalid form data",
		})
	}

	username := form.Value["username"][0]
	password := form.Value["password"][0]

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(RegisterResponse{
			Success: false,
			Message: "Failed to hash password",
		})
	}

	// Handle profile photo
	var profilePhotoPath string
	if profileFiles, ok := form.File["profilePhoto"]; ok && len(profileFiles) > 0 {
		profileFile := profileFiles[0]
		fileName := fmt.Sprintf("profile_%s%s", username, filepath.Ext(profileFile.Filename))
		path := filepath.Join("../database/uploads/profile", fileName)

		if err := saveUploadedFile(profileFile, path); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(RegisterResponse{
				Success: false,
				Message: "Failed to save profile photo",
			})
		}
		profilePhotoPath = filepath.Join("uploads/profile", fileName)
	}

	// Handle ID photo
	var idPhotoPath string
	if idFiles, ok := form.File["idPhoto"]; ok && len(idFiles) > 0 {
		idFile := idFiles[0]
		fileName := fmt.Sprintf("id_%s%s", username, filepath.Ext(idFile.Filename))
		path := filepath.Join("../database/uploads/id", fileName)

		if err := saveUploadedFile(idFile, path); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(RegisterResponse{
				Success: false,
				Message: "Failed to save ID photo",
			})
		}
		idPhotoPath = filepath.Join("uploads/id", fileName)
	}

	// Insert into database
	_, err = db.Exec(
		"INSERT INTO users (username, password, profile_photo_path, id_photo_path) VALUES (?, ?, ?, ?)",
		username, hashedPassword, profilePhotoPath, idPhotoPath,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).JSON(RegisterResponse{
				Success: false,
				Message: "Username already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(RegisterResponse{
			Success: false,
			Message: "Failed to create user",
		})
	}

	return c.JSON(RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

func saveUploadedFile(file *fiber.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func main() {
	app := fiber.New()

	// Initialize database
	db, err := sql.Open("sqlite3", "../database/login.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	// Routes
	app.Post("/api/login", func(c *fiber.Ctx) error {
		return loginHandler(c, db)
	})

	app.Post("/api/register", func(c *fiber.Ctx) error {
		return registerHandler(c, db)
	})

	fmt.Println("Fiber server running on http://localhost:7002")
	if err := app.Listen(":7002"); err != nil {
		panic(err)
	}
}