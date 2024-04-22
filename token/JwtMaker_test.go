package token

import (
	"testing"
	"time"

	"github.com/MathHRM/simple_bank/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)



func TestJwtCreateToken(t *testing.T) {
	secreKey := util.RandomString(32)
	maker, err := NewJwtMaker(secreKey)
	require.NoError(t, err)

	user := util.RandomName()
	duration := time.Minute

	createdAt := time.Now()
	expiredAt := createdAt.Add(duration)

	token, err := maker.CreateToken(user, duration)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, user)
	require.WithinDuration(t, payload.CreatedAt, createdAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	secreKey := util.RandomString(32)
	maker, err := NewJwtMaker(secreKey)
	require.NoError(t, err)

	user := util.RandomName()
	duration := time.Minute

	token, err := maker.CreateToken(user, -duration)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ExpiredError.Error())
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	payload, err := NewPayload(util.RandomName(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, InvalidTokenError.Error())
	require.Nil(t, payload)
}