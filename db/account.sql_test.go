package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-backend-practice/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T, q *Queries) Account {
	arg := CreateAccountParams{
		Owner:    "test_" + util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := q.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotEmpty(t, account.ID)

	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.LastModifiedAt)

	if !arg.Balance.Equal(account.Balance) {
		panic("Create random account balance not equal!")
	}

	return account
}

func Test_GetAccount(t *testing.T) {
	txerr := ExecTestingTx(context.Background(), testTx, func(q *Queries) error {
		var getAccountError error
		account1 := CreateRandomAccount(t, q)
		account2, err := q.GetAccount(context.Background(), account1.ID)

		require.NoError(t, err)
		require.NotEmpty(t, account2)

		require.Equal(t, account1.ID, account2.ID)

		if !account1.Balance.Equal(account2.Balance) {
			panic("Create new account error!")
		}
		require.Equal(t, account1.Currency, account2.Currency)
		require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

		return getAccountError
	}, true)

	require.NoError(t, txerr)
}

func Test_UpdateAccount(t *testing.T) {
	var updateAccount error
	txerr := ExecTestingTx(context.Background(), testTx, func(q *Queries) error {
		account1 := CreateRandomAccount(t, q)

		arg := UpdateAccountParams{
			ID:      account1.ID,
			Balance: account1.Balance.Add(util.RandomBalance()).Round(3),
		}

		account2, err := q.UpdateAccount(context.Background(), arg)

		require.NoError(t, err)
		require.NotEmpty(t, account2)

		require.Equal(t, account1.ID, account2.ID)

		if !account2.Balance.Equal(arg.Balance) {
			panic("Update balance not equal!")
		}
		return updateAccount
	}, true)

	require.NoError(t, txerr)
}

func Test_DeleteAccount(t *testing.T) {
	var testDeleteAccount error
	txerr := ExecTestingTx(context.Background(), testTx, func(q *Queries) error {
		account1 := CreateRandomAccount(t, q)
		err := q.DeleteAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		account2, err := q.GetAccount(context.Background(), account1.ID)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, account2)

		return testDeleteAccount
	}, true)

	require.NoError(t, txerr)
}

func Test_ListAccount(t *testing.T) {
	var testListAccount error
	txerr := ExecTestingTx(context.Background(), testTx, func(q *Queries) error {
		for i := 0; i < 10; i++ {
			CreateRandomAccount(t, q)
		}

		arg := ListAccountsParams{
			Limit:  5,
			Offset: 5,
		}

		accounts, err := q.ListAccounts(context.Background(), arg)
		require.NoError(t, err)
		require.Len(t, accounts, 5)

		for _, account := range accounts {
			require.NotEmpty(t, account)
		}
		return testListAccount
	}, true)

	require.NoError(t, txerr)
}
