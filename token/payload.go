package token

import (
	"errors"
	"time"

	uuid "github.com/google/uuid"
)

// Different types of errors that can be returned by VerifyToken 
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// contains the payload data of the token
type Payload struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expires_at"`
}

// creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// checks for expired tokens
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}