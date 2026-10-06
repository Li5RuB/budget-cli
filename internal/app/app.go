package app

import (
	"database/sql"
	"fmt"

	"budget-cli/internal/db"
)

type App struct {
	store *db.Store
}

func NewApp(database *sql.DB) *App {
	return &App{
		store: db.NewStore(database),
	}
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return nil
	}

	subcommand := args[0]

	subArgs := args[1:]

	switch subcommand {
	case "add":
		return a.handleAdd(subArgs)
	case "report":
		return a.handleReport(subArgs)
	case "remove":
		return a.handleRemove(subArgs)
	case "get":
		return a.handleGet(subArgs)
	case "help", "-h", "--help":
		a.printHelp()
		return nil
	default:
		return fmt.Errorf("Unknow comand: %q", subcommand)
	}
}

func (a *App) printHelp() {
	fmt.Println("budget <comand> [args]")
	fmt.Println("comands:")
	fmt.Println("add    ----add transaction")
	fmt.Println("report    ----report by transactions")
	fmt.Println("budget add -amount 1500 -tags 'food,cafe' -comment 'lunch'")
}
