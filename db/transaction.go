package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Transaction struct {
	*Queries
	db *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{
		Queries: New(db),
		db:      db,
	}
}

func (transaction *Transaction) ExecTx(ctx context.Context, fn func(*Queries) error, isTest bool) error {
	var rbErr error
	tx, err := transaction.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if isTest {
			rbErr = tx.Rollback()
		}
		rbErr = tx.Commit()
	}()

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	
	return rbErr
}
