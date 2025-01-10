package auth

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
	duration     time.Duration
}

func New(symmetricKey string, duration time.Duration) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("SymmetricKey too short should be: %v", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		duration:     duration,
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(user *domain.User) (string, error) {
	payload, err := domain.NewPayload(*user, maker.duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) ParseToken(token string) (*domain.TokenPayload, error) {
	payload := &domain.TokenPayload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*domain.TokenPayload, error) {
	payload, err := maker.ParseToken(token)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
