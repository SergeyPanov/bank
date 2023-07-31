package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/SergeyPanov/bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	acc := createRandomAccount(t)

	args := CreateEntryParams{
		AccountID: sql.NullInt64{
			Int64: acc.ID,
			Valid: true,
		},
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	createdEntry := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, createdEntry.ID, entry.ID)
	require.Equal(t, createdEntry.AccountID, entry.AccountID)
	require.Equal(t, createdEntry.Amount, entry.Amount)
	require.Equal(t, createdEntry.CreatedAt, entry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	createdEntry := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, createdEntry.ID, entry.ID)
	require.Equal(t, createdEntry.AccountID, entry.AccountID)
	require.Equal(t, createdEntry.Amount, entry.Amount)
	require.Equal(t, createdEntry.CreatedAt, entry.CreatedAt)
}
