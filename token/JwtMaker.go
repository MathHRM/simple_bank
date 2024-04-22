package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)



const MIN_SECRET_KEY_LENGTH = 32

type JwtMaker struct {
	SecretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if( len(secretKey) < MIN_SECRET_KEY_LENGTH ) {
		return nil, fmt.Errorf("Tamanho invalido da key")
	}

	return &JwtMaker{secretKey}, nil
}

func (jwtMaker *JwtMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString( []byte(jwtMaker.SecretKey) )
}

func (jwtMaker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidTokenError
		}
		return []byte(jwtMaker.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		varErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(varErr.Inner, ExpiredError) {
			return nil, ExpiredError
		}
		return nil, InvalidTokenError
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, InvalidTokenError
	}

	return payload, nil
}