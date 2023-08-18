package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-backend-practice/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    "test_" + util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQuries.CreateAccount(context.Background(), arg)

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
	account1 := CreateRandomAccount(t)
	account2, err := testQuries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)

	if !account1.Balance.Equal(account2.Balance) {
		panic("Create new account error!")
	}
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func Test_UpdateAccount(t *testing.T) {

	account1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: account1.Balance.Add(util.RandomBalance()).Round(3),
	}

	account2, err := testQuries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)

	if !account2.Balance.Equal(arg.Balance) {
		panic("Update balance not equal!")
	}
}

func Test_DeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQuries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQuries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func Test_ListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQuries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
