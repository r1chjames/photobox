package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

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

		service := NewShareService(mockRepo)
		share, err := service.CreateShare("album", "album-123", "user-1", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, share)
		assert.NotEmpty(t, share.Token)
		assert.Equal(t, "album", share.ResourceType)
		assert.Equal(t, "album-123", share.ResourceId)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockShareRepository)
			tt.mockSetup(mockRepo)

			service := NewShareService(mockRepo)
			result, err := service.GetSharedResource(tt.token, tt.password)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListShares tests listing shared links
func TestListShares(t *testing.T) {
	t.Run("successfully list shares", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		shares := []*domain.SharedLink{
			{Token: "token-1", ResourceType: "album"},
			{Token: "token-2", ResourceType: "photo"},
		}
		mockRepo.On("ListShares").Return(shares, nil)

		service := NewShareService(mockRepo)
		result, err := service.ListShares()

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

// TestRevokeShare tests revoking a shared link
func TestRevokeShare(t *testing.T) {
	t.Run("successfully revoke share", func(t *testing.T) {
		mockRepo := new(MockShareRepository)
		mockRepo.On("DeleteShare", "token-1").Return(nil)

		service := NewShareService(mockRepo)
		err := service.RevokeShare("token-1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
