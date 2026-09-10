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
	apiKeyHandler *ApiKeyHandler,
	importHandler *ImportHandler,
	wsHandler *WebSocketHandler,
	workspaceHandler *WorkspaceHandler,
	workspaceService port.WorkspaceService) (*Router, error) {

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
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma", "X-Workspace-ID"}
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

	defineResources(appConfig, router, token, authHandler, photoHandler, albumHandler, utilityHandler, healthHandler, userHandler, searchHandler, shareHandler, apiKeyHandler, importHandler, authLimiter, thumbnailLimiter, wsHandler, workspaceHandler, workspaceService)

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
	apiKeyHandler *ApiKeyHandler,
	importHandler *ImportHandler,
	authLimiter *IPRateLimiter,
	thumbnailLimiter *IPRateLimiter,
	wsHandler *WebSocketHandler,
	workspaceHandler *WorkspaceHandler,
	workspaceService port.WorkspaceService) {

	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	// Health endpoints — unauthenticated, for Kubernetes probes.
	health := router.Group(urlBasePath)
	{
		health.GET("/health", healthHandler.Liveness)
		health.GET("/ready", healthHandler.Readiness)
		health.GET("/healthz", healthHandler.Startup)
	}

	// WebSocket endpoint — upgrades after auth. The workspace middleware
	// binds the connection to a workspace (via the workspace_id query param,
	// since browsers cannot set headers on a WS upgrade); the hub then
	// delivers only that workspace's events (issue #74).
	router.GET(fmt.Sprintf("%s/ws", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), wsHandler.HandleUpgrade)

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

	// Album get. Previously unauthenticated (a public read of any album by
	// ID — a cross-tenant exposure once tenancy exists, plan D3). The webapp
	// already sends a bearer token on this path, so requiring auth is not a
	// client-visible break.
	router.GET(fmt.Sprintf("%s/album/:id", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), albumHandler.GetAlbum)

	albums := router.Group(fmt.Sprintf("%s/albums", urlBasePath)).Use(authMiddleware(token), workspaceMiddleware(workspaceService))
	{
		albums.GET("", albumHandler.ListAlbums)
		albums.GET("/count", albumHandler.AlbumCount)
		albums.POST("", albumHandler.CreateAlbum)
		albums.POST("/smart", albumHandler.CreateSmartAlbum)
		albums.PATCH("/:id/smart", albumHandler.UpdateSmartAlbum)
		albums.POST("/:id/photos", photoHandler.AddPhotosToAlbum)
		albums.PATCH("/:id", albumHandler.UpdateAlbum)
		albums.DELETE("/:id", albumHandler.DeleteAlbum)
	}

	photo := router.Group(fmt.Sprintf("%s/photo", urlBasePath)).Use(authMiddleware(token), workspaceMiddleware(workspaceService))
	{
		photo.GET("/info/*id", photoHandler.GetPhoto)
		photo.GET("/bin/*id", photoHandler.GetPhotoBin)
		photo.POST("", photoHandler.UploadPhoto)
	}

	// Thumbnail endpoint with dedicated stricter rate limiter
	router.GET(fmt.Sprintf("%s/photo/thumbnail/*id", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), rateLimitMiddleware(thumbnailLimiter), photoHandler.GetPhotoThumbnail)

	// Capability-URL thumbnails (issue #74 D1): unauthenticated by design —
	// the random thumb_cap is the credential. Immutable-cacheable so a CDN
	// edge can absorb thumbnail reads.
	router.GET(fmt.Sprintf("%s/t/:cap/:size", urlBasePath), photoHandler.GetThumbnailByCap)

	photos := router.Group(fmt.Sprintf("%s/photos", urlBasePath)).Use(authMiddleware(token), workspaceMiddleware(workspaceService))
	{
		photos.GET("", photoHandler.ListPhotos)
		photos.GET("/count", photoHandler.GetPhotoCount)
		photos.GET("/trash", photoHandler.ListTrashPhotos)
		photos.DELETE("/trash/empty", photoHandler.EmptyTrash)
		photos.POST("/trash/restore/:id", photoHandler.RestorePhoto)
		photos.POST("/download", photoHandler.DownloadPhotos)
	photos.POST("/batch", photoHandler.BatchPhotos)
		photos.GET("/timeline", photoHandler.GetTimeline)
		photos.GET("/memories", photoHandler.GetMemories)
	photos.GET("/geodata", photoHandler.GetGeodata)
	photos.GET("/tags", photoHandler.GetAllTags)
	photos.POST("/tags/batch", photoHandler.BatchUpdatePhotoTags)
	photos.GET("/duplicates", photoHandler.GetDuplicatePhotos)
	}

	// Individual photo actions
	router.DELETE(fmt.Sprintf("%s/photos/:id", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.DeletePhoto)
	router.PATCH(fmt.Sprintf("%s/photos/:id/favorite", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.SetFavorite)
	router.PATCH(fmt.Sprintf("%s/photos/:id/tags", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.UpdatePhotoTags)
	router.PATCH(fmt.Sprintf("%s/photos/:id/metadata", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.UpdatePhotoMetadata)
	router.GET(fmt.Sprintf("%s/photos/:id/live-video", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), rateLimitMiddleware(thumbnailLimiter), photoHandler.GetPhotoLiveVideo)
	router.GET(fmt.Sprintf("%s/photos/:id/location", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.GetPhotoLocation)
	router.POST(fmt.Sprintf("%s/photos/:id/edit", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.EditPhoto)
	router.DELETE(fmt.Sprintf("%s/photos/:id/edit", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.ClearEdits)
	router.POST(fmt.Sprintf("%s/photos/:id/rotate", urlBasePath), authMiddleware(token), workspaceMiddleware(workspaceService), photoHandler.RotatePhoto)

	// Photo system jobs (admin-only)
	router.POST(fmt.Sprintf("%s/photos/jobs/:type/start", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.StartJob)
	router.POST(fmt.Sprintf("%s/photos/jobs/:type/stop", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.StopJob)
	router.GET(fmt.Sprintf("%s/photos/jobs", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.GetJobStatuses)
	router.POST(fmt.Sprintf("%s/photos/jobs/stop", urlBasePath), authMiddleware(token), requireRole(domain.ADMINISTRATOR), photoHandler.StopAllJobs)

	// Workspaces (issue #74). Listing/creating workspaces needs only auth —
	// the middleware is not applied because these routes operate across the
	// caller's workspaces, not within one.
	workspaces := router.Group(fmt.Sprintf("%s/workspaces", urlBasePath)).Use(authMiddleware(token))
	{
		workspaces.GET("", workspaceHandler.List)
		workspaces.POST("", workspaceHandler.Create)
		workspaces.GET("/:id", workspaceHandler.Get)
		workspaces.PATCH("/:id", workspaceHandler.Update)
		workspaces.DELETE("/:id", workspaceHandler.Delete)
		workspaces.GET("/:id/members", workspaceHandler.ListMembers)
		workspaces.POST("/:id/members", workspaceHandler.AddMember)
		workspaces.DELETE("/:id/members/:userId", workspaceHandler.RemoveMember)
		workspaces.PATCH("/:id/members/:userId", workspaceHandler.UpdateMemberRole)
	}

	// Search (workspace middleware is applied together with query scoping)
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

	// API key management (admin-only, issue #112)
	apiKeys := router.Group(fmt.Sprintf("%s/api-keys", urlBasePath)).Use(authMiddleware(token), requireRole(domain.ADMINISTRATOR))
	{
		apiKeys.POST("", apiKeyHandler.CreateApiKey)
		apiKeys.GET("", apiKeyHandler.ListApiKeys)
		apiKeys.DELETE("/:id", apiKeyHandler.RevokeApiKey)
	}

	// Google Takeout import (admin-only, issue #106)
	importGroup := router.Group(fmt.Sprintf("%s/import", urlBasePath)).Use(authMiddleware(token), requireRole(domain.ADMINISTRATOR))
	{
		importGroup.POST("/takeout", importHandler.ImportTakeout)
		importGroup.POST("/takeout/scan", importHandler.ScanTakeout)
		importGroup.GET("/takeout/progress", importHandler.GetProgress)
	}

}
