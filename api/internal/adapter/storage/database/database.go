package database

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Env struct {
	Db *gorm.DB
}

func InitDbConnection(appConfig *appconfig.AppConfig) *Env {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: appConfig.DbUrl,
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false,
			}})

	if err != nil {
		slog.Error("Failed to connect database", "error", err)
		os.Exit(1)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get database instance", "error", err)
		os.Exit(1)
	}

	// Set maximum number of idle connections in the pool
	sqlDB.SetMaxIdleConns(10)

	// Set maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// Set maximum lifetime of a connection (reuse connections for up to 1 hour)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Env{Db: db}
}

func (dbEnv *Env) PerformDbSetup() {
	// Ensure the photobox schema exists before creating tables
	if err := dbEnv.Db.Exec("CREATE SCHEMA IF NOT EXISTS photobox").Error; err != nil {
		slog.Error("Failed to create database schema", "error", err)
		os.Exit(1)
	}

	// Migrate the schema
	err := dbEnv.Db.AutoMigrate(&domain.Album{}, &domain.Photo{}, &domain.PhotoTag{}, &domain.PhotoAnalysis{}, &domain.Setting{}, &domain.Job{}, &domain.User{}, &domain.SharedLink{}, &domain.ApiKey{})
	if err != nil {
		slog.Error("Failed to perform database migration", "error", err)
		os.Exit(1)
	}

	// Ensure photo_analysis table exists (AutoMigrate sometimes misses custom TableName models)
	if err := dbEnv.Db.Exec(`CREATE TABLE IF NOT EXISTS photobox.photo_analysis (
		photo_id TEXT PRIMARY KEY,
		model TEXT,
		caption TEXT,
		tags JSONB,
		objects JSONB,
		is_nsfw BOOLEAN DEFAULT FALSE,
		is_portrait BOOLEAN DEFAULT FALSE,
		status TEXT DEFAULT 'pending',
		attempts INTEGER DEFAULT 0,
		error_message TEXT,
		retryable BOOLEAN DEFAULT TRUE,
		started_at TIMESTAMP,
		last_attempt_at TIMESTAMP,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	)`).Error; err != nil {
		slog.Error("Failed to create photo_analysis table", "error", err)
		os.Exit(1)
	}

	// Add columns missing on existing tables (safe to run on new tables too)
	for _, stmt := range []string{
		`ALTER TABLE photobox.photo_analysis ADD COLUMN IF NOT EXISTS started_at TIMESTAMP`,
		`ALTER TABLE photobox.photo_analysis ADD COLUMN IF NOT EXISTS error_message TEXT`,
		`ALTER TABLE photobox.photo_analysis ADD COLUMN IF NOT EXISTS retryable BOOLEAN DEFAULT TRUE`,
	} {
		if err := dbEnv.Db.Exec(stmt).Error; err != nil {
			slog.Error("Failed to add column to photo_analysis", "error", err)
			os.Exit(1)
		}
	}

	// Create GIN indexes for full-text search
	dbEnv.createSearchIndexes()

	// Create GiST index for geospatial queries
	dbEnv.createGeoIndexes()

	// Migrate created_epoch from EXIF dates for photos indexed before the fix
	dbEnv.migratePhotoEpochs()

	// Backfill year and month columns for existing photos
	dbEnv.migratePhotoYearMonth()

	// Drop legacy tags column now that PhotoTag junction table is the sole source of truth
	dbEnv.migrateDropLegacyTagsColumn()

	// Enable query performance tracking
	if err := dbEnv.Db.Exec("CREATE EXTENSION IF NOT EXISTS pg_stat_statements").Error; err != nil {
		slog.Warn("Failed to enable pg_stat_statements", "error", err)
	}
}

func (dbEnv *Env) createSearchIndexes() {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_photo_search ON photobox.photos USING GIN (to_tsvector('english', coalesce(name, '')))`,
		`CREATE INDEX IF NOT EXISTS idx_album_search ON photobox.albums USING GIN (to_tsvector('english', coalesce(name, '')))`,
	}
	for _, idx := range indexes {
		if err := dbEnv.Db.Exec(idx).Error; err != nil {
			slog.Warn("Failed to create search index", "error", err, "index", idx)
		}
	}
}

func (dbEnv *Env) createGeoIndexes() {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_photos_lat_lng ON photobox.photos USING btree (latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL`,
	}
	for _, idx := range indexes {
		if err := dbEnv.Db.Exec(idx).Error; err != nil {
			slog.Warn("Failed to create geo index", "error", err, "index", idx)
		}
	}
}

