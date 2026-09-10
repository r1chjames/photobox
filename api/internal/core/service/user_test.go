package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockUserRepository is a mock implementation of port.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		// Return the modified input user when mock returns nil
		return user, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) CreateUserWithPersonalWorkspace(user *domain.User, workspace *domain.Workspace) error {
	args := m.Called(user, workspace)
	return args.Error(0)
}

func (m *MockUserRepository) ListUsers(pageNumber, pageSize int) ([]domain.User, error) {
	args := m.Called(pageNumber, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetApprovedUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestRegisterWithPersonalWorkspace tests the signup flow that atomically
// creates a user + personal workspace + owner membership (issue #74 §8).
func TestRegisterWithPersonalWorkspace(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("CreateUserWithPersonalWorkspace", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			user := args.Get(0).(*domain.User)
			user.ID = "user-123"
			ws := args.Get(1).(*domain.Workspace)
			ws.ID = "ws-123"
			ws.CreatedAt = time.Now()
		}).Return(nil)

	service := NewUserService(mockRepo)
	user, ws, err := service.RegisterWithPersonalWorkspace(&domain.User{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "password123",
	}, "Alice's Photos")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "alice", user.Username)
	assert.Equal(t, domain.VIEWER, user.Role)
	assert.NotEqual(t, "password123", user.Password, "password should be hashed")
	assert.NotNil(t, ws)
	assert.Equal(t, "Alice's Photos", ws.Name)
	assert.Equal(t, "alices-photos", ws.Slug)
	assert.NotEmpty(t, ws.ID)
	mockRepo.AssertExpectations(t)
}

// TestRegisterWithPersonalWorkspace_RepoError asserts a repository failure
// surfaces as ErrInternal and no partial registration is reported.
func TestRegisterWithPersonalWorkspace_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("CreateUserWithPersonalWorkspace", mock.Anything, mock.Anything).
		Return(errors.New("tx failed"))

	service := NewUserService(mockRepo)
	user, ws, err := service.RegisterWithPersonalWorkspace(&domain.User{
		Username: "bob",
		Email:    "bob@example.com",
		Password: "password123",
	}, "Bob's Photos")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInternal))
	assert.Nil(t, user)
	assert.Nil(t, ws)
	mockRepo.AssertExpectations(t)
}

// TestRegister tests the user registration functionality
func TestRegister(t *testing.T) {
	tests := []struct {
		name          string
		inputUser     *domain.User
		mockSetup     func(*MockUserRepository)
		expectedError error
		validate      func(*testing.T, *domain.User, error)
	}{
		{
			name: "successful registration",
			inputUser: &domain.User{
				Username: "testuser",
				Password: "plainPassword123",
				Email:    "test@example.com",
			},
			mockSetup: func(m *MockUserRepository) {
				// Mock will return the input user after setting ID
				m.On("CreateUser", mock.MatchedBy(func(u *domain.User) bool {
					return u.Username == "testuser" && u.Password != "plainPassword123"
				})).Run(func(args mock.Arguments) {
					// Set ID on the user that was passed in
					user := args.Get(0).(*domain.User)
					user.ID = "generated-id"
				}).Return(nil, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "testuser", user.Username)
				// Password should be hashed (not equal to plain password)
				assert.NotEqual(t, "plainPassword123", user.Password)
				assert.NotEmpty(t, user.Password)
			},
		},
		{
			name: "registration with conflicting username",
			inputUser: &domain.User{
				Username: "existinguser",
				Password: "password123",
				Email:    "test@example.com",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything).Return(nil, domain.ErrConflictingData)
			},
			expectedError: domain.ErrConflictingData,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrConflictingData))
				assert.Nil(t, user)
			},
		},
		{
			name: "registration with repository error",
			inputUser: &domain.User{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrInternal))
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)

			// Execute
			result, err := service.Register(tt.inputUser)

			// Validate
			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestCreateUser tests the CreateUser method
