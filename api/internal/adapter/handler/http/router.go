package http

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"strings"
	"time"
)

var dbEnv *database.Env

func NewDBEnv(env *database.Env) *database.Env {
	dbEnv = env
	return dbEnv
}

type Router struct {
	*gin.Engine
}

func NewRouter(
	appConfig appconfig.AppConfig,
	token port.TokenService,
	authHandler AuthHandler,
	photoHandler PhotoHandler,
	albumHandler AlbumHandler,
	utilityHandler UtilityHandler,
	userHandler UserHandler) (*Router, error) {

	router := gin.Default()
	//router.Use(cors.Default())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "UserResponse-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	defineResources(appConfig, router, token, authHandler, photoHandler, albumHandler, utilityHandler, userHandler)

	return &Router{
		router,
	}, nil
}

func defineResources(
	appConfig appconfig.AppConfig,
	router *gin.Engine,
	token port.TokenService,
	authHandler AuthHandler,
	photoHandler PhotoHandler,
	albumHandler AlbumHandler,
	utilityHandler UtilityHandler,
	userHandler UserHandler) {

	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	router.POST(fmt.Sprintf("%s/login", urlBasePath), authHandler.Login)

	user := router.Group(fmt.Sprintf("%s/user", urlBasePath))
	{
		user.POST("/register", userHandler.Register)
		authUser := user.Use(authMiddleware(token))
		{
			authUser.POST("/update", userHandler.UpdateUser)
		}
	}

	router.GET(fmt.Sprintf("%s/album/:id", urlBasePath), albumHandler.GetAlbum)
	albums := router.Group(fmt.Sprintf("%s/albums", urlBasePath)).Use(authMiddleware(token))
	{
		albums.GET("", albumHandler.ListAlbums)
		albums.GET("/count", albumHandler.AlbumCount)
	}

	photos := router.Group(fmt.Sprintf("%s/photos", urlBasePath)).Use(authMiddleware(token))
	{
		photos.GET("", photoHandler.ListPhotos)
		photos.GET("/count", photoHandler.GetPhotoCount)
		photos.POST("/index", photoHandler.IndexPhotos)
	}

	photo := router.Group(fmt.Sprintf("%s/photo", urlBasePath)).Use(authMiddleware(token))
	{
		//photo.POST("", addPhoto)
		photo.GET("/:id", photoHandler.GetPhoto)
		photo.GET("/:id/thumbnail", photoHandler.GetPhotoThumbnail)
		photo.GET("/:id/bin", photoHandler.GetPhotoBin)
	}

	settings := router.Group(urlBasePath).Use(authMiddleware(token))
	{
		settings.GET("/health", utilityHandler.HealthCheck)
		settings.GET("/settings", utilityHandler.ListAllSettings)
		settings.POST("/settings", utilityHandler.UpdateSettings)
	}

}

func setDbEnv(env *database.Env) {
	dbEnv = env
}
