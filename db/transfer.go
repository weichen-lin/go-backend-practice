package db

import (
	"context"

	"github.com/shopspring/decimal"
)

type TransferError struct {
	Msg string
}

func (e *TransferError) Error() string {
	return e.Msg
}

func (q *Queries) AccountTransfer(accoun1 Account, account2 Account, amount decimal.Decimal) (FromAccount UpdateAccountRow, ToAccount UpdateAccountRow, err error) {
	// Always check that account update would happen in the same order
	if accoun1.ID.String() > account2.ID.String() {
		FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:     accoun1.ID,
			Amount: amount.Neg(),
		})

		if err != nil {
			return
		}

		ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:     account2.ID,
			Amount: amount,
		})

		if err != nil {
			return
		}
	} else {
		ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:     account2.ID,
			Amount: amount,
		})

		if err != nil {
			return
		}
		FromAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
			ID:     accoun1.ID,
			Amount: amount.Neg(),
		})

		if err != nil {
			return
		}
	}

	return FromAccount, ToAccount, nil
}
