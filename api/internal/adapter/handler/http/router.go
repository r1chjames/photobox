package http

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
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
	photoHandler *PhotoHandler,
	albumHandler AlbumHandler,
	utilityHandler UtilityHandler,
	healthHandler *HealthHandler,
	userHandler UserHandler,
	searchHandler SearchHandler,
	shareHandler ShareHandler,
	wsHandler *WebSocketHandler) (*Router, error) {

	router := gin.New()
	router.MaxMultipartMemory = 32 << 20 // 32 MB

	// Custom panic recovery with structured logging and JSON error response
	router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		slog.Error("panic recovered", "error", err, "path", c.Request.URL.Path)
		c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
	}))

	// Gin's default logger
	router.Use(gin.Logger())

	// Configure CORS with environment-based allowed origins
	config := cors.DefaultConfig()
	config.AllowOrigins = appConfig.CorsAllowedOrigins
	config.AllowMethods = []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(requestIDMiddleware())
	router.Use(contentTypeMiddleware())

	// Add request timeout middleware (30 seconds for most requests)
	router.Use(timeoutMiddleware(30 * time.Second))

	// Create rate limiters
	// Global rate limiter: 100 requests per second with burst of 200
	globalLimiter := NewIPRateLimiter(100, 200)
	router.Use(rateLimitMiddleware(globalLimiter))

	// Add security headers to all responses
	router.Use(securityHeadersMiddleware())

	// Limit request body size to 50 MB to prevent OOM attacks
	router.Use(maxBodySizeMiddleware(50 * 1024 * 1024))

	// Strict rate limiter for auth endpoints: 5 requests per second with burst of 10
	authLimiter := NewIPRateLimiter(5, 10)

	// Dedicated rate limiter for thumbnail endpoints: 30 req/s with burst of 60
	thumbnailLimiter := NewIPRateLimiter(30, 60)

	defineResources(appConfig, router, token, authHandler, photoHandler, albumHandler, utilityHandler, healthHandler, userHandler, searchHandler, shareHandler, authLimiter, thumbnailLimiter, wsHandler)

	return &Router{
		router,
	}, nil
}

