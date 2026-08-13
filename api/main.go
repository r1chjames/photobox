package main

import (
	"context"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/auth"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/http"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/repository"
	filesystemRepos "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/filesystem/repository"
	thumbFs "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/thumbnail/filesystem"
	thumbS3 "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/thumbnail/s3"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/components"
	ws "gitlab.com/r1chjames/photobox/api/internal/components/websocket"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"log/slog"
	"os"
)

// version is set at build time via -ldflags "-X main.version=..."
var version = "0.0.0-dev"

func main() {
	appConfig := appconfig.New()

	dbEnv := database.InitDbConnection(appConfig)
	dbEnv.PerformDbSetup()
	if appConfig.ThumbnailStorage != "s3" {
		dbEnv.MigrateThumbnailsToFilesystem(appConfig.PhotoDir)
	}

	services := setupAppServices(dbEnv, appConfig)
	if err := services.utilityService.CreateBaseSettings(appConfig.ResetSettings); err != nil {
		slog.Error("Failed to create base settings", "error", err)
		os.Exit(1)
	}
	services.jobService.CreateBaseJobs()
	services.scheduler.StopAllRunningJobs()
	services.scheduler.AddScheduledJobs()
	addDefaultAdminUser(services.userService, appConfig)

	router, err := setupHttpHandlers(appConfig, services)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	err = router.Run()
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}
}

func addDefaultAdminUser(userService *service.UserService, config *appconfig.AppConfig) {
	_, err := userService.CreateUser(&domain.User{
		Username: config.AdminUsername,
		Password: config.AdminPassword,
		Role:     domain.ADMINISTRATOR,
	})

	if err != nil {
		slog.Info("Error creating default admin user", "error", err)
	}
}

type AppServices struct {
	scheduler         *components.Scheduler
	tokenService      port.TokenService
	userService       *service.UserService
	authService       *service.AuthService
	photoService      *service.PhotoService
	albumService      *service.AlbumService
	jobService        *service.JobService
	utilityService    *service.UtilityService
	filesystemService *service.FilesystemService
	shareService      *service.ShareService
	cacheService      *service.CacheService
	wsHub             *ws.Hub
}

func setupAppServices(dbEnv *database.Env, config *appconfig.AppConfig) *AppServices {
	token, err := auth.New(config.Token, config.TokenDuration)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	// User
	userRepo := repository.NewUserRepository(dbEnv)
	userService := service.NewUserService(userRepo)

	// Auth
	authService := service.NewAuthService(userRepo, token)

	// Album
	albumRepo := repository.NewAlbumRepository(dbEnv)
	albumService := service.NewAlbumService(albumRepo, *config)

	// Utility
	utilityRepo := repository.NewUtilityRepository(dbEnv)
	utilityService := service.NewUtilityService(utilityRepo)

	// Job
	jobRepo := repository.NewJobRepository(dbEnv)
	jobService := service.NewJobService(jobRepo)

	filesystemRepo := filesystemRepos.NewFilesystemRepository(*config, jobService)
	filesystemService := service.NewFilesystemService(filesystemRepo, jobService, utilityService, config.PhotoIndexWorkers, config.ThumbnailStorage)

	// Cache
	cacheService := service.NewCacheService(*config)

	// AI
	var aiService port.AIService
	if config.AIEnabled {
		aiService = service.NewOllamaClient(*config)
		slog.Info("AI service enabled", "model", config.OllamaModel, "host", config.OllamaHost)
	}

	// WebSocket hub for real-time events
	wsHub := ws.NewHub(context.Background())

	// Photo
	photoRepo := repository.NewPhotoRepository(dbEnv)
	// Thumbnail storage adapter
	var thumbnailStorage port.ThumbnailStorage
	switch config.ThumbnailStorage {
	case "s3":
		s3store, err := thumbS3.New(thumbS3.Config{
			Endpoint:  config.S3Endpoint,
			AccessKey: config.S3AccessKey,
			SecretKey: config.S3SecretKey,
			Bucket:    config.S3Bucket,
			UseSSL:    config.S3UseSSL,
		})
		if err != nil {
			slog.Error("Failed to create S3 thumbnail storage", "error", err)
			os.Exit(1)
		}
		thumbnailStorage = s3store
	default:
		thumbnailStorage = thumbFs.New(config.ThumbnailDir)
	}

	photoService := service.NewPhotoService(photoRepo, albumService, filesystemService, cacheService, aiService, *config, thumbnailStorage, wsHub, jobService)

	// Share
	shareRepo := repository.NewShareRepository(dbEnv)
	shareService := service.NewShareService(shareRepo, photoService, albumService)

	// Cron
	return &AppServices{
		components.NewScheduler(utilityService, jobService, photoService, *config),
		token,
		userService,
		authService,
		photoService,
		albumService,
		jobService,
		utilityService,
		filesystemService,
		shareService,
		cacheService,
		wsHub,
	}
}

func setupHttpHandlers(
	config *appconfig.AppConfig,
	appServices *AppServices) (*http.Router, error) {

	userHandler := http.NewUserHandler(appServices.userService)
	authHandler := http.NewAuthHandler(appServices.authService, config.TokenDuration)
	photoHandler := http.NewPhotoHandler(appServices.photoService, appServices.jobService)
	albumHandler := http.NewAlbumHandler(appServices.albumService)
	utilityHandler := http.NewUtilityHandler(appServices.utilityService, version)
	healthHandler := http.NewHealthHandler(
		appServices.utilityService.Ping,
		config.PhotoDir,
		appServices.cacheService.Ping,
		version,
	)
	searchHandler := http.NewSearchHandler(appServices.photoService, appServices.albumService)
	shareHandler := http.NewShareHandler(appServices.shareService)
	wsHandler := http.NewWebSocketHandler(appServices.wsHub)

	return http.NewRouter(
		*config,
		appServices.tokenService,
		*authHandler,
		photoHandler,
		*albumHandler,
		*utilityHandler,
		healthHandler,
		*userHandler,
		*searchHandler,
		*shareHandler,
		wsHandler,
	)
}
