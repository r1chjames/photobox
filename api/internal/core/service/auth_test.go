package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/auth"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockTokenService is a mock implementation of port.TokenService
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) CreateToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) VerifyToken(token string) (*domain.TokenPayload, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenPayload), args.Error(1)
}

// Helper function to create a valid password hash for testing
func createTestPasswordHash(password string) string {
	hash, _ := auth.CreateHash(password, &auth.Argon2idHash{
		Iterations:  1,
		Memory:      64 * 1024,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})
	return hash
}

// TestLogin tests the authentication login functionality
func TestLogin(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		mockSetup     func(*MockUserRepository, *MockTokenService)
		expectedError error
		validate      func(*testing.T, string, error)
	}{
		{
			name:     "successful login",
			username: "testuser",
			password: "correctPassword123",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				passwordHash := createTestPasswordHash("correctPassword123")
				userRepo.On("GetUserByUsername", "testuser").Return(&domain.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: passwordHash,
				}, nil)
				tokenSvc.On("CreateToken", mock.MatchedBy(func(u *domain.User) bool {
					return u.ID == "user123" && u.Username == "testuser"
				})).Return("valid.jwt.token", nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "valid.jwt.token", token)
			},
		},
		{
			name:     "user not found",
			username: "nonexistent",
			password: "anyPassword",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				userRepo.On("GetUserByUsername", "nonexistent").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrInvalidCredentials,
			validate: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInvalidCredentials, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "incorrect password",
			username: "testuser",
			password: "wrongPassword",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				correctPasswordHash := createTestPasswordHash("correctPassword123")
				userRepo.On("GetUserByUsername", "testuser").Return(&domain.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: correctPasswordHash,
				}, nil)
			},
			expectedError: domain.ErrInvalidCredentials,
			validate: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInvalidCredentials, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "repository internal error",
			username: "testuser",
			password: "anyPassword",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				userRepo.On("GetUserByUsername", "testuser").Return(nil, domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "token creation fails",
			username: "testuser",
			password: "correctPassword123",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				passwordHash := createTestPasswordHash("correctPassword123")
				userRepo.On("GetUserByUsername", "testuser").Return(&domain.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: passwordHash,
				}, nil)
				tokenSvc.On("CreateToken", mock.Anything).Return("", domain.ErrTokenCreation)
			},
			expectedError: domain.ErrTokenCreation,
			validate: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrTokenCreation, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				userRepo.On("GetUserByUsername", "").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrInvalidCredentials,
			validate: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInvalidCredentials, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "empty password",
			username: "testuser",
			password: "",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				correctPasswordHash := createTestPasswordHash("correctPassword123")
				userRepo.On("GetUserByUsername", "testuser").Return(&domain.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: correctPasswordHash,
				}, nil)
			},
			expectedError: domain.ErrInvalidCredentials,
			validate: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInvalidCredentials, err)
				assert.Empty(t, token)
			},
		},
		{
			name:     "password with special characters",
			username: "testuser",
			password: "P@ssw0rd!#$%",
			mockSetup: func(userRepo *MockUserRepository, tokenSvc *MockTokenService) {
				passwordHash := createTestPasswordHash("P@ssw0rd!#$%")
				userRepo.On("GetUserByUsername", "testuser").Return(&domain.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: passwordHash,
				}, nil)
				tokenSvc.On("CreateToken", mock.Anything).Return("valid.jwt.token", nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "valid.jwt.token", token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockTokenSvc := new(MockTokenService)
			tt.mockSetup(mockUserRepo, mockTokenSvc)

			service := NewAuthService(mockUserRepo, mockTokenSvc)
			result, err := service.Login(context.Background(), tt.username, tt.password)

			tt.validate(t, result, err)
			mockUserRepo.AssertExpectations(t)
			mockTokenSvc.AssertExpectations(t)
		})
	}
}
