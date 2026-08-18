package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/auth"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

func strPtr(s string) *string { return &s }

// MockShareRepository is a mock implementation of port.ShareRepository
type MockShareRepository struct {
	mock.Mock
}

func (m *MockShareRepository) CreateShare(share *domain.SharedLink) error {
	args := m.Called(share)
	return args.Error(0)
}

func (m *MockShareRepository) GetShareByToken(token string) (*domain.SharedLink, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SharedLink), args.Error(1)
}

func (m *MockShareRepository) ListSharesByOwner(createdBy string) ([]*domain.SharedLink, error) {
	args := m.Called(createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.SharedLink), args.Error(1)
}

func (m *MockShareRepository) ListShares() ([]*domain.SharedLink, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.SharedLink), args.Error(1)
}

func (m *MockShareRepository) DeleteShare(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockShareRepository) IncrementViewCount(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

// TestCreateShare tests creating a shared link
func TestCreateShare(t *testing.T) {
	t.Run("successfully create share", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		mockRepo.On("CreateShare", mock.AnythingOfType("*domain.SharedLink")).Return(nil)

		service := NewShareService(mockRepo, nil, nil)
		share, err := service.CreateShare("album", "album-123", "user-1", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, share)
		assert.NotEmpty(t, share.Token)
		assert.Equal(t, "album", share.ResourceType)
		assert.Equal(t, "album-123", share.ResourceId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("share password is hashed at rest", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		var stored *domain.SharedLink
		mockRepo.On("CreateShare", mock.AnythingOfType("*domain.SharedLink")).Return(nil).Run(func(args mock.Arguments) {
			stored = args.Get(0).(*domain.SharedLink)
		})

		service := NewShareService(mockRepo, nil, nil)
		password := "secret123"
		_, err := service.CreateShare("photo", "photo-1", "user-1", nil, &password)

		assert.NoError(t, err)
		assert.NotNil(t, stored)
		assert.NotEqual(t, "secret123", stored.PasswordHash, "password must not be stored plaintext")
		assert.Contains(t, stored.PasswordHash, "$argon2id$", "password should be an argon2id hash")
		mockRepo.AssertExpectations(t)
	})
}

// TestGetSharedResource tests retrieving a shared resource
func TestGetSharedResource(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		password  *string
		mockSetup func(*MockShareRepository)
		validate  func(*testing.T, *domain.SharedLink, error)
	}{
		{
			name:  "successfully get shared resource",
			token: "valid-token",
			mockSetup: func(m *MockShareRepository) {
				m.On("GetShareByToken", "valid-token").Return(&domain.SharedLink{
					Token:        "valid-token",
					ResourceType: "album",
					ResourceId:   "album-123",
				}, nil)
				m.On("IncrementViewCount", "valid-token").Return(nil)
			},
			validate: func(t *testing.T, link *domain.SharedLink, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, link)
				assert.Equal(t, "album-123", link.ResourceId)
			},
		},
		{
			name:  "expired shared link",
			token: "expired-token",
			mockSetup: func(m *MockShareRepository) {
				expired := time.Now().Add(-time.Hour)
				m.On("GetShareByToken", "expired-token").Return(&domain.SharedLink{
					Token:        "expired-token",
					ResourceType: "album",
					ResourceId:   "album-123",
					Expiry:       &expired,
				}, nil)
			},
			validate: func(t *testing.T, link *domain.SharedLink, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrSharedLinkExpired, err)
				assert.Nil(t, link)
			},
		},
		{
			name:     "password required",
			token:    "password-token",
			password: nil,
			mockSetup: func(m *MockShareRepository) {
				m.On("GetShareByToken", "password-token").Return(&domain.SharedLink{
					Token:        "password-token",
					ResourceType: "photo",
					ResourceId:   "photo-123",
					PasswordHash: "secret",
				}, nil)
			},
			validate: func(t *testing.T, link *domain.SharedLink, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrSharedLinkPasswordRequired, err)
				assert.Nil(t, link)
			},
		},
		{
			name:  "not found",
			token: "missing-token",
			mockSetup: func(m *MockShareRepository) {
				m.On("GetShareByToken", "missing-token").Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, link *domain.SharedLink, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrDataNotFound, err)
				assert.Nil(t, link)
			},
		},
		{
			name:     "correct hashed password grants access",
			token:    "hashed-token",
			password: strPtr("secret123"),
			mockSetup: func(m *MockShareRepository) {
				hash, _ := auth.CreateHash("secret123", auth.DefaultArgon2idHash())
				m.On("GetShareByToken", "hashed-token").Return(&domain.SharedLink{
					Token:        "hashed-token",
					ResourceType: "photo",
					ResourceId:   "photo-123",
					PasswordHash: hash,
				}, nil)
				m.On("IncrementViewCount", "hashed-token").Return(nil)
			},
			validate: func(t *testing.T, link *domain.SharedLink, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, link)
			},
		},
		{
			name:     "wrong password denied",
			token:    "wrong-token",
			password: strPtr("wrongpass"),
			mockSetup: func(m *MockShareRepository) {
				hash, _ := auth.CreateHash("secret123", auth.DefaultArgon2idHash())
				m.On("GetShareByToken", "wrong-token").Return(&domain.SharedLink{
					Token:        "wrong-token",
					ResourceType: "photo",
					ResourceId:   "photo-123",
					PasswordHash: hash,
				}, nil)
			},
			validate: func(t *testing.T, link *domain.SharedLink, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrSharedLinkPasswordRequired, err)
				assert.Nil(t, link)
			},
		},
		{
			name:     "legacy plaintext password still verifies",
			token:    "legacy-token",
			password: strPtr("oldsecret"),
			mockSetup: func(m *MockShareRepository) {
				m.On("GetShareByToken", "legacy-token").Return(&domain.SharedLink{
					Token:        "legacy-token",
					ResourceType: "photo",
					ResourceId:   "photo-123",
					PasswordHash: "oldsecret", // pre-hashing plaintext row
				}, nil)
				m.On("IncrementViewCount", "legacy-token").Return(nil)
			},
			validate: func(t *testing.T, link *domain.SharedLink, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, link)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockShareRepository)
			tt.mockSetup(mockRepo)

			service := NewShareService(mockRepo, nil, nil)
			result, err := service.GetSharedResource(tt.token, tt.password)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListShares tests listing shared links
