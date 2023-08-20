package db

import (
	"context"
	"testing"

	"github.com/go-backend-practice/util"
	"github.com/stretchr/testify/require"
)

func createTwoAccountForTestTransfer(t *testing.T, q *Queries) []Account {
	account1 := CreateRandomAccount(t, q)
	account2 := CreateRandomAccount(t, q)

	return []Account{account1, account2}
}

func createTransferFlow(t *testing.T, q *Queries, accounts []Account) error {
	var transferErr error

	account1 := accounts[0]
	account2 := accounts[1]

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	transfer, tferr := q.CreateTransfer(context.Background(), arg)
	require.NoError(t, tferr)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	if !arg.Amount.Equal(transfer.Amount) {
		panic("Create transfer amount not equal!")
	}

	entryFrom, entryFromErr := q.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account1.ID,
		Amount:    arg.Amount.Neg(),
	})

	require.NoError(t, entryFromErr)
	require.NotEmpty(t, entryFrom)
	require.Equal(t, entryFrom.AccountID, account1.ID)
	if !entryFrom.Amount.Equal(arg.Amount.Neg()) {
		panic("Create EntryFrom amount not equal!")
	}
	entryTo, entryToErr := q.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account2.ID,
		Amount:    arg.Amount,
	})

	require.NoError(t, entryToErr)
	require.NotEmpty(t, entryTo)
	require.Equal(t, entryTo.AccountID, account2.ID)
	if !entryTo.Amount.Equal(arg.Amount) {
		panic("Create EntryTo amount not equal!")
	}

	// TODO: update account1 and account2 balance

	return transferErr
}

func Test_CreateTransfer(t *testing.T) {
	q := New(sharedConn)
	accounts := createTwoAccountForTestTransfer(t, q)

	errs := make(chan error)

	for i := 0; i < 10; i++ {
		go func() {
			tx := NewTransaction(sharedConn)
			err := tx.ExecTx(context.Background(), func(q *Queries) error {
				createTransforErr := createTransferFlow(t, q, accounts)
				return createTransforErr
			}, false)

			errs <- err
		}()
	}

	for i := 0; i < 10; i++ {
		err := <-errs
		require.NoError(t, err)
	}
}
