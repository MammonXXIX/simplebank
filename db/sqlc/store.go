package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store {
		Queries: New(db),
		db: db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx,nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil { return err }

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil { return err }

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil { return err }

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = updateBalance(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = updateBalance(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)	
		}

		return nil
	})

	return result, err
}

func updateBalance(ctx context.Context, q *Queries, fromAccountID, fromAccountBalance, toAccountID, toAccountBalance int64) (fromAccount Account, toAccount Account, err error) {
	fromAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID: fromAccountID,
		Balance: fromAccountBalance,
	})
	if err != nil { return }

	toAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID: toAccountID,
		Balance: toAccountBalance,
	})
	if err != nil { return }

	return
}
