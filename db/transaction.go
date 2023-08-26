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

	q := New(tx)

	err := fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("execute err: %v, rb err: %v", err, rbErr)
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

func (transaction *Transaction) ExecTx(ctx context.Context, fn func(q *Queries) error, needRollback bool) error {
	tx, txerr := transaction.db.BeginTx(ctx, nil)

	if txerr != nil {
		return txerr
	}

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
