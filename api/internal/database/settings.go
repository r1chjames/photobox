package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
	"log"
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


func UpdateAllSettings(settings *[]Setting) error {
	result := dbConn.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&settings)
	return result.Error
}

func overwriteAllSettings(settings *[]Setting) error {
	result := dbConn.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&settings)
	return result.Error
}

func createBaseSettings(resetSettings bool) {
	var settings = []Setting{
		{
			Key:          "thumbnail_width",
			FriendlyName: "Thumbnail width",
			Category:     "Photo",
			Type:         "Choice",
			Options:      []string{"400", "600", "800", "1000", "1200"},
			Description:  "Thumbnail width used during thumbnail generation",
			Value:        "600",
		},
		{
			Key:          "thumbnail_height",
			FriendlyName: "Thumbnail height",
			Category:     "Photo",
			Type:         "Choice",
			Options:      []string{"400", "600", "800", "1000", "1200"},
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
		err := UpdateAllSettings(&settings)
		if err != nil {
			log.Fatal("Unable to create initial settings")
		}
	} else {
		err := overwriteAllSettings(&settings)
		if err != nil {
			log.Fatal("Unable to create initial settings")
		}
	}
}
