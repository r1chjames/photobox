package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
)

// ImportHandler handles HTTP requests for Google Takeout import operations.
type ImportHandler struct {
	importer *service.TakeoutImporter
}

// NewImportHandler creates a new ImportHandler instance.
func NewImportHandler(importer *service.TakeoutImporter) *ImportHandler {
	return &ImportHandler{importer: importer}
}

// ScanTakeout handles POST /api/import/takeout/scan — inspects an export
// without modifying the library and returns preview totals.
func (ih *ImportHandler) ScanTakeout(ctx *gin.Context) {
	var req domain.TakeoutImportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	if req.Path == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	result, err := ih.importer.ScanTakeout(req.Path)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// ImportTakeout handles POST /api/import/takeout — starts a background import.
// Returns 202 Accepted immediately; progress is available via the progress
// endpoint and WebSocket events.
func (ih *ImportHandler) ImportTakeout(ctx *gin.Context) {
	var req domain.TakeoutImportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	if req.Path == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	go func() {
		_ = ih.importer.ImportTakeout(context.Background(), req)
	}()

	ctx.JSON(http.StatusAccepted, gin.H{"message": "import started"})
}

// GetProgress handles GET /api/import/takeout/progress — returns the current
// or last import progress snapshot.
func (ih *ImportHandler) GetProgress(ctx *gin.Context) {
	progress := ih.importer.GetProgress()
	ctx.JSON(http.StatusOK, progress)
}