func TestListShares(t *testing.T) {
	t.Run("successfully list all shares (admin)", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		shares := []*domain.SharedLink{
			{Token: "token-1", ResourceType: "album"},
			{Token: "token-2", ResourceType: "photo"},
		}
		mockRepo.On("ListShares").Return(shares, nil)

		service := NewShareService(mockRepo, nil, nil)
		result, err := service.ListShares("")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("list only the requesting user's shares", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		shares := []*domain.SharedLink{{Token: "token-1", CreatedBy: "user-1"}}
		mockRepo.On("ListSharesByOwner", "user-1").Return(shares, nil)

		service := NewShareService(mockRepo, nil, nil)
		result, err := service.ListShares("user-1")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].CreatedBy)
		mockRepo.AssertExpectations(t)
	})
}

// TestRevokeShare tests revoking a shared link
func TestRevokeShare(t *testing.T) {
	t.Run("successfully revoke own share", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		mockRepo.On("GetShareByToken", "token-1").Return(&domain.SharedLink{Token: "token-1", CreatedBy: "user-1"}, nil)
		mockRepo.On("DeleteShare", "token-1").Return(nil)

		service := NewShareService(mockRepo, nil, nil)
		err := service.RevokeShare("token-1", "user-1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("cannot revoke another user's share", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		mockRepo.On("GetShareByToken", "token-1").Return(&domain.SharedLink{Token: "token-1", CreatedBy: "user-2"}, nil)

		service := NewShareService(mockRepo, nil, nil)
		err := service.RevokeShare("token-1", "user-1")

		assert.ErrorIs(t, err, domain.ErrDataNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("admin can revoke any share", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		mockRepo.On("DeleteShare", "token-1").Return(nil)

		service := NewShareService(mockRepo, nil, nil)
		err := service.RevokeShare("token-1", "")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
