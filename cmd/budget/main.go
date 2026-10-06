package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"budget-cli/internal/app"
	"budget-cli/internal/db"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Can't take user directory: %v", err)
	}

	dbPath := filepath.Join(homeDir, ".budget.db")
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("init db error: %v", err)
	}

	defer database.Close()

	cliApp := app.NewApp(database)

	err = cliApp.Run(os.Args[1:])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
