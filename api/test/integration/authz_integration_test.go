//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/auth"
	httpHandler "gitlab.com/r1chjames/photobox/api/internal/adapter/handler/http"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/repository"
	filesystemRepos "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/filesystem/repository"
	thumbFs "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/thumbnail/filesystem"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	ws "gitlab.com/r1chjames/photobox/api/internal/components/websocket"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"gitlab.com/r1chjames/photobox/api/test/testutil"
)

// buildAuthzRouter wires the real HTTP router against the test DB so the
// AuthZ suite exercises the actual middleware + handlers (issue #152).
func buildAuthzRouter(t *testing.T, env *database.Env, config *appconfig.AppConfig) *httpHandler.Router {
	t.Helper()

	token, err := auth.New(config.Token, config.TokenDuration)
	assert.NoError(t, err)

	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, token)

	albumRepo := repository.NewAlbumRepository(env)
	albumService := service.NewAlbumService(albumRepo, *config)

	utilityRepo := repository.NewUtilityRepository(env)
	utilityService := service.NewUtilityService(utilityRepo)

	jobRepo := repository.NewJobRepository(env)
	jobService := service.NewJobService(jobRepo)

	filesystemRepo := filesystemRepos.NewFilesystemRepository(*config, jobService)
	filesystemService := service.NewFilesystemService(filesystemRepo, jobService, utilityService, config.PhotoIndexWorkers, config.ThumbnailStorage)

	cacheService := service.NewCacheService(*config)
	wsHub := ws.NewHub(context.Background())

	photoRepo := repository.NewPhotoRepository(env)
	thumbnailStorage := thumbFs.New(config.ThumbnailDir)
	photoService := service.NewPhotoService(photoRepo, albumService, filesystemService, cacheService, nil, *config, thumbnailStorage, wsHub, jobService)

	shareRepo := repository.NewShareRepository(env)
	shareService := service.NewShareService(shareRepo, photoService, albumService)

	apiKeyRepo := repository.NewApiKeyRepository(env)
	apiKeyService := service.NewApiKeyService(apiKeyRepo)
	httpHandler.SetApiKeyService(apiKeyService)

	userHandler := httpHandler.NewUserHandler(userService)
	authHandler := httpHandler.NewAuthHandler(authService, config.TokenDuration)
	photoHandler := httpHandler.NewPhotoHandler(photoService, jobService)
	albumHandler := httpHandler.NewAlbumHandler(albumService)
	utilityHandler := httpHandler.NewUtilityHandler(utilityService, "test")
	healthHandler := httpHandler.NewHealthHandler(utilityService.Ping, config.PhotoDir, cacheService.Ping, "test")
	searchHandler := httpHandler.NewSearchHandler(photoService, albumService)
	shareHandler := httpHandler.NewShareHandler(shareService)
	apiKeyHandler := httpHandler.NewApiKeyHandler(apiKeyService)
	importHandler := httpHandler.NewImportHandler(service.NewTakeoutImporter(photoService, filesystemService, photoRepo, config, wsHub))
	wsHandler := httpHandler.NewWebSocketHandler(wsHub)

	router, err := httpHandler.NewRouter(
		*config,
		token,
		*authHandler,
		photoHandler,
		*albumHandler,
		*utilityHandler,
		healthHandler,
		*userHandler,
		*searchHandler,
		*shareHandler,
		apiKeyHandler,
		importHandler,
		wsHandler,
	)
	assert.NoError(t, err)
	return router
}

// loginAs logs in and returns a bearer token for the given user.
func loginAs(t *testing.T, router *httpHandler.Router, username, password string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "login should succeed for %s", username)
	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp.Data.Token
}

// TestAuthZ_Unauthenticated_401 asserts protected endpoints reject anonymous
// requests (issue #152).
func TestAuthZ_Unauthenticated_401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	timezone, _ := time.LoadLocation("UTC")
	config := &appconfig.AppConfig{
		PhotoDir:         t.TempDir(),
		Timezone:         timezone,
		Token:            "0123456789abcdef0123456789abcdef",
		TokenDuration:    time.Hour,
		ThumbnailDir:     t.TempDir(),
		ThumbnailStorage: "filesystem",
		CorsAllowedOrigins: []string{"http://localhost"},
		ApiBasePath: "/api",
	}
	router := buildAuthzRouter(t, env, config)

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/photos"},
		{http.MethodGet, "/api/albums"},
		{http.MethodGet, "/api/shares"},
		{http.MethodGet, "/api/users"},
		{http.MethodGet, "/api/settings"},
	}
	for _, ep := range endpoints {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(ep.method, ep.path, nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "%s %s should be 401 without auth", ep.method, ep.path)
	}
}

