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

var testQuries *Queries

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

func ExecTestingTx(ctx context.Context, transaction *Transaction, fn func() error, isTest bool) error {
	var rbErr error

	tx, txerr := transaction.db.BeginTx(ctx, nil)

	if txerr != nil {
		return txerr
	}

	fmt.Printf("tx start : %v\n", tx)

	testQuries = New(tx)

	err := fn()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if isTest {
		fmt.Printf("tx end and rollback because of test : %v\n", tx)
		rbErr = tx.Rollback()
	} else {
		rbErr = tx.Commit()
	}

	return rbErr
}
