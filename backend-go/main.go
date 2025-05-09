package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("sqlite3", "../database/login.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()
	app.Use(cors.New())

	app.Post("/api/login", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid input"})
		}

		var hashedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", body.Username).Scan(&hashedPassword)
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{"success": false, "message": "User not found"})
		} else if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Database error"})
		}

		if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body.Password)) != nil {
			return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid credentials"})
		}

		return c.JSON(fiber.Map{"success": true, "message": "Login successful"})
	})

	app.Post("/api/register", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid form"})
		}

		username := form.Value["username"][0]
		password := form.Value["password"][0]

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// Upload file
		var profilePath, idPath string
		if files := form.File["profilePhoto"]; len(files) > 0 {
			filename := fmt.Sprintf("profile_%s.jpg", username)
			path := filepath.Join("../database/uploads/profile", filename)
			profilePath = "uploads/profile/" + filename
			if err := c.SaveFile(files[0], path); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to save profile photo"})
			}
		}
		if files := form.File["idPhoto"]; len(files) > 0 {
			filename := fmt.Sprintf("id_%s.jpg", username)
			path := filepath.Join("../database/uploads/id", filename)
			idPath = "uploads/id/" + filename
			if err := c.SaveFile(files[0], path); err != nil {
				return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to save ID photo"})
			}
		}

		stmt, err := db.Prepare("INSERT INTO users (username, password, profile_photo_path, id_photo_path) VALUES (?, ?, ?, ?)")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to prepare query"})
		}
		_, err = stmt.Exec(username, string(hashedPassword), profilePath, idPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to create user / duplicate username"})
		}

		return c.JSON(fiber.Map{"success": true, "message": "User registered successfully"})
	})

	log.Println("Fiber server running at http://localhost:3001")
	log.Fatal(app.Listen(":3001"))
}
