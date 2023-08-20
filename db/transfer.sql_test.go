package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTwoAccountForTestTransfer(t *testing.T) []Account {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	return []Account{account1, account2}
}

func createTransferFlow(t *testing.T, accounts []Account) {
	account1 := accounts[0]
	account2 := accounts[1]

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        account1.Balance,
	}

	transfer, tferr := testQuries.CreateTransfer(context.Background(), arg)
	require.NoError(t, tferr)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	entryFrom, entryFromErr := testQuries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account1.ID,
		Amount:    arg.Amount.Neg(),
	})

	require.NoError(t, entryFromErr)
	require.NotEmpty(t, entryFrom)
	require.Equal(t, entryFrom.AccountID, account1.ID)
	require.Equal(t, entryFrom.Amount, arg.Amount.Neg())

	entryTo, entryToErr := testQuries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account2.ID,
		Amount:    arg.Amount,
	})

	require.NoError(t, entryToErr)
	require.NotEmpty(t, entryTo)
	require.Equal(t, entryTo.AccountID, account2.ID)
	require.Equal(t, entryTo.Amount, arg.Amount)

	// TODO: update account1 and account2 balance
}

func Test_CreateTransfer(t *testing.T) {
	txerr := ExecTestingTx(context.Background(), testTx, func() error {
		accounts := createTwoAccountForTestTransfer(t)

		for i := 0; i < 10; i++ {
			go createTransferFlow(t, accounts)
		}

		return nil
	}, true)
	require.NoError(t, txerr)
}
