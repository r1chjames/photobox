package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performHealthRequest(t *testing.T, h *HealthHandler, method string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api"+method, nil)
	switch method {
	case "/health":
		h.Liveness(ctx)
	case "/ready":
		h.Readiness(ctx)
	case "/healthz":
		h.Startup(ctx)
	default:
		t.Fatalf("unknown health method %q", method)
	}
	return w
}

type healthBody struct {
	Status  string            `json:"status"`
	Version string            `json:"version"`
	Checks  map[string]string `json:"checks"`
}

func TestHealthHandler_Liveness(t *testing.T) {
	h := NewHealthHandler(func() error { return nil }, "/tmp", nil, "1.2.3")

	w := performHealthRequest(t, h, "/health")

	assert.Equal(t, http.StatusOK, w.Code)

	var body healthBody
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body.Status)
	assert.Equal(t, "1.2.3", body.Version)
	assert.Nil(t, body.Checks, "liveness must not perform dependency checks")
}

func TestHealthHandler_Readiness_OK(t *testing.T) {
	dir := t.TempDir()
	h := NewHealthHandler(func() error { return nil }, dir, nil, "1.2.3")

	w := performHealthRequest(t, h, "/ready")

	assert.Equal(t, http.StatusOK, w.Code)

	var body healthBody
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body.Status)
	assert.Equal(t, "ok", body.Checks["database"])
	assert.Equal(t, "ok", body.Checks["storage"])
}

func TestHealthHandler_Readiness_DatabaseDown(t *testing.T) {
	dir := t.TempDir()
	h := NewHealthHandler(func() error { return errors.New("connection refused") }, dir, nil, "1.2.3")

	w := performHealthRequest(t, h, "/ready")

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var body healthBody
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "unavailable", body.Status)
	assert.Equal(t, "down", body.Checks["database"])
	assert.Equal(t, "ok", body.Checks["storage"])
}

func TestHealthHandler_Readiness_StorageMissing(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "does-not-exist")
	h := NewHealthHandler(func() error { return nil }, dir, nil, "1.2.3")

	w := performHealthRequest(t, h, "/ready")

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var body healthBody
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "unavailable", body.Status)
	assert.Equal(t, "ok", body.Checks["database"])
	assert.Equal(t, "down", body.Checks["storage"])
}

func TestHealthHandler_Readiness_CacheNonBlocking(t *testing.T) {
	dir := t.TempDir()
	// Cache is down, but that must not fail readiness.
	h := NewHealthHandler(func() error { return nil }, dir, func() error { return errors.New("cache unreachable") }, "1.2.3")

	w := performHealthRequest(t, h, "/ready")

	assert.Equal(t, http.StatusOK, w.Code)

	var body healthBody
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body.Status)
	assert.Equal(t, "down", body.Checks["cache"])
}

func TestHealthHandler_Readiness_CacheOK(t *testing.T) {
	dir := t.TempDir()
	h := NewHealthHandler(func() error { return nil }, dir, func() error { return nil }, "1.2.3")

	w := performHealthRequest(t, h, "/ready")

	assert.Equal(t, http.StatusOK, w.Code)

	var body healthBody
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body.Checks["cache"])
}

func TestHealthHandler_Startup_DelegatesToReadiness(t *testing.T) {
	dir := t.TempDir()
	h := NewHealthHandler(func() error { return nil }, dir, nil, "1.2.3")

	w := performHealthRequest(t, h, "/healthz")

	assert.Equal(t, http.StatusOK, w.Code)

	var body healthBody
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "ok", body.Status)
	assert.Equal(t, "ok", body.Checks["database"])
}

func TestCheckStorageDir(t *testing.T) {
	dir := t.TempDir()
	assert.NoError(t, checkStorageDir(dir))

	// A regular file is not a valid storage directory.
	f := filepath.Join(t.TempDir(), "file")
	assert.NoError(t, os.WriteFile(f, []byte("x"), 0o644))
	assert.Error(t, checkStorageDir(f))

	// A missing path errors.
	assert.Error(t, checkStorageDir(filepath.Join(t.TempDir(), "missing")))
}
