package db

import (
	"context"
	"testing"
	"time"

	"github.com/SergeyPanov/bank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:   user.Username,
		Balance: util.RandomMoney(),
		Currency: pgtype.Text{
			String: util.RandomCurrency(),
			Valid:  true,
		},
	}

	account, err := testStore.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc, err := testStore.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, acc1.ID, acc.ID)
	require.Equal(t, acc1.Balance, acc.Balance)
	require.Equal(t, acc1.Currency, acc.Currency)
	require.Equal(t, acc1.Owner, acc.Owner)
	require.WithinDuration(t, acc1.CreatedAt, acc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc2.Balance, arg.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	err := testStore.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, acc)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
		require.Equal(t, lastAccount.Owner, acc.Owner)
	}
}
