package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	minBalance int64 = 1000
	maxBalance int64 = 10000
)

func GenerateAccountBalance() int64 {
	return randomEntity.RandomInt(minBalance, maxBalance)
}

func GenerateCurrency() string {
	return "USD"
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: GenerateAccountBalance(),
		Currency: GenerateCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)
	
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)

	return acc
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountByID(t *testing.T) {
	createdAccount := createRandomAccount(t)

	id := createdAccount.ID
	acc, err := testQueries.GetAccountByID(context.Background(), id)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, createdAccount.Owner, acc.Owner)
	require.Equal(t, createdAccount.Balance, acc.Balance)
	require.Equal(t, createdAccount.Currency, acc.Currency)
	require.WithinDuration(t, createdAccount.CreatedAt, acc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	arg := UpdateAccountByIDParams{
		ID: createdAccount.ID,
		Balance: GenerateAccountBalance(),
	}

	updatedAccount, err := testQueries.UpdateAccountByID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccount.Currency)
	require.WithinDuration(t, createdAccount.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	createAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), createAccount.ID)
	require.NoError(t, err)

	acc, err := testQueries.GetAccountByID(context.Background(), createAccount.ID)
	require.Error(t, err)
	require.Empty(t, acc)
}