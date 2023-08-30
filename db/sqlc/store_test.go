package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	r := require.New(t)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	amount := randomEntity.RandomInt(1, 10)
	numTransaction := 2

	errs := make(chan error)
	results := make(chan *TransferTxResults)
	for i := 0; i < numTransaction; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < numTransaction; i++ {
		err := <-errs
		r.NoError(err)

		result := <-results
		r.NotEmpty(result)

		// Check transfer
		transfer := result.Transfer
		r.NotEmpty(transfer)
		r.Equal(account1.ID, transfer.FromAccountID)
		r.Equal(account2.ID, transfer.ToAccountID)
		r.Equal(amount, transfer.Amount)
		r.NotZero(transfer.ID)
		r.NotZero(transfer.CreatedAt)
		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		r.NoError(err)

		// Check entries
		fromEntry := result.FromEntry
		r.NotEmpty(fromEntry)
		r.Equal(account1.ID, fromEntry.AccountID)
		r.Equal(amount, -fromEntry.Amount)
		r.NotZero(fromEntry.CreatedAt)
		r.NotZero(fromEntry.ID)
		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		r.NoError(err)

		toEntry := result.ToEntry
		r.NotEmpty(toEntry)
		r.Equal(account2.ID, toEntry.AccountID)
		r.Equal(amount, toEntry.Amount)
		r.NotZero(toEntry.CreatedAt)
		r.NotZero(toEntry.ID)
		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		r.NoError(err)

		// Check accounts
		fromAccount := result.FromAccount
		diff1 := account1.Balance - fromAccount.Balance
		r.NotEmpty(fromAccount)
		r.Equal(account1.ID, fromAccount.ID)
		r.True(diff1 > 0)
		r.True(diff1%amount==0)
		
		toAccount := result.ToAccount
		diff2 := toAccount.Balance - account2.Balance
		r.NotEmpty(toAccount)
		r.Equal(account2.ID, toAccount.ID)
		r.Equal(diff1, diff2)
	}

	updatedAccount1, err := testStore.GetAccountByID(context.Background(), account1.ID)
	r.NoError(err)
	r.Equal(account1.Balance - amount*int64(numTransaction), updatedAccount1.Balance)

	updatedAccount2, err := testStore.GetAccountByID(context.Background(), account2.ID)
	r.NoError(err)
	r.Equal(account2.Balance + amount*int64(numTransaction), updatedAccount2.Balance)
}

// Test deadlock: 10 transaction, 5 from 1 to 2, 5 from 2 to 1 -> balances remain the same
func TestTransferTxDeadlock(t *testing.T) {
	r := require.New(t)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	amount := randomEntity.RandomInt(1, 10)
	numTransaction := 2

	errs := make(chan error)
	for i := 0; i < numTransaction; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2==1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < numTransaction; i++ {
		err := <-errs
		r.NoError(err)
	}

	updatedAccount1, err := testStore.GetAccountByID(context.Background(), account1.ID)
	r.NoError(err)
	r.Equal(updatedAccount1.Balance, account1.Balance)
	
	updatedAccount2, err := testStore.GetAccountByID(context.Background(), account2.ID)
	r.NoError(err)
	r.Equal(updatedAccount2.Balance, account2.Balance)
}