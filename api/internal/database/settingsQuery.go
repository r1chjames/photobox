package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
)

func GetAllSettings() ([]Setting, error) {
	var setting []Setting
	result := dbConn.Find(&setting)
	return setting, result.Error
}

func UpdateAllSettings(key string, value string) error {
	var setting Setting
	setting.Key = key
	setting.Value = value
	result := dbConn.Save(setting)
	return result.Error
}
