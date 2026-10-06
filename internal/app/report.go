package app

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"budget-cli/internal/engine"
)

func (a *App) handleReport(args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)

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

	r := engine.NewReport()

	tagReport, err := r.CalculateTagReport(transactions)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

	fmt.Println("📊 Report")
	fmt.Fprintln(w, "Tags\tAmount")

	for tag, amount := range tagReport {
		fmt.Fprintf(w, "%v\t%v\n",
			tag, amount,
		)
	}

	fmt.Fprintf(w, "---------------------------------\n")

	fmt.Fprintln(w, "Wallet\tAmount")

	walletReport, err := r.CalculateWalletReport(transactions)
	if err != nil {
		return err
	}

	for wallet, amount := range walletReport {
		fmt.Fprintf(w, "%v\t%v\n",
			wallet, amount)
	}

	fmt.Fprintf(w, "---------------------------------\n")
	fmt.Fprintln(w, "Mounth\tAmount")

	monthlyReport, mounths, err := r.GenerateMonthlyReport(transactions)
	if err != nil {
		return err
	}

	for _, month := range mounths {
		fmt.Fprintf(w, "%s\t%.2f\n", month, monthlyReport[month])
	}

	fmt.Fprintf(w, "---------------------------------\n")
	fmt.Fprintln(w, "TotalLost\tTotalGain\tTotal")

	total, err := r.GetTotal(transactions)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%v\t%v\t%v\n", total.TotalLost, total.TotalGain, total.Total)

	w.Flush()

	return nil
}
