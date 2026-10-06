package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("Can't open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Can't connect to db: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("Can't create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
	transaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
	cur_id INTEGER DEFAULT 0,
	wallet TEXT NOT NULL,
	amount REAL NOT NULL,
	tags TEXT NOT NULL,
	comment TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}
