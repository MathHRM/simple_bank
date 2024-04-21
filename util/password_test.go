package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)



func TestHashPassword(t *testing.T) {
	password := RandomString(8)
	wrongPassword := RandomString(8)

	hash1, err := HashPassword( password )
	require.NoError(t, err)
	err = CompareHashPassword(password, hash1)
	require.NoError(t, err)

	hash2, err := HashPassword( password )
	require.NoError(t, err)
	err = CompareHashPassword(password, hash2)
	require.NoError(t, err)

	require.NotEqual(t, hash1, hash2)
	
	err = CompareHashPassword(wrongPassword, hash1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}