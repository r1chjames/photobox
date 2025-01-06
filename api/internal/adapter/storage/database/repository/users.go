package repository

import (
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	dbEnv *database.Env
}

func NewUserRepository(dbEnv *database.Env) *UserRepository {
	return &UserRepository{
		dbEnv,
	}
}

func (ur *UserRepository) ListUsers(skip, limit uint64) ([]domain.User, error) {
	var user []domain.User
	result := ur.dbEnv.Db.Find(&user)
	return user, result.Error
}

func (ur *UserRepository) GetUserById(id string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "id = ? ", id)
	return user, result.Error
}

func (ur *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "username = ? ", username)
	return user, result.Error
}

func (ur *UserRepository) GetApprovedUserByUsername(username string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "username = ? AND approved = true", username)
	return user, result.Error
}

func (ur *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	result := ur.dbEnv.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&user)
	return user, result.Error
}

func (ur *UserRepository) UpdateUser(user *domain.User) error {
	result := ur.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Updates(&user)
	return result.Error
}

func (ur *UserRepository) DeleteUser(id string) error {
	var user *domain.User
	user.ID = id
	result := ur.dbEnv.Db.Delete(&user)
	return result.Error
}