func defineResources(
	appConfig appconfig.AppConfig,
	router *gin.Engine,
	token port.TokenService,
	authHandler AuthHandler,
	photoHandler *PhotoHandler,
	albumHandler AlbumHandler,
	utilityHandler UtilityHandler,
	healthHandler *HealthHandler,
	userHandler UserHandler,
	searchHandler SearchHandler,
	shareHandler ShareHandler,
	authLimiter *IPRateLimiter,
	thumbnailLimiter *IPRateLimiter,
	wsHandler *WebSocketHandler) {

	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	// Health endpoints — unauthenticated, for Kubernetes probes.
	health := router.Group(urlBasePath)
	{
		health.GET("/health", healthHandler.Liveness)
		health.GET("/ready", healthHandler.Readiness)
		health.GET("/healthz", healthHandler.Startup)
	}

	// WebSocket endpoint — upgrades after auth
	router.GET(fmt.Sprintf("%s/ws", urlBasePath), authMiddleware(token), wsHandler.HandleUpgrade)

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

	// Admin user management
	users := router.Group(fmt.Sprintf("%s/users", urlBasePath)).Use(authMiddleware(token), requireRole(domain.ADMINISTRATOR))
	{
		users.GET("", userHandler.ListUsers)
	users.GET("/:id", userHandler.GetUser)
	users.PATCH("/:id", userHandler.UpdateUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
	}

	// Public album get (no auth)
	router.GET(fmt.Sprintf("%s/album/:id", urlBasePath), albumHandler.GetAlbum)

	albums := router.Group(fmt.Sprintf("%s/albums", urlBasePath)).Use(authMiddleware(token))
	{
		albums.GET("", albumHandler.ListAlbums)
		albums.GET("/count", albumHandler.AlbumCount)
		albums.POST("", albumHandler.CreateAlbum)
		albums.PATCH("/:id", albumHandler.UpdateAlbum)
		albums.DELETE("/:id", albumHandler.DeleteAlbum)
	}

	photo := router.Group(fmt.Sprintf("%s/photo", urlBasePath)).Use(authMiddleware(token))
	{
		photo.GET("/info/*id", photoHandler.GetPhoto)
		photo.GET("/bin/*id", photoHandler.GetPhotoBin)
	}

	// Thumbnail endpoint with dedicated stricter rate limiter
	router.GET(fmt.Sprintf("%s/photo/thumbnail/*id", urlBasePath), authMiddleware(token), rateLimitMiddleware(thumbnailLimiter), photoHandler.GetPhotoThumbnail)

	photos := router.Group(fmt.Sprintf("%s/photos", urlBasePath)).Use(authMiddleware(token))
	{
		photos.GET("", photoHandler.ListPhotos)
		photos.GET("/count", photoHandler.GetPhotoCount)
		photos.GET("/trash", photoHandler.ListTrashPhotos)
		photos.DELETE("/trash/empty", photoHandler.EmptyTrash)
		photos.POST("/trash/restore/:id", photoHandler.RestorePhoto)
		photos.POST("/download", photoHandler.DownloadPhotos)
		photos.GET("/timeline", photoHandler.GetTimeline)
	photos.GET("/geodata", photoHandler.GetGeodata)
	photos.GET("/tags", photoHandler.GetAllTags)
	photos.POST("/tags/batch", photoHandler.BatchUpdatePhotoTags)
	photos.GET("/duplicates", photoHandler.GetDuplicatePhotos)
	}

	// Individual photo actions
	router.DELETE(fmt.Sprintf("%s/photos/:id", urlBasePath), authMiddleware(token), photoHandler.DeletePhoto)
	router.PATCH(fmt.Sprintf("%s/photos/:id/favorite", urlBasePath), authMiddleware(token), photoHandler.SetFavorite)
	router.PATCH(fmt.Sprintf("%s/photos/:id/tags", urlBasePath), authMiddleware(token), photoHandler.UpdatePhotoTags)
	router.POST(fmt.Sprintf("%s/photos/:id/rotate", urlBasePath), authMiddleware(token), photoHandler.RotatePhoto)

	// Photo system jobs (admin-only)
	router.POST(fmt.Sprintf("%s/photos/jobs/:type/start", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.StartJob)
	router.POST(fmt.Sprintf("%s/photos/jobs/:type/stop", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.StopJob)
	router.GET(fmt.Sprintf("%s/photos/jobs", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.GetJobStatuses)
	router.POST(fmt.Sprintf("%s/photos/jobs/stop", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.StopAllJobs)

	// Search
	search := router.Group(fmt.Sprintf("%s/search", urlBasePath)).Use(authMiddleware(token))
	{
		search.GET("", searchHandler.Search)
	}

	// Sharing
	share := router.Group(fmt.Sprintf("%s/share", urlBasePath)).Use(authMiddleware(token))
	{
		share.POST("", shareHandler.CreateShare)
	}

	// Public shared view (no auth)
	router.GET(fmt.Sprintf("%s/shared/:token", urlBasePath), shareHandler.GetShared)
	router.GET(fmt.Sprintf("%s/shared/:token/resource", urlBasePath), shareHandler.GetSharedResourceData)

	// Admin share management
	shares := router.Group(fmt.Sprintf("%s/shares", urlBasePath)).Use(authMiddleware(token), requireRole(domain.ADMINISTRATOR))
	{
		shares.GET("", shareHandler.ListShares)
		shares.DELETE("/:token", shareHandler.RevokeShare)
	}

	// Settings management is admin-only
	router.GET(fmt.Sprintf("%s/settings", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), utilityHandler.ListAllSettings)
	router.POST(fmt.Sprintf("%s/settings", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), utilityHandler.UpdateSettings)

}
