package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromE := res.FromEntry
		require.NotEmpty(t, fromE)
		require.Equal(t, account1.ID, fromE.AccountID.Int64)
		require.Equal(t, -amount, fromE.Amount)
		require.NotZero(t, fromE.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromE.ID)
		require.NoError(t, err)

		toE := res.ToEntry
		require.NotEmpty(t, fromE)
		require.Equal(t, account2.ID, toE.AccountID.Int64)
		require.Equal(t, amount, toE.Amount)
		require.NotZero(t, toE.CreatedAt)

		// TODO: check acc balance

		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)

		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n)
	}

	updatedAcc1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAcc2.Balance)

}
