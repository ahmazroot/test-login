package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("sqlite3", "login.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			profile_photo_path TEXT,
			id_photo_path TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create test users with bcrypt hashed passwords
	users := []struct {
		username string
		password string
	}{
		{"test", "test123"},
		{"admin", "admin123"},
	}

	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(
			"INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)",
			user.username,
			string(hashedPassword),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Database initialized successfully with test users:")
	log.Println("1. username: test, password: test123")
	log.Println("2. username: admin, password: admin123")
}
