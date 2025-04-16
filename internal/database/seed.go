package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

type UserSeed struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Department string `json:"department"`
	UserStatus string `json:"userStatus"`
	UserName   string `json:"username"`
}

func SeedDB(db *sql.DB) error {
	seedData, err := os.ReadFile("./seed.json")
	if err != nil {
		return fmt.Errorf("failed to read seed data: %w", err)
	}

	var users []UserSeed
    if err := json.Unmarshal(seedData, &users); err != nil {
        return fmt.Errorf("could not parse seed JSON: %w", err)
    }

	for _, user := range users {
		_, err := db.Exec(
			"INSERT OR IGNORE INTO users (first_name, last_name, email, department, user_status, user_name) VALUES (?, ?, ?, ?, ?, ?)",
			user.FirstName, user.LastName, user.Email, user.Department, user.UserStatus, user.UserName)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	return nil
}
