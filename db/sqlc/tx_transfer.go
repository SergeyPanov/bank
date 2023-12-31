package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"entry"`
	ToEntry     Entry    `json:"entry"`
}

var txKey = struct{}{}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: pgtype.Int8{
				Int64: arg.FromAccountId,
				Valid: true,
			},
			ToAccountID: pgtype.Int8{
				Int64: arg.ToAccountId,
				Valid: true,
			},
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: pgtype.Int8{
				Int64: arg.FromAccountId,
				Valid: true,
			},
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: pgtype.Int8{
				Int64: arg.ToAccountId,
				Valid: true,
			},
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, accId1 int64, amount1 int64, accId2 int64, amount2 int64) (acc1 Account, acc2 Account, err error) {
	acc1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	acc2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accId2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
