package app

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"budget-cli/internal/engine"
)

func (a *App) handleGet(args []string) error {

	fs := flag.NewFlagSet("get", flag.ContinueOnError)

	var tagFilter string
	var fromStr string
	var toStr string

	fs.StringVar(&tagFilter, "tag", "", "Filter by tags")
	fs.StringVar(&fromStr, "from", "", "From date: YYYY-MM-DD")
	fs.StringVar(&toStr, "to", "", "To date: YYYY-MM-DD")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var opts engine.FilterOptions
	opts.Tag = tagFilter

	const dateLayout = "2006-01-02"
	if fromStr != "" {
		t, err := time.Parse(dateLayout, fromStr)
		if err != nil {
			return fmt.Errorf("date format: YYYY-MM-DD: %v", err)
		}
		opts.FromDate = t
	}
	if toStr != "" {
		t, err := time.Parse(dateLayout, toStr)
		if err != nil {
			return fmt.Errorf("date format: YYYY-MM-DD: %v", err)
		}
		opts.ToDate = t
	}

	transactions, err := a.store.GetFilteredTransactions(opts)
	if err != nil {
		return fmt.Errorf("Can't get transactions: %w", err)
	}

	if len(transactions) == 0 {
		fmt.Println("📭 transactions not found 'budget add'")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Id\tAmount\tWallet\tTags\tCur_id\tComment\tCreateAt")

	for _, tx := range transactions {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%s\t%s\n",
			tx.TransactionId,
			tx.Amount,
			tx.Wallet,
			tx.Tags,
			tx.CurId,
			tx.Comment,
			tx.CreateAt.Format("2006-01-02"),
		)
	}

	w.Flush()

	return nil
}
