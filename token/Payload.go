package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ExpiredError      = errors.New("Token expirado")
	InvalidTokenError = errors.New("Token invalido")
)

type Payload struct {
	ID        uuid.UUID `json:id`
	Username  string    `json:username`
	CreatedAt time.Time `json:created_at`
	ExpiredAt time.Time `json:expired_at`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ExpiredError
	}
	return nil
}
