package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type ApiKeyService struct {
	repo port.ApiKeyRepository
}

func NewApiKeyService(repo port.ApiKeyRepository) *ApiKeyService {
	return &ApiKeyService{repo}
}

// generateKey creates a random API key and returns it plus its hash.
func generateKey() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	key := "pb_" + hex.EncodeToString(bytes)
	sum := sha256.Sum256([]byte(key))
	return key, hex.EncodeToString(sum[:]), nil
}

// CreateKey generates a new API key, stores its hash, and returns the
// plaintext key (shown only once).
func (s *ApiKeyService) CreateKey(name string, scope domain.ApiKeyScope, createdBy string) (string, *domain.ApiKey, error) {
	key, hash, err := generateKey()
	if err != nil {
		return "", nil, err
	}
	apiKey := &domain.ApiKey{
		ID:        uuid.New().String(),
		Name:      name,
		KeyHash:   hash,
		Scope:     scope,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(apiKey); err != nil {
		return "", nil, err
	}
	return key, apiKey, nil
}

// ValidateKey checks a plaintext key against the stored hash and returns the
// key record (or nil if invalid/revoked).
func (s *ApiKeyService) ValidateKey(key string) (*domain.ApiKey, error) {
	sum := sha256.Sum256([]byte(key))
	hash := hex.EncodeToString(sum[:])
	apiKey, err := s.repo.GetByHash(hash)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	_ = s.repo.TouchLastUsed(apiKey.ID)
	return apiKey, nil
}

func (s *ApiKeyService) ListKeys() ([]*domain.ApiKey, error) {
	return s.repo.List()
}

func (s *ApiKeyService) RevokeKey(id string) error {
	return s.repo.Revoke(id)
}

// hashApiKey is a helper for tests / external use.
func hashApiKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