// TestAuthZ_AdminOnly_403 asserts non-admin users are denied admin-only
// routes (issue #152).
func TestAuthZ_AdminOnly_403(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	timezone, _ := time.LoadLocation("UTC")
	config := &appconfig.AppConfig{
		PhotoDir:         t.TempDir(),
		Timezone:         timezone,
		Token:            "0123456789abcdef0123456789abcdef",
		TokenDuration:    time.Hour,
		ThumbnailDir:     t.TempDir(),
		ThumbnailStorage: "filesystem",
		CorsAllowedOrigins: []string{"http://localhost"},
		ApiBasePath: "/api",
	}
	router := buildAuthzRouter(t, env, config)

	// Create a non-admin user (service hashes the password so login works)
	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	_, err := userService.CreateUser(&domain.User{
		Username: "viewer1",
		Email:    "viewer1@example.com",
		Password: "password123",
		Role:     domain.VIEWER,
		Approved: true,
	})
	assert.NoError(t, err)

	token := loginAs(t, router, "viewer1", "password123")

	// Non-admin hitting admin-only /api/shares should get 403
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/shares", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code, "non-admin should get 403 on /api/shares")
}

// TestAuthZ_ShareOwnerScoping_404 asserts a user cannot revoke another
// user's share (issue #151/#152).
func TestAuthZ_ShareOwnerScoping_404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	userRepo := repository.NewUserRepository(env)
	userService := service.NewUserService(userRepo)
	userA, err := userService.CreateUser(&domain.User{Username: "usera", Email: "a@example.com", Password: "password123", Role: domain.CONTRIBUTOR, Approved: true})
	assert.NoError(t, err)
	userB, err := userService.CreateUser(&domain.User{Username: "userb", Email: "b@example.com", Password: "password123", Role: domain.CONTRIBUTOR, Approved: true})
	assert.NoError(t, err)

	// User A creates a share
	shareRepo := repository.NewShareRepository(env)
	share := &domain.SharedLink{Token: "token-a", ResourceType: "photo", ResourceId: "p1", CreatedBy: userA.ID}
	assert.NoError(t, shareRepo.CreateShare(share))

	// User B (non-admin) tries to revoke A's share via the service layer
	// (the route is admin-only, so this exercises the service scoping).
	shareService := service.NewShareService(shareRepo, nil, nil)
	err = shareService.RevokeShare("token-a", userB.ID)
	assert.ErrorIs(t, err, domain.ErrDataNotFound, "user B must not revoke user A's share")

	// User A can revoke their own share
	err = shareService.RevokeShare("token-a", userA.ID)
	assert.NoError(t, err)
}

// TestAuthZ_ApiKeyAuth verifies API key authentication (issue #112): a key
// created via the service can access protected endpoints via X-API-Key.
func TestAuthZ_ApiKeyAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	timezone, _ := time.LoadLocation("UTC")
	config := &appconfig.AppConfig{
		PhotoDir:         t.TempDir(),
		Timezone:         timezone,
		Token:            "0123456789abcdef0123456789abcdef",
		TokenDuration:    time.Hour,
		ThumbnailDir:     t.TempDir(),
		ThumbnailStorage: "filesystem",
		CorsAllowedOrigins: []string{"http://localhost"},
		ApiBasePath: "/api",
	}
	router := buildAuthzRouter(t, env, config)

	// Create an API key directly via the service
	apiKeyRepo := repository.NewApiKeyRepository(env)
	apiKeySvc := service.NewApiKeyService(apiKeyRepo)
	key, _, err := apiKeySvc.CreateKey("test-key", domain.ApiKeyReadOnly, "user-1")
	assert.NoError(t, err)

	// Access a protected endpoint with the API key
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/photos", nil)
	req.Header.Set("X-API-Key", key)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "valid API key should access protected endpoint")

	// Invalid API key is rejected
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/photos", nil)
	req2.Header.Set("X-API-Key", "pb_invalid")
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code, "invalid API key should be rejected")
}
