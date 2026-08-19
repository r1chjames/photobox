package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockApiKeyRepository is a mock implementation of port.ApiKeyRepository.
type MockApiKeyRepository struct {
	mock.Mock
}

func (m *MockApiKeyRepository) Create(key *domain.ApiKey) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockApiKeyRepository) GetByHash(hash string) (*domain.ApiKey, error) {
	args := m.Called(hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ApiKey), args.Error(1)
}

func (m *MockApiKeyRepository) List() ([]*domain.ApiKey, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ApiKey), args.Error(1)
}

func (m *MockApiKeyRepository) Revoke(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockApiKeyRepository) TouchLastUsed(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestApiKeyService(t *testing.T) {
	t.Run("create key returns plaintext once and stores hash", func(t *testing.T) {
		mockRepo := new(MockApiKeyRepository)
		var stored *domain.ApiKey
		mockRepo.On("Create", mock.AnythingOfType("*domain.ApiKey")).Return(nil).Run(func(args mock.Arguments) {
			stored = args.Get(0).(*domain.ApiKey)
		})

		svc := NewApiKeyService(mockRepo)
		key, apiKey, err := svc.CreateKey("test-key", domain.ApiKeyReadOnly, "user-1")

		assert.NoError(t, err)
		assert.NotEmpty(t, key)
		assert.True(t, len(key) > 10, "key should be a long random string")
		assert.NotEqual(t, key, stored.KeyHash, "plaintext must not equal stored hash")
		assert.Equal(t, domain.ApiKeyReadOnly, apiKey.Scope)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validate key matches stored hash", func(t *testing.T) {
		mockRepo := new(MockApiKeyRepository)
		svc := NewApiKeyService(mockRepo)

		var created *domain.ApiKey
		mockRepo.On("Create", mock.AnythingOfType("*domain.ApiKey")).Return(nil).Run(func(args mock.Arguments) {
			created = args.Get(0).(*domain.ApiKey)
		})
		key, apiKey, err := svc.CreateKey("test-key", domain.ApiKeyReadWrite, "user-1")
		assert.NoError(t, err)
		_ = created

		mockRepo.On("GetByHash", hashApiKey(key)).Return(apiKey, nil)
		mockRepo.On("TouchLastUsed", apiKey.ID).Return(nil)

		validated, err := svc.ValidateKey(key)
		assert.NoError(t, err)
		assert.Equal(t, apiKey.ID, validated.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid key is unauthorized", func(t *testing.T) {
		mockRepo := new(MockApiKeyRepository)
		mockRepo.On("GetByHash", mock.Anything).Return(nil, domain.ErrDataNotFound)

		svc := NewApiKeyService(mockRepo)
		_, err := svc.ValidateKey("bad-key")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
		mockRepo.AssertExpectations(t)
	})
}