func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		inputUser     *domain.User
		mockSetup     func(*MockUserRepository)
		expectedError error
		validate      func(*testing.T, *domain.User, error)
	}{
		{
			name: "successful user creation",
			inputUser: &domain.User{
				Username: "adminuser",
				Password: "securePassword",
				Email:    "admin@example.com",
				Role:     domain.ADMINISTRATOR,
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything).Return(&domain.User{
					ID:       "admin-123",
					Username: "adminuser",
					Email:    "admin@example.com",
					Role:     domain.ADMINISTRATOR,
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "adminuser", user.Username)
				assert.Equal(t, domain.ADMINISTRATOR, user.Role)
				// Password should be hashed
				assert.NotEqual(t, "securePassword", user.Password)
			},
		},
		{
			name: "user creation with conflict",
			inputUser: &domain.User{
				Username: "duplicate",
				Password: "password",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything).Return(nil, domain.ErrConflictingData)
			},
			expectedError: domain.ErrConflictingData,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrConflictingData))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)

			result, err := service.CreateUser(tt.inputUser)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetUser tests retrieving a user by ID
func TestGetUser(t *testing.T) {
	tests := []struct {
		name          string
		userId        string
		mockSetup     func(*MockUserRepository)
		expectedError error
		validate      func(*testing.T, *domain.User, error)
	}{
		{
			name:   "successfully get existing user",
			userId: "user-123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
					Email:    "test@example.com",
					Role:     domain.VIEWER,
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "user-123", user.ID)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name:   "user not found",
			userId: "nonexistent",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "nonexistent").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrDataNotFound,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrDataNotFound))
				assert.Nil(t, user)
			},
		},
		{
			name:   "repository error",
			userId: "user-123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(nil, errors.New("database error"))
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrInternal))
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)

			result, err := service.GetUser(tt.userId)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListUsers tests listing users with pagination
