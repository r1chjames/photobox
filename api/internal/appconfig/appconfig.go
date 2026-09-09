package appconfig

import (
	"fmt"
	"log/slog"
	"os"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	PhotoDir           string
	ApiBasePath        string
	DbUrl              string
	DbPassword         string
	ResetSettings      bool
	DebugMode          bool
	Timezone           *time.Location
	Token              string
	TokenDuration      time.Duration
	AdminUsername      string
	AdminPassword      string
	CorsAllowedOrigins []string
	CacheHost          string
	CachePort          string
	CachePassword      string
	CacheEnabled       bool
	CacheDB            int
	AIEnabled          bool
	OllamaHost         string
	OllamaModel        string
	FaceEngineEnabled  bool
	FaceEngineURL      string
	ThumbnailStorage   string
	ThumbnailDir       string
	S3Endpoint         string
	S3AccessKey        string
	S3SecretKey        string
	S3Bucket           string
	S3UseSSL           bool
	PhotoIndexWorkers  int
	TrashRetentionDays int
	GeocodeEndpoint    string
	FFmpegPoolSize     int
}

func New() *AppConfig {
	cfg := load()

	if err, warnings := cfg.validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	} else {
		for _, w := range warnings {
			slog.Warn("configuration warning", "warning", w)
		}
	}

	cfg.logEffectiveConfig()

	return cfg
}

// load parses environment variables into an AppConfig without validating or
// exiting. It is exported for tests and called by New.
func load() *AppConfig {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "photobox")
	dbPassword := utils.GetEnv("DB_PASSWORD", "photobox")
	dbName := utils.GetEnv("DB_NAME", "photobox")
	dbSSLMode := utils.GetEnv("DB_SSL_MODE", "require")

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=photobox,public", dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)
	resetSettings, _ := strconv.ParseBool(utils.GetEnv("RESET_SETTINGS", "false"))
	debugMode, _ := strconv.ParseBool(utils.GetEnv("DEBUG_MODE", "false"))
	timezone, _ := time.LoadLocation(utils.GetEnv("TIMEZONE", "Europe/London"))

	// TOKEN and DEFAULT_ADMIN_PASSWORD are required and validated below; use
	// os.Getenv here so a missing value produces a clean aggregated error
	// rather than a panic from utils.GetEnv.
	token := os.Getenv("TOKEN")
	tokenDuration, _ := time.ParseDuration(utils.GetEnv("TOKEN_DURATION", "1h"))

	adminUsername := utils.GetEnv("DEFAULT_ADMIN_USERNAME", "admin")
	adminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")

	// Parse CORS allowed origins - comma-separated list
	corsOriginsStr := utils.GetEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	var corsOrigins []string
	if corsOriginsStr != "" {
		corsOrigins = strings.Split(corsOriginsStr, ",")
		// Trim whitespace from each origin
		for i, origin := range corsOrigins {
			corsOrigins[i] = strings.TrimSpace(origin)
		}
	}

	cacheEnabled, _ := strconv.ParseBool(utils.GetEnv("CACHE_ENABLED", "false"))
	cacheDB, _ := strconv.Atoi(utils.GetEnv("CACHE_DB", "0"))

	aiEnabled, _ := strconv.ParseBool(utils.GetEnv("AI_ENABLED", "false"))

	faceEngineEnabled, _ := strconv.ParseBool(utils.GetEnv("FACE_ENGINE_ENABLED", "false"))
	faceEngineURL := utils.GetEnv("FACE_ENGINE_URL", "http://localhost:5010")

	thumbnailStorage := utils.GetEnv("THUMBNAIL_STORAGE", "filesystem")
	thumbnailDir := utils.GetEnv("THUMBNAIL_DIR", "/thumbnails")

	// S3/MinIO config — optional, only used when THUMBNAIL_STORAGE=s3.
	// Use os.Getenv (not utils.GetEnv) so they don't panic when unset.
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		s3Bucket = "photobox-thumbnails"
	}
	s3UseSSL, _ := strconv.ParseBool(os.Getenv("S3_USE_SSL"))

	indexWorkers, _ := strconv.Atoi(utils.GetEnv("PHOTO_INDEX_WORKERS", fmt.Sprintf("%d", runtime.NumCPU())))
	if indexWorkers < 1 {
		indexWorkers = 1
	}
	if indexWorkers > runtime.NumCPU() {
		indexWorkers = runtime.NumCPU()
	}

	// Trash retention in days; 0 disables automatic trash cleanup.
	trashRetentionDays, _ := strconv.Atoi(utils.GetEnv("TRASH_RETENTION_DAYS", "60"))
	if trashRetentionDays < 0 {
		trashRetentionDays = 0
	}

	// Nominatim-compatible reverse geocoding endpoint (empty = public OSM).
	geocodeEndpoint := os.Getenv("GEOCODE_ENDPOINT")

	// Concurrent ffmpeg processes for video thumbnails; 0/absent = 2.
	ffmpegPoolSize, _ := strconv.Atoi(utils.GetEnv("FFMPEG_POOL_SIZE", "2"))
	if ffmpegPoolSize < 1 {
		ffmpegPoolSize = 2
	}

	cfg := &AppConfig{
		PhotoDir:           utils.GetEnv("PHOTO_DIR", "/photos"),
		ApiBasePath:        utils.GetEnv("API_BASE_PATH", "/api"),
		DbUrl:              dbURL,
		DbPassword:         dbPassword,
		ResetSettings:      resetSettings,
		DebugMode:          debugMode,
		Timezone:           timezone,
		Token:              token,
		TokenDuration:      tokenDuration,
		AdminUsername:      adminUsername,
		AdminPassword:      adminPassword,
		CorsAllowedOrigins: corsOrigins,
		CacheHost:          utils.GetEnv("CACHE_HOST", "localhost"),
		CachePort:          utils.GetEnv("CACHE_PORT", "6379"),
		CachePassword:      utils.GetEnv("CACHE_PASSWORD", "-"),
		CacheEnabled:       cacheEnabled,
		CacheDB:            cacheDB,
		AIEnabled:          aiEnabled,
		OllamaHost:         utils.GetEnv("OLLAMA_HOST", "http://localhost:11434"),
		OllamaModel:        utils.GetEnv("OLLAMA_MODEL", "moondream"),
		FaceEngineEnabled:  faceEngineEnabled,
		FaceEngineURL:      faceEngineURL,
		ThumbnailStorage:   thumbnailStorage,
		ThumbnailDir:       thumbnailDir,
		S3Endpoint:         s3Endpoint,
		S3AccessKey:        s3AccessKey,
		S3SecretKey:        s3SecretKey,
		S3Bucket:           s3Bucket,
		S3UseSSL:           s3UseSSL,
		PhotoIndexWorkers:  indexWorkers,
		TrashRetentionDays: trashRetentionDays,
		GeocodeEndpoint:    geocodeEndpoint,
		FFmpegPoolSize:     ffmpegPoolSize,
	}

	return cfg
}
