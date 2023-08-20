package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateEntry(t *testing.T) {
	var testCreateEntryError error

	txerr := ExecTestingTx(context.Background(), testTx, func() error {
		account1 := CreateRandomAccount(t)

		arg := CreateEntryParams{
			AccountID: account1.ID,
			Amount:    account1.Balance,
		}

		entry1, err := testQuries.CreateEntry(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, entry1)

		require.Equal(t, arg.AccountID, entry1.AccountID)

		require.NotEmpty(t, entry1.ID)
		require.NotEmpty(t, entry1.CreatedAt)

		if !arg.Amount.Equal(entry1.Amount) {
			panic("Create entry amount not equal!")
		}
		return testCreateEntryError
	}, true)

	require.NoError(t, txerr)
}
