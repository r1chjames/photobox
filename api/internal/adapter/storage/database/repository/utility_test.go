package repository

import (
	"testing"

	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// TestGetAllSettings_Success tests successful retrieval of all settings
func TestGetAllSettings_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedSettings := sqlmock.NewRows([]string{"key", "value", "friendly_name", "category", "type", "options", "description"}).
		AddRow("thumbnail_width", "600", "Thumbnail width", "Photo", "Choice", "400,600,800,1000,1200", "Thumbnail width used during thumbnail generation").
		AddRow("thumbnail_height", "600", "Thumbnail height", "Photo", "Choice", "400,600,800,1000,1200", "Thumbnail height used during thumbnail generation").
		AddRow("default_new_albums_dir", "/photos", "New album storage location", "System", "Text", "", "Default location on disk to store new albums")

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "settings"`, expectedSettings)

	repo := NewUtilityRepository(env)
	settings, err := repo.GetAllSettings()

	if err != nil {
		t.Errorf("error was not expected while getting all settings: %s", err)
	}

	if len(settings) != 3 {
		t.Errorf("expected 3 settings, got %d", len(settings))
	}

	if settings[0].Key != "thumbnail_width" {
		t.Errorf("expected first setting key to be thumbnail_width, got %s", settings[0].Key)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetAllSettings_Empty tests when no settings exist
func TestGetAllSettings_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedSettings := sqlmock.NewRows([]string{"key", "value", "friendly_name", "category", "type", "options", "description"})

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "settings"`, expectedSettings)

	repo := NewUtilityRepository(env)
	settings, err := repo.GetAllSettings()

	if err != nil {
		t.Errorf("error was not expected for empty result: %s", err)
	}

	if settings == nil {
		t.Error("expected empty slice, got nil")
	}

	if len(settings) != 0 {
		t.Errorf("expected 0 settings, got %d", len(settings))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetSetting_Success tests successful single setting retrieval
func TestGetSetting_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	settingKey := "default_new_albums_dir"
	expectedSetting := sqlmock.NewRows([]string{"key", "value", "friendly_name", "category", "type", "options", "description"}).
		AddRow(settingKey, "/photos", "New album storage location", "System", "Text", "", "Default location on disk to store new albums")

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "settings"`, expectedSetting)

	repo := NewUtilityRepository(env)
	setting, err := repo.GetSetting(settingKey)

	if err != nil {
		t.Errorf("error was not expected while getting setting: %s", err)
	}

	if setting == nil {
		t.Error("expected setting to be returned, got nil")
	}

	if setting.Key != settingKey {
		t.Errorf("expected setting key %s, got %s", settingKey, setting.Key)
	}

	if setting.Value != "/photos" {
		t.Errorf("expected setting value /photos, got %s", setting.Value)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetSetting_NotFound tests when setting doesn't exist
func TestGetSetting_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	settingKey := "nonexistent_key"
	expectedSetting := sqlmock.NewRows([]string{"key", "value", "friendly_name", "category", "type", "options", "description"})

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "settings"`, expectedSetting)

	repo := NewUtilityRepository(env)
	setting, err := repo.GetSetting(settingKey)

	if err != nil {
		t.Errorf("error was not expected for empty result: %s", err)
	}

	if setting == nil {
		t.Error("expected setting struct, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateSetting_Success tests successful single setting update
func TestUpdateSetting_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	setting := &domain.Setting{
		Key:          "thumbnail_width",
		Value:        "800",
		FriendlyName: "Thumbnail width",
		Category:     "Photo",
		Type:         "Choice",
		Options:      "400,600,800,1000,1200",
		Description:  "Thumbnail width used during thumbnail generation",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "settings"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewUtilityRepository(env)
	err := repo.UpdateSetting(setting)

	if err != nil {
		t.Errorf("error was not expected while updating setting: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateAllSettings_Success tests successful bulk settings update
func TestUpdateAllSettings_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	settings := []*domain.Setting{
		{
			Key:          "thumbnail_width",
			Value:        "800",
			FriendlyName: "Thumbnail width",
			Category:     "Photo",
			Type:         "Choice",
			Options:      "400,600,800,1000,1200",
			Description:  "Thumbnail width used during thumbnail generation",
		},
		{
			Key:          "thumbnail_height",
			Value:        "800",
			FriendlyName: "Thumbnail height",
			Category:     "Photo",
			Type:         "Choice",
			Options:      "400,600,800,1000,1200",
			Description:  "Thumbnail height used during thumbnail generation",
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "settings"`).
		WillReturnResult(sqlmock.NewResult(2, 2))
	mock.ExpectCommit()

	repo := NewUtilityRepository(env)
	err := repo.UpdateAllSettings(settings)

	if err != nil {
		t.Errorf("error was not expected while updating settings: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateBaseSettings_WithoutReset tests creating base settings without reset
func TestCreateBaseSettings_WithoutReset(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	// Without reset, uses overwriteAllSettings with DO NOTHING
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "settings"`).
		WillReturnResult(sqlmock.NewResult(4, 4))
	mock.ExpectCommit()

	repo := NewUtilityRepository(env)
	err := repo.CreateBaseSettings(false)

	if err != nil {
		t.Errorf("error was not expected while creating base settings: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateBaseSettings_WithReset tests creating base settings with reset (upsert)
func TestCreateBaseSettings_WithReset(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	// With reset, uses UpdateAllSettings with UPDATE ALL
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "settings"`).
		WillReturnResult(sqlmock.NewResult(4, 4))
	mock.ExpectCommit()

	repo := NewUtilityRepository(env)
	err := repo.CreateBaseSettings(true)

	if err != nil {
		t.Errorf("error was not expected while creating base settings with reset: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
