package test_db

import (
	"context"
	"testing"

	"github.com/go-backend-practice/db"
	"github.com/stretchr/testify/require"
)

func Test_CreateEntry(t *testing.T) {
	var testCreateEntryError error

	txerr := db.ExecTestingTx(context.Background(), testTx, func(q *db.Queries) error {
		account1 := CreateRandomAccount(t, q)

		arg := db.CreateEntryParams{
			AccountID: account1.ID,
			Amount:    account1.Balance,
		}

		entry1, err := q.CreateEntry(context.Background(), arg)

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
