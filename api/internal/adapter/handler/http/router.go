package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

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

	// Configure CORS with environment-based allowed origins
	config := cors.DefaultConfig()
	config.AllowOrigins = appConfig.CorsAllowedOrigins
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "UserResponse-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(requestIDMiddleware())
	router.Use(contentTypeMiddleware())

	// Add request timeout middleware (30 seconds for most requests)
	router.Use(timeoutMiddleware(30 * time.Second))

	// Create rate limiters
	// Global rate limiter: 100 requests per second with burst of 200
	globalLimiter := NewIPRateLimiter(100, 200)
	router.Use(rateLimitMiddleware(globalLimiter))

	// Strict rate limiter for auth endpoints: 5 requests per second with burst of 10
	authLimiter := NewIPRateLimiter(5, 10)

	defineResources(appConfig, router, token, authHandler, photoHandler, albumHandler, utilityHandler, userHandler, authLimiter)

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
	userHandler UserHandler,
	authLimiter *IPRateLimiter) {

	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	// Apply strict rate limiting to login endpoint to prevent brute force attacks
	router.POST(fmt.Sprintf("%s/login", urlBasePath), rateLimitMiddleware(authLimiter), authHandler.Login)

	user := router.Group(fmt.Sprintf("%s/user", urlBasePath))
	{
		// Apply strict rate limiting to registration endpoint
		user.POST("/register", rateLimitMiddleware(authLimiter), userHandler.Register)
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

	photo := router.Group(fmt.Sprintf("%s/photo", urlBasePath)).Use(authMiddleware(token))
	{
		photo.GET("/info/*id", photoHandler.GetPhoto)
		photo.GET("/thumbnail/*id", photoHandler.GetPhotoThumbnail)
		photo.GET("/bin/*id", photoHandler.GetPhotoBin)
	}

	photos := router.Group(fmt.Sprintf("%s/photos", urlBasePath)).Use(authMiddleware(token))
	{
		photos.GET("", photoHandler.ListPhotos)
		photos.GET("/count", photoHandler.GetPhotoCount)
	}

	// Photo indexing is admin-only as it's a system operation
	router.POST(fmt.Sprintf("%s/photos/index", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.IndexPhotos)

	settings := router.Group(urlBasePath).Use(authMiddleware(token))
	{
		settings.GET("/health", utilityHandler.HealthCheck)
	}

	// Settings management is admin-only
	router.GET(fmt.Sprintf("%s/settings", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), utilityHandler.ListAllSettings)
	router.POST(fmt.Sprintf("%s/settings", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), utilityHandler.UpdateSettings)

}
