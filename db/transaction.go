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
		db:      db,
	}
}

func (transaction *Transaction) ExecTx(ctx context.Context, fn func() error, isTest bool) error {
	var rbErr error
	tx, err := transaction.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Println("We are at a db transaction")
	// q := New(tx)
	err = fn()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	
	if isTest {
		rbErr = tx.Rollback()
	} else {
		rbErr = tx.Commit()
	}

	return rbErr
}
