package port

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// ApiKeyRepository persists API keys (hashed at rest).
type ApiKeyRepository interface {
	Create(key *domain.ApiKey) error
	GetByHash(hash string) (*domain.ApiKey, error)
	List() ([]*domain.ApiKey, error)
	Revoke(id string) error
	TouchLastUsed(id string) error
}

// ApiKeyService manages programmatic API keys.
type ApiKeyService interface {
	CreateKey(name string, scope domain.ApiKeyScope, createdBy string) (string, *domain.ApiKey, error)
	ValidateKey(key string) (*domain.ApiKey, error)
	ListKeys() ([]*domain.ApiKey, error)
	RevokeKey(id string) error
}
