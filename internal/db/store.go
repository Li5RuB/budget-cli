package db

import (
	"database/sql"
	"strings"
	"time"

	"budget-cli/internal/engine"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) SaveTransaction(t engine.Transaction) error {
	tagsString := strings.Join(t.Tags, ",")

	query := `
	INSERT INTO transactions (amount, tags, comment, created_at, cur_id, wallet)
	VALUES (?,?,?,?,?,?);`

	if t.CreateAt.IsZero() {
		t.CreateAt = time.Now()
	}

	_, err := s.db.Exec(query, t.Amount, tagsString, t.Comment, t.CreateAt, t.CurId, t.Wallet)
	return err
}

func (s *Store) GetFilteredTransactions(opts engine.FilterOptions) ([]engine.Transaction, error) {
	query := `SELECT transaction_id, amount, tags, comment, created_at, cur_id, wallet FROM transactions WHERE 1=1`
	var args []interface{}

	if opts.Tag != "" {
		query += ` AND (tags = ? OR tags LIKE ? OR tags LIKE ? OR tags LIKE ?)`
		args = append(args,
			opts.Tag,
			opts.Tag+",%",
			"%,"+opts.Tag,
			"%,"+opts.Tag+",%",
		)
	}

	if !opts.FromDate.IsZero() {
		query += ` AND created_at >= ?`
		args = append(args, opts.FromDate)
	}
	if !opts.ToDate.IsZero() {
		query += ` AND created_at <= ?`
		args = append(args, opts.ToDate)
	}

	query += ` ORDER BY created_at DESC;`

	return GetTransactions(s, query, args)
}

func (s *Store) RemoveTransaction(id int64) error {
	query := `DELETE FROM Transactions WHERE transaction_id = ?`

	_, err := s.db.Query(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetTransactions(s *Store, query string, args []interface{}) ([]engine.Transaction, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []engine.Transaction

	for rows.Next() {
		var t engine.Transaction
		var tagsString string

		err := rows.Scan(&t.TransactionId, &t.Amount, &tagsString, &t.Comment, &t.CreateAt, &t.CurId, &t.Wallet)
		if err != nil {
			return nil, err
		}

		if tagsString != "" {
			t.Tags = strings.Split(tagsString, ",")
			for i := range t.Tags {
				t.Tags[i] = strings.TrimSpace(t.Tags[i])
			}
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
