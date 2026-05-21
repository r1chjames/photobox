package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiryAt  time.Time `json:"expiry_at"`
}

func NewPayload(user User, duration time.Duration) (*TokenPayload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &TokenPayload{
		ID:        tokenId,
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: time.Now(),
		ExpiryAt:  time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiryAt) {
		return errors.New("token has expired")
	}
	return nil
}
