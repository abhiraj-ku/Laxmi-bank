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
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTX executes a function within a database transactions
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v , rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// this TransferTXParams contains all necessary details in order to do a transfer of money
type TransferTXParams struct {
	FromAccountID int64 `json:"from_accoutn_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTXResult is the result of TransferTX/  transfer transaction
type TransferTXResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// transferTX performs a money transfer from one account to another
// it creates a transfer record , add acc entries, and update acc balance with single db call/transaction
func (store *Store) TransferTX(ctx context.Context, arg TransferTXParams) (TransferTXResult, error) {
	var result TransferTXResult

	err := store.execTX(ctx, func(q *Queries) error {
		var err error

		// create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// substract money from 'Fromaccount/senders account'
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Add money to 'ToAccount/reciever's account'
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: create update account and implement locking mechanism to handle deadlock

		// transfer/move money from account 1
		account1, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return nil
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		// add money to account 2
		account2, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return nil
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
