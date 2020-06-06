package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

func transformRowsIntoSettingsArray(results *sql.Rows) interface{} {
	defer results.Close()
	var records []SettingRecord
	for results.Next() {
		var record SettingRecord
		err := results.Scan(
			&record.Setting,
			&record.Value)
		if err != nil {
			log.Fatal(fmt.Sprintf("An error parsing settings record: %s", err))
		}
		records = append(records, record)
	}
	return records
}

func GetAllSettings(appConfig AppConfig) ([]SettingRecord, error) {
	allPhotoFields := "SETTING, VALUE"
	query := fmt.Sprintf("SELECT %s FROM settings", allPhotoFields)

	queryResults := executeTransformDbQuery(appConfig, query, transformRowsIntoSettingsArray).([]SettingRecord)
	if len(queryResults) > 0 {
		return queryResults, nil
	}
	return []SettingRecord{}, errors.New("no settings found")
}

func UpdateAllSettings(appConfig AppConfig, c *gin.Context) error {
	var settings []SettingRecord
	err := c.BindJSON(&settings)
	if err != nil {
		return err
	}

	for _, setting := range settings {
		settingsInsertQuery := fmt.Sprintf("INSERT IGNORE INTO settings VALUES ('%s','%s')",
			setting.Value,
			setting.Value)
		executeDbInsert(appConfig, settingsInsertQuery)
	}
	return nil
}
