package repository

import (
	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	dbEnv *db.Env
}

func NewUserRepository(dbEnv *db.Env) *UserRepository {
	return &UserRepository{
		dbEnv,
	}
}

func (ur *UserRepository) ListUsers(pageNumber, pageSize int) ([]domain.User, error) {
	var user []domain.User
	result := ur.dbEnv.Db.Scopes(db.Paginate(pageNumber, pageSize)).Find(&user)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetUserById(id string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "id = ? ", id)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "username = ? ", username)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetApprovedUserByUsername(username string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "username = ? AND approved = true", username)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	ur.dbEnv.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&user)
	return user, nil
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
