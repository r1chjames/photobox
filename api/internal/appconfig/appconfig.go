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
	ThumbnailStorage   string
	ThumbnailDir       string
	S3Endpoint         string
	S3AccessKey        string
	S3SecretKey        string
	S3Bucket           string
	S3UseSSL           bool
	PhotoIndexWorkers  int
}

func New() *AppConfig {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "photobox")
	dbPassword := utils.GetEnv("DB_PASSWORD", "photobox")
	dbName := utils.GetEnv("DB_NAME", "photobox")

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	resetSettings, _ := strconv.ParseBool(utils.GetEnv("RESET_SETTINGS", "false"))
	debugMode, _ := strconv.ParseBool(utils.GetEnv("DEBUG_MODE", "false"))
	timezone, _ := time.LoadLocation(utils.GetEnv("TIMEZONE", "Europe/London"))

	token := utils.GetEnv("TOKEN", "")
	tokenDuration, _ := time.ParseDuration(utils.GetEnv("TOKEN_DURATION", "1h"))

	adminUsername := utils.GetEnv("DEFAULT_ADMIN_USERNAME", "admin")
	adminPassword := utils.GetEnv("DEFAULT_ADMIN_PASSWORD", "")
	if adminPassword == "" {
		slog.Error("DEFAULT_ADMIN_PASSWORD environment variable must be set")
		os.Exit(1)
	}

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

	return &AppConfig{
		PhotoDir:           utils.GetEnv("PHOTO_DIR", "/photos"),
		ApiBasePath:        utils.GetEnv("API_BASE_PATH", "/api"),
		DbUrl:              dbURL,
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
		ThumbnailStorage:   thumbnailStorage,
		ThumbnailDir:       thumbnailDir,
		S3Endpoint:         s3Endpoint,
		S3AccessKey:        s3AccessKey,
		S3SecretKey:        s3SecretKey,
		S3Bucket:           s3Bucket,
		S3UseSSL:           s3UseSSL,
		PhotoIndexWorkers:  indexWorkers,
	}
}
