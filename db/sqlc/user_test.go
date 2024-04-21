package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MathHRM/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	name := util.RandomName()
	lastName := util.RandomName()
	fullName := fmt.Sprintf("%s %s", name, lastName)
	hash, err := util.HashPassword( util.RandomString(8) )
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       name,
		FullName:       fullName,
		Email:          util.RandomEmail(name),
		HashedPassword: hash,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.NotEmpty(t, user.Username)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	getUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)

	require.Equal(t, user.Username, getUser.Username)
	require.Equal(t, user.FullName, getUser.FullName)
	require.Equal(t, user.Email, getUser.Email)
	require.WithinDuration(t, user.CreatedAt, getUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, getUser.PasswordChangedAt, time.Second)
}