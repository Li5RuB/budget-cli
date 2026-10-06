package app

import (
	"flag"
	"fmt"
)

func (a *App) handleRemove(args []string) error {

	fs := flag.NewFlagSet("remove", flag.ContinueOnError)

	var id int64

	fs.Int64Var(&id, "id", 0, "id required")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := a.store.RemoveTransaction(id); err != nil {
		return fmt.Errorf("Ca: %w", err)
	}

	fmt.Printf("✅ Success! Removed transaction %v\n", id)
	return nil
}
