package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
	"log"
)

func (dbEnv *Env) GetAllSettings() ([]Setting, error) {
	var setting []Setting
	result := dbEnv.Db.Find(&setting)
	return setting, result.Error
}

func (dbEnv *Env) GetSetting(key string) (Setting, error) {
	var setting Setting
	setting.Key = key
	result := dbEnv.Db.Find(&setting)
	return setting, result.Error
}

func (dbEnv *Env) UpdateSetting(setting Setting) error {
	result := dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&setting)
	return result.Error
}


func (dbEnv *Env) UpdateAllSettings(settings *[]Setting) error {
	result := dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&settings)
	return result.Error
}

func (dbEnv *Env) overwriteAllSettings(settings *[]Setting) error {
	result := dbEnv.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&settings)
	return result.Error
}

func (dbEnv *Env) createBaseSettings(resetSettings bool) {
	var settings = []Setting{
		{
			Key:          "thumbnail_width",
			FriendlyName: "Thumbnail width",
			Category:     "Photo",
			Type:         "Choice",
			Options:      "400,600,800,1000,1200",
			Description:  "Thumbnail width used during thumbnail generation",
			Value:        "600",
		},
		{
			Key:          "thumbnail_height",
			FriendlyName: "Thumbnail height",
			Category:     "Photo",
			Type:         "Choice",
			Options:      "400,600,800,1000,1200",
			Description:  "Thumbnail height used during thumbnail generation",
			Value:        "600",
		},
		{
			Key:          "default_new_albums_dir",
			FriendlyName: "New album storage location",
			Category:     "System",
			Type:         "Text",
			Description:  "Default location on disk to store new albums",
			Value:        "/photos",
		},
		{
			Key:          "index_frequency_cron",
			FriendlyName: "CRON expression for indexing",
			Category:     "System",
			Type:         "Text",
			Description:  "CRON expression used to initiate indexing",
			Value:        "0 0 23 1/1 * ? *",
		},
	}
	if resetSettings {
		err := dbEnv.UpdateAllSettings(&settings)
		if err != nil {
			log.Fatal("Unable to create initial settings")
		}
	} else {
		err := dbEnv.overwriteAllSettings(&settings)
		if err != nil {
			log.Fatal("Unable to create initial settings")
		}
	}
}
