package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	wordLength          = 6
	passwordLength = 36
	emailLength         = 10
	emailDomainLength   = 5
)

func GenerateUsername() string {
	return randomEntity.RandomString(wordLength)
}

func GenerateFullname() string {
	return fmt.Sprintf("%s %s %s", randomEntity.RandomString(wordLength), randomEntity.RandomString(wordLength), randomEntity.RandomString(wordLength))
}

func GenerateEmail() string {
	return fmt.Sprintf("%s@%s.com", randomEntity.RandomString(emailLength), randomEntity.RandomString(emailDomainLength))
}

func GeneratePassword() string {
	return randomEntity.RandomString(passwordLength)
}

func createRandomUser(t *testing.T) User {
	password := GeneratePassword()
	hashedPassword, err := passwordManager.HashPassword(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       GenerateUsername(),
		FullName:       GenerateFullname(),
		Email:          GenerateEmail(),
		HashedPassword: hashedPassword,
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	err = passwordManager.CheckPassword(password, hashedPassword)
	require.NoError(t, err)
	
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByUsername(t *testing.T) {
	createdUser := createRandomUser(t)

	username := createdUser.Username
	user, err := testQueries.GetUserByUsername(context.Background(), username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createdUser.Username, user.Username)
	require.Equal(t, createdUser.FullName, user.FullName)
	require.Equal(t, createdUser.Email, user.Email)
	require.Equal(t, createdUser.HashedPassword, user.HashedPassword)

	require.WithinDuration(t, createdUser.CreatedAt, user.CreatedAt, time.Second)
	require.WithinDuration(t, createdUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
}
