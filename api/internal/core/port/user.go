package port

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

//go:generate mockgen -source=user.go -destination=mock/user.go -package=mock

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(user *domain.User) (*domain.User, error)
	// CreateUserWithPersonalWorkspace inserts a user, their personal workspace,
	// and the owner membership in a single transaction (issue #74 §8).
	CreateUserWithPersonalWorkspace(user *domain.User, workspace *domain.Workspace) error
	// ListUsers selects a list of users with pagination
	ListUsers(pageNumber, pageSize int) ([]domain.User, error)
	// GetUserById selects a user by id
	GetUserById(id string) (*domain.User, error)
	// GetUserByUsername selects a user by id
	GetUserByUsername(username string) (*domain.User, error)
	// GetApprovedUserByUsername selects a user by email
	GetApprovedUserByUsername(username string) (*domain.User, error)
	// UpdateUser updates a user
	UpdateUser(user *domain.User) error
	// DeleteUser deletes a user
	DeleteUser(id string) error
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	// Register registers a new user
	Register(user *domain.User) (*domain.User, error)
	// RegisterWithPersonalWorkspace registers a user plus their personal
	// workspace and owner membership atomically (issue #74 §8).
	RegisterWithPersonalWorkspace(user *domain.User, workspaceName string) (*domain.User, *domain.Workspace, error)
	// CreateUser registers a new user
	CreateUser(user *domain.User) (*domain.User, error)
	// GetUser returns a user by id
	GetUser(id string) (*domain.User, error)
	// ListUsers returns a list of users with pagination
	ListUsers(pageNumber, pageSize int) ([]domain.User, error)
	// UpdateUser updates a user
	UpdateUser(user *domain.User) (*domain.User, error)
	// DeleteUser deletes a user
	DeleteUser(id string) error
}
