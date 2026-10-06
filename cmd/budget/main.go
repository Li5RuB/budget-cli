package main

import (
	"fmt"
	"log"
	"os"

	"budget-cli/internal/app"
	"budget-cli/internal/db"
)

func main() {
	database, err := db.InitDB("budget.db")
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
