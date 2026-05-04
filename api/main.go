package main

import (
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/auth"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/http"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/repository"
	filesystemRepos "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/filesystem/repository"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/components"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"log/slog"
	"os"
)

func main() {
	appConfig := appconfig.New()

	dbEnv := database.InitDbConnection(appConfig)
	dbEnv.PerformDbSetup()
	dbEnv.MigrateThumbnailsToFilesystem(appConfig.PhotoDir)

	services := setupAppServices(dbEnv, appConfig)
	services.utilityService.CreateBaseSettings(appConfig.ResetSettings)
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
	filesystemService := service.NewFilesystemService(filesystemRepo, jobService, utilityService)

	// Cache
	cacheService := service.NewCacheService(*config)

	// Photo
	photoRepo := repository.NewPhotoRepository(dbEnv)
	photoService := service.NewPhotoService(photoRepo, albumService, filesystemService, cacheService, *config)

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
	}
}

func setupHttpHandlers(
	config *appconfig.AppConfig,
	appServices *AppServices) (*http.Router, error) {

	userHandler := http.NewUserHandler(appServices.userService)
	authHandler := http.NewAuthHandler(appServices.authService)
	photoHandler := http.NewPhotoHandler(appServices.photoService, appServices.jobService)
	albumHandler := http.NewAlbumHandler(appServices.albumService)
	utilityHandler := http.NewUtilityHandler(appServices.utilityService)
	searchHandler := http.NewSearchHandler(appServices.photoService, appServices.albumService)
	shareHandler := http.NewShareHandler(appServices.shareService)

	return http.NewRouter(
		*config,
		appServices.tokenService,
		*authHandler,
		*photoHandler,
		*albumHandler,
		*utilityHandler,
		*userHandler,
		*searchHandler,
		*shareHandler,
	)
}
