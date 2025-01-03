package database

import (
	"gorm.io/gorm/clause"
)

func (dbEnv *Env) GetAllUsers() ([]User, error) {
	var user []User
	result := dbEnv.Db.Find(&user)
	return user, result.Error
}

func (dbEnv *Env) GetUserByUsername(username string) (User, error) {
	var user User
	result := dbEnv.Db.Find(&user, "username = ? ", username)
	return user, result.Error
}

func (dbEnv *Env) GetApprovedUserByUsername(username string) (User, error) {
	var user User
	result := dbEnv.Db.Find(&user, "username = ? AND approved = true", username)
	return user, result.Error
}

func (dbEnv *Env) CreateUser(user User) error {
	result := dbEnv.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&user)
	return result.Error
}

func (dbEnv *Env) UpdateUser(user User) error {
	result := dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Updates(&user)
	return result.Error
}

func (dbEnv *Env) DeleteUser(id string) error {
	var user User
	user.ID = id
	result := dbEnv.Db.Delete(&user)
	return result.Error
}
