package service

import (
	"errors"

	"github.com/google/uuid"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/auth"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type UserService struct {
	repo port.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo,
	}
}

// Register creates a new user with viewer permissions
func (us *UserService) Register(user *domain.User) (*domain.User, error) {
	hashedPassword, err := auth.CreateHash(user.Password, auth.DefaultArgon2idHash())
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.ID = uuid.New().String()
	user.Password = hashedPassword
	user.Role = domain.VIEWER

	user, err = us.repo.CreateUser(user)
	if err != nil {
		if errors.Is(err, domain.ErrConflictingData) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

// CreateUser creates a new user
func (us *UserService) CreateUser(user *domain.User) (*domain.User, error) {
	hashedPassword, err := auth.CreateHash(user.Password, auth.DefaultArgon2idHash())
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.ID = uuid.New().String()
	user.Password = hashedPassword

	user, err = us.repo.CreateUser(user)
	if err != nil {
		if errors.Is(err, domain.ErrConflictingData) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

// GetUser gets a user by ID
func (us *UserService) GetUser(id string) (*domain.User, error) {
	var user *domain.User

	user, err := us.repo.GetUserById(id)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

// ListUsers lists all users
func (us *UserService) ListUsers(pageNumber, pageSize int) ([]domain.User, error) {
	var users []domain.User

	users, err := us.repo.ListUsers(pageNumber, pageSize)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return users, nil
}

// UpdateUser updates a user's name, email, and password
func (us *UserService) UpdateUser(user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserById(user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	emptyData := user.Username == "" &&
		user.Email == "" &&
		user.Password == "" &&
		user.Role == ""
	sameData := existingUser.Username == user.Username &&
		existingUser.Email == user.Email &&
		existingUser.Role == user.Role
	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = auth.CreateHash(user.Password, auth.DefaultArgon2idHash())
		if err != nil {
			return nil, domain.ErrInternal
		}
	}

	user.Password = hashedPassword

	err = us.repo.UpdateUser(user)
	if err != nil {
		if errors.Is(err, domain.ErrConflictingData) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(id string) error {
	_, err := us.repo.GetUserById(id)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return err
		}
		return domain.ErrInternal
	}

	return us.repo.DeleteUser(id)
}
