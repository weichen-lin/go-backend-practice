package db

import (
	"context"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func createTwoAccountForTestTransfer(t *testing.T, q *Queries) []Account {
	account1 := CreateRandomAccount(t, q)
	account2 := CreateRandomAccount(t, q)

	return []Account{account1, account2}
}

func createTransferFlow(t *testing.T, q *Queries, accounts []Account) (decimal.Decimal, error) {
	var transferErr error

	account1 := accounts[0]
	account2 := accounts[1]

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        decimal.NewFromFloat(123.456),
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

	getAccount1, err1 := q.GetAccountForUpdate(context.Background(), accounts[0].ID)
	getAccount2, err2 := q.GetAccountForUpdate(context.Background(), accounts[1].ID)

	require.NoError(t, err1)
	require.NoError(t, err2)

	FromAccount, ToAccount, updateErr := q.AccountTransfer(getAccount1, getAccount2, arg.Amount)

	require.NotEmpty(t, FromAccount)
	require.NotEmpty(t, ToAccount)
	require.NoError(t, updateErr)

	if !FromAccount.Balance.Equal(getAccount1.Balance.Add(arg.Amount.Neg())) {
		panic("Update Account1 balance not equal!")
	}

	if !ToAccount.Balance.Equal(getAccount2.Balance.Add(arg.Amount)) {
		panic("Update Account2 balance not equal!")
	}

	return arg.Amount, transferErr
}

func Test_CreateTransfer(t *testing.T) {

	wg := sync.WaitGroup{}

	q := New(sharedConn)
	accounts := createTwoAccountForTestTransfer(t, q)
	require.Len(t, accounts, 2)

	repeat := 10

	errs := make(chan error, repeat)
	amounts := make(chan decimal.Decimal, repeat)

	for i := 0; i < repeat; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tx := NewTransaction(sharedConn)
			err := tx.ExecTx(context.Background(), func(q *Queries) error {
				amountTransfer, createTransforErr := createTransferFlow(t, q, accounts)
				amounts <- amountTransfer
				return createTransforErr
			}, false)

			errs <- err
		}()
	}

	wg.Wait()

	close(errs)
	close(amounts)

	var finalAmount decimal.Decimal

	for amount := range amounts {
		finalAmount = finalAmount.Add(amount)
	}

	err := <-errs
	require.NoError(t, err)

	check1, check1Err := q.GetAccountForUpdate(context.Background(), accounts[0].ID)
	check2, check2Err := q.GetAccountForUpdate(context.Background(), accounts[1].ID)

	require.NoError(t, check1Err)
	require.NoError(t, check2Err)

	if !check1.Balance.Equal(accounts[0].Balance.Add(finalAmount.Neg())) {
		panic("Final Account1 balance not correct!")
	}

	if !check2.Balance.Equal(accounts[1].Balance.Add(finalAmount)) {
		panic("Final Account2 balance not correct!")
	}

	deleteTransferErr := q.DeleteTransfer(context.Background(), accounts[0].ID)

	require.NoError(t, deleteTransferErr)

	deleteEntry1Err := q.DeleteEntry(context.Background(), accounts[0].ID)
	deleteEntry2Err := q.DeleteEntry(context.Background(), accounts[1].ID)

	require.NoError(t, deleteEntry1Err)
	require.NoError(t, deleteEntry2Err)

	deleteAccount1Err := q.DeleteAccount(context.Background(), accounts[0].ID)
	deleteAccount2Err := q.DeleteAccount(context.Background(), accounts[1].ID)

	require.NoError(t, deleteAccount1Err)
	require.NoError(t, deleteAccount2Err)
}
