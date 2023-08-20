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
		db: db,
	}
}

func ExecTestingTx(ctx context.Context, transaction *Transaction, fn func(*Queries) error, isTest bool) error {
	var rbErr error

	tx, txerr := transaction.db.BeginTx(ctx, nil)

	if txerr != nil {
		return txerr
	}

	fmt.Printf("tx start : %v\n", tx)

	q := New(tx)

	err := fn(q)
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

func (transaction *Transaction) ExecTx(ctx context.Context, fn func(*Queries) error, needRollback bool) error {
	tx, txerr := transaction.db.BeginTx(ctx, nil)

	if txerr != nil {
		return txerr
	}

	fmt.Printf("tx start : %v\n", tx)

	q := New(tx)

	err := fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if needRollback {
		return tx.Rollback()
	} else {
		return tx.Commit()
	}
}