func TestListUsers(t *testing.T) {
	tests := []struct {
		name       string
		pageNumber int
		pageSize   int
		mockSetup  func(*MockUserRepository)
		validate   func(*testing.T, []domain.User, error)
	}{
		{
			name:       "successfully list users",
			pageNumber: 1,
			pageSize:   10,
			mockSetup: func(m *MockUserRepository) {
				users := []domain.User{
					{ID: "user-1", Username: "user1", Role: domain.VIEWER},
					{ID: "user-2", Username: "user2", Role: domain.CONTRIBUTOR},
					{ID: "user-3", Username: "user3", Role: domain.ADMINISTRATOR},
				}
				m.On("ListUsers", 1, 10).Return(users, nil)
			},
			validate: func(t *testing.T, users []domain.User, err error) {
				assert.NoError(t, err)
				assert.Len(t, users, 3)
				assert.Equal(t, "user1", users[0].Username)
				assert.Equal(t, "user2", users[1].Username)
			},
		},
		{
			name:       "empty user list",
			pageNumber: 5,
			pageSize:   10,
			mockSetup: func(m *MockUserRepository) {
				m.On("ListUsers", 5, 10).Return([]domain.User{}, nil)
			},
			validate: func(t *testing.T, users []domain.User, err error) {
				assert.NoError(t, err)
				assert.Empty(t, users)
			},
		},
		{
			name:       "repository error",
			pageNumber: 1,
			pageSize:   10,
			mockSetup: func(m *MockUserRepository) {
				m.On("ListUsers", 1, 10).Return(nil, errors.New("database error"))
			},
			validate: func(t *testing.T, users []domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrInternal))
				assert.Nil(t, users)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)

			result, err := service.ListUsers(tt.pageNumber, tt.pageSize)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateUser tests updating user information
func TestUpdateUser(t *testing.T) {
	existingUser := &domain.User{
		ID:       "user-123",
		Username: "oldusername",
		Email:    "old@example.com",
		Role:     domain.VIEWER,
		Password: "hashedOldPassword",
	}

	tests := []struct {
		name          string
		updateUser    *domain.User
		mockSetup     func(*MockUserRepository)
		expectedError error
		validate      func(*testing.T, *domain.User, error)
	}{
		{
			name: "successfully update username and email",
			updateUser: &domain.User{
				ID:       "user-123",
				Username: "newusername",
				Email:    "new@example.com",
				Role:     domain.VIEWER,
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(existingUser, nil)
				m.On("UpdateUser", mock.MatchedBy(func(u *domain.User) bool {
					return u.Username == "newusername" && u.Email == "new@example.com"
				})).Return(nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "newusername", user.Username)
				assert.Equal(t, "new@example.com", user.Email)
			},
		},
		{
			name: "successfully update password with different username",
			updateUser: &domain.User{
				ID:       "user-123",
				Username: "newusername", // Changed from oldusername
				Email:    "old@example.com",
				Password: "newPlainPassword",
				Role:     domain.VIEWER,
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(existingUser, nil)
				m.On("UpdateUser", mock.MatchedBy(func(u *domain.User) bool {
					// Password should be hashed and username changed
					return u.Password != "" && u.Password != "newPlainPassword" && u.Username == "newusername"
				})).Return(nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				// Password should be hashed
				assert.NotEqual(t, "newPlainPassword", user.Password)
				assert.NotEmpty(t, user.Password)
				assert.Equal(t, "newusername", user.Username)
			},
		},
		{
			name: "update user not found",
			updateUser: &domain.User{
				ID:       "nonexistent",
				Username: "newusername",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "nonexistent").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrDataNotFound,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrDataNotFound))
				assert.Nil(t, user)
			},
		},
		{
			name: "no updated data - empty fields",
			updateUser: &domain.User{
				ID:       "user-123",
				Username: "",
				Email:    "",
				Password: "",
				Role:     "",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(existingUser, nil)
			},
			expectedError: domain.ErrNoUpdatedData,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrNoUpdatedData))
				assert.Nil(t, user)
			},
		},
		{
			name: "no updated data - same data",
			updateUser: &domain.User{
				ID:       "user-123",
				Username: "oldusername",
				Email:    "old@example.com",
				Role:     domain.VIEWER,
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(existingUser, nil)
			},
			expectedError: domain.ErrNoUpdatedData,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrNoUpdatedData))
				assert.Nil(t, user)
			},
		},
		{
			name: "update with conflicting data",
			updateUser: &domain.User{
				ID:       "user-123",
				Username: "duplicateusername",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(existingUser, nil)
				m.On("UpdateUser", mock.Anything).Return(domain.ErrConflictingData)
			},
			expectedError: domain.ErrConflictingData,
			validate: func(t *testing.T, user *domain.User, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrConflictingData))
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)

			result, err := service.UpdateUser(tt.updateUser)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestDeleteUser tests deleting a user
func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userId        string
		mockSetup     func(*MockUserRepository)
		expectedError error
	}{
		{
			name:   "successfully delete user",
			userId: "user-123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(&domain.User{
					ID:       "user-123",
					Username: "testuser",
				}, nil)
				m.On("DeleteUser", "user-123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "delete non-existent user",
			userId: "nonexistent",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "nonexistent").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrDataNotFound,
		},
		{
			name:   "repository error on get",
			userId: "user-123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(nil, errors.New("database error"))
			},
			expectedError: domain.ErrInternal,
		},
		{
			name:   "repository error on delete",
			userId: "user-123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserById", "user-123").Return(&domain.User{
					ID: "user-123",
				}, nil)
				m.On("DeleteUser", "user-123").Return(errors.New("deletion failed"))
			},
			expectedError: errors.New("deletion failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)

			err := service.DeleteUser(tt.userId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, domain.ErrDataNotFound) {
					assert.True(t, errors.Is(err, domain.ErrDataNotFound))
				} else if errors.Is(tt.expectedError, domain.ErrInternal) {
					assert.True(t, errors.Is(err, domain.ErrInternal))
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
