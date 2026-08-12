package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// HealthHandler exposes liveness, readiness, and startup probes for Kubernetes.
//
// Liveness (/health) performs no dependency checks so it is cheap and cannot be
// dragged down by a slow or unavailable database. Readiness (/ready) verifies
// critical dependencies and returns 503 while they are unreachable so the pod
// is removed from the Service. Startup (/healthz) gates the pod until
// initialization is complete.
type HealthHandler struct {
	dbPing     func() error
	storageDir string
	cachePing  func() error // optional; treated as non-blocking in readiness
	version    string
}

// NewHealthHandler creates a HealthHandler. cachePing may be nil when the cache
// is disabled, in which case the cache check is omitted from readiness.
func NewHealthHandler(dbPing func() error, storageDir string, cachePing func() error, version string) *HealthHandler {
	return &HealthHandler{
		dbPing:     dbPing,
		storageDir: storageDir,
		cachePing:  cachePing,
		version:    version,
	}
}

// Liveness reports whether the process is alive and can serve requests.
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": h.version,
	})
}

// Readiness reports whether critical dependencies are reachable. The cache is
// reported but treated as non-blocking to match the cache service's graceful
// degradation behaviour.
func (h *HealthHandler) Readiness(c *gin.Context) {
	checks := map[string]string{}
	ok := true

	if err := h.dbPing(); err != nil {
		checks["database"] = "down"
		ok = false
	} else {
		checks["database"] = "ok"
	}

	if err := checkStorageDir(h.storageDir); err != nil {
		checks["storage"] = "down"
		ok = false
	} else {
		checks["storage"] = "ok"
	}

	if h.cachePing != nil {
		if err := h.cachePing(); err != nil {
			checks["cache"] = "down"
		} else {
			checks["cache"] = "ok"
		}
	}

	status := http.StatusOK
	body := gin.H{"status": "ok", "checks": checks, "version": h.version}
	if !ok {
		status = http.StatusServiceUnavailable
		body = gin.H{"status": "unavailable", "checks": checks, "version": h.version}
	}
	c.JSON(status, body)
}

// Startup reports whether the process has finished initializing. Migrations and
// base setup complete before the router starts serving, so by the time this
// endpoint is reachable the service is ready; it delegates to Readiness.
func (h *HealthHandler) Startup(c *gin.Context) {
	h.Readiness(c)
}

func checkStorageDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}
	return nil
}