func (dbEnv *Env) MigrateThumbnailsToFilesystem(photoDir string) {
	var photos []domain.Photo
	result := dbEnv.Db.Where("thumbnail_path = ? OR thumbnail_path IS NULL", "").Where("thumbnail IS NOT NULL AND length(thumbnail) > 0").Find(&photos)
	if result.Error != nil {
		slog.Warn("Failed to query photos for thumbnail migration", "error", result.Error)
		return
	}
	if len(photos) == 0 {
		return
	}
	slog.Info("Migrating thumbnails to filesystem", "count", len(photos))
	thumbsDir := filepath.Join(photoDir, ".thumbnails")
	for _, p := range photos {
		safeId := strings.ReplaceAll(p.ID, "/", "_")
		safeId = strings.ReplaceAll(safeId, "+", "-")
		safeId = strings.ReplaceAll(safeId, "=", "")
		path := filepath.Join(thumbsDir, safeId+".webp")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			slog.Warn("Failed to create thumbnail directory", "error", err)
			continue
		}
		if err := os.WriteFile(path, p.Thumbnail, 0644); err != nil {
			slog.Warn("Failed to write thumbnail to disk", "photo", p.ID, "error", err)
			continue
		}
		dbEnv.Db.Model(&domain.Photo{}).Where("id = ?", p.ID).Update("thumbnail_path", path)
	}
	slog.Info("Thumbnail migration complete")
}

// migratePhotoEpochs updates created_epoch for existing photos from EXIF
// DateTimeOriginal where available. This fixes photos indexed before the
// getPhotoEpoch helper was introduced, which stored the index timestamp
// instead of the actual photo date.
func (dbEnv *Env) migratePhotoEpochs() {
	slog.Info("Starting photo epoch migration")
	result := dbEnv.Db.Exec(`
		UPDATE photobox.photos
		SET created_epoch = CASE
			WHEN metadata->'exif'->>'DateTimeOriginal' ~ '^\d{4}:\d{2}:\d{2}'
			THEN (EXTRACT(EPOCH FROM to_timestamp((metadata->'exif'->>'DateTimeOriginal')::text, 'YYYY:MM:DD HH24:MI:SS'))::bigint * 1000)
			ELSE created_epoch
		END
		WHERE deleted_at IS NULL
	`)
	if result.Error != nil {
		slog.Error("Failed to migrate photo epochs", "error", result.Error)
	} else {
		slog.Info("Photo epoch migration complete", "rows", result.RowsAffected)
	}
}

// migratePhotoYearMonth backfills the year and month columns for existing photos
// using the same EXIF parsing logic as the timeline query. Photos without valid
// EXIF dates fall back to created_epoch.
func (dbEnv *Env) migratePhotoYearMonth() {
	slog.Info("Starting photo year/month migration")
	result := dbEnv.Db.Exec(`
		UPDATE photobox.photos
		SET year = EXTRACT(YEAR FROM COALESCE(
				CASE
					WHEN metadata->'exif'->>'DateTimeOriginal' ~ '^\d{4}:\d{2}:\d{2}'
					THEN to_timestamp((metadata->'exif'->>'DateTimeOriginal')::text, 'YYYY:MM:DD HH24:MI:SS')
					ELSE NULL
				END,
				to_timestamp(created_epoch / 1000)
			))::int,
			month = EXTRACT(MONTH FROM COALESCE(
				CASE
					WHEN metadata->'exif'->>'DateTimeOriginal' ~ '^\d{4}:\d{2}:\d{2}'
					THEN to_timestamp((metadata->'exif'->>'DateTimeOriginal')::text, 'YYYY:MM:DD HH24:MI:SS')
					ELSE NULL
				END,
				to_timestamp(created_epoch / 1000)
			))::int
		WHERE deleted_at IS NULL AND (year IS NULL OR year = 0)
	`)
	if result.Error != nil {
		slog.Error("Failed to migrate photo year/month", "error", result.Error)
	} else {
		slog.Info("Photo year/month migration complete", "rows", result.RowsAffected)
	}
}

// migrateDropLegacyTagsColumn drops the legacy comma-separated tags column
// from the photos table. Tags are now stored exclusively in the PhotoTag
// junction table.
func (dbEnv *Env) migrateDropLegacyTagsColumn() {
	slog.Info("Checking for legacy tags column")
	var count int64
	result := dbEnv.Db.Raw(`
		SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema = 'photobox' AND table_name = 'photos' AND column_name = 'tags'
	`).Scan(&count)
	if result.Error != nil {
		slog.Error("Failed to check for legacy tags column", "error", result.Error)
		return
	}
	if count == 0 {
		return
	}
	result = dbEnv.Db.Exec(`ALTER TABLE photobox.photos DROP COLUMN IF EXISTS tags`)
	if result.Error != nil {
		slog.Error("Failed to drop legacy tags column", "error", result.Error)
	} else {
		slog.Info("Legacy tags column dropped successfully")
	}
}

// Ping checks the database connection is alive
func (dbEnv *Env) Ping() error {
	sqlDB, err := dbEnv.Db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func HandleError(result *gorm.DB) error {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.ErrDataNotFound
	} else if result.Error != nil {
		return result.Error
	}
	return nil
}

func Paginate(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		var offset int
		if page == 1 {
			offset = 0
		} else {
			offset = (page - 1) * limit
		}

		return db.Offset(offset).Limit(limit)
	}
}
