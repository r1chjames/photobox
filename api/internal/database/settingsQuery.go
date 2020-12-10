package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
)

func GetAllSettings() ([]Setting, error) {
	var setting []Setting
	result := dbConn.Find(&setting)
	return setting, result.Error
}

func UpdateSetting(setting Setting) error {
	result := dbConn.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&setting)
	return result.Error
}

func UpdateAllSettings(settings []Setting) error {
	result := dbConn.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&settings)
	return result.Error
}
