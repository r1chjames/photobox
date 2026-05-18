package repository

import (
	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
	"os"
	"strings"
)

type UtilityRepository struct {
	dbEnv *db.Env
}

func NewUtilityRepository(dbEnv *db.Env) *UtilityRepository {
	return &UtilityRepository{
		dbEnv,
	}
}

func (ur *UtilityRepository) Ping() error {
	return ur.dbEnv.Ping()
}

func (ur *UtilityRepository) GetAllSettings() ([]*domain.Setting, error) {
	var setting []*domain.Setting
	result := ur.dbEnv.Db.Find(&setting)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func (ur *UtilityRepository) GetSetting(key string) (*domain.Setting, error) {
	var setting domain.Setting
	setting.Key = key
	result := ur.dbEnv.Db.Find(&setting)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (ur *UtilityRepository) UpdateSetting(setting *domain.Setting) error {
	result := ur.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&setting)
	return result.Error
}

func (ur *UtilityRepository) UpdateAllSettings(settings []*domain.Setting) error {
	result := ur.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&settings)
	return result.Error
}

func (ur *UtilityRepository) overwriteAllSettings(settings []*domain.Setting) error {
	result := ur.dbEnv.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&settings)
	return result.Error
}

func (ur *UtilityRepository) CreateBaseSettings(reset bool) error {
	thumbStorage := strings.TrimSpace(os.Getenv("THUMBNAIL_STORAGE"))
	if thumbStorage == "" {
		thumbStorage = "filesystem"
	}

	var settings = []*domain.Setting{
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
			Value:        "0 1 * * *",
		},
		{
			Key:          "thumbnail_storage",
			FriendlyName: "Thumbnail storage backend",
			Category:     "Photo",
			Type:         "Choice",
			Options:      "filesystem,s3",
			Description:  "Where thumbnails are stored. Requires restart to take effect. 'filesystem' stores on disk under THUMBNAIL_DIR, 's3' stores in an S3-compatible object store (MinIO or AWS S3).",
			Value:        thumbStorage,
		},
	}
	if reset {
		return ur.UpdateAllSettings(settings)
	}
	return ur.overwriteAllSettings(settings)
}
