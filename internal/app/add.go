package app

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"budget-cli/internal/engine"
)

func (a *App) handleAdd(args []string) error {

	fs := flag.NewFlagSet("add", flag.ContinueOnError)

	var amount float64
	var tagsRaw string
	var comment string
	var createdAtStr string
	var curId int64
	var wallet string

	fs.Float64Var(&amount, "a", 0.0, "amount required")
	fs.StringVar(&wallet, "w", "main", "wallet")
	fs.StringVar(&tagsRaw, "t", "", "'food,cafe'")
	fs.StringVar(&comment, "com", "", "comment")
	fs.Int64Var(&curId, "c", 0, "NBRB cur_id")
	fs.StringVar(&createdAtStr, "d", "", "YYYY-MM-DD")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var txDate time.Time
	if createdAtStr != "" {
		var err error
		txDate, err = time.Parse("2006-01-02", createdAtStr)
		if err != nil {
			return fmt.Errorf("date format: YYYY-MM-DD: %v", err)
		}
	} else {
		txDate = time.Now()
	}

	var tags []string
	if tagsRaw != "" {
		rawSlice := strings.Split(tagsRaw, ",")
		for _, tag := range rawSlice {
			cleanTag := strings.TrimSpace(tag)
			if cleanTag != "" {
				tags = append(tags, cleanTag)
			}
		}
	}

	tx := engine.Transaction{
		Amount:   amount,
		Tags:     tags,
		Comment:  comment,
		CreateAt: txDate,
		CurId:    curId,
		Wallet:   wallet,
	}

	if err := a.store.SaveTransaction(tx); err != nil {
		return fmt.Errorf("Ca: %w", err)
	}

	fmt.Printf("✅ Success! Amount: %.2f, Tags: %v\n", tx.Amount, tx.Tags)
	return nil
}
