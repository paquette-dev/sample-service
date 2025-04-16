package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
    "fmt"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name VARCHAR(50) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		department VARCHAR(255),
		user_status VARCHAR(1) NOT NULL
	);`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}


