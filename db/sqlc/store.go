package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface{
	TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResults, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// This func execute all transaction
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
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
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResults struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Create a transfer, 2 entries (from and to), update 2 accounts (from and to)
// nolint:gosimple
func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTxResults, error) {
	var result TransferTxResults

	err := s.execTx(ctx, func(q *Queries) error {
		// Verify amount transfer > 0
		if arg.Amount <= 0 {
			return fmt.Errorf("amount should be greater than 0, got %d", arg.Amount)
		}

		// Create transfer
		transfer, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.Transfer = transfer

		// Create 2 entries
		fromEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry = fromEntry

		toEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry = toEntry

		// Update 2 accounts
		// advoid deadlock: update account with smaller id first
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = adjustBalance(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = adjustBalance(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func adjustBalance(
	ctx context.Context,
	q *Queries,
	account1ID, amount1, account2ID, amount2 int64,
) (account1, account2 Account, err error) {
	checkAccountID := account1ID
	if amount2 < 0 {
		checkAccountID = account2ID
	}
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	// verify balance after transfer > 0
	checkFromAccount, err := q.GetAccountByID(ctx, checkAccountID)
	if err != nil {
		return
	}
	if checkFromAccount.Balance < 0 {
		err = fmt.Errorf("the amount transfer is greater than its balance")
		return 
	}
	return
}
