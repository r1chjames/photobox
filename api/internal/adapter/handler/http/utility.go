package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type UtilityHandler struct {
	svc     port.UtilityService
	version string
}

// NewUtilityHandler creates a new UtilityHandler instance
func NewUtilityHandler(svc port.UtilityService, version string) *UtilityHandler {
	return &UtilityHandler{
		svc:     svc,
		version: version,
	}
}

func (uh *UtilityHandler) ListAllSettings(ctx *gin.Context) {
	resp, err := uh.svc.ListAllSettings()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, resp)
}

func (uh *UtilityHandler) UpdateSettings(c *gin.Context) {
	var settings []*domain.Setting
	if err := c.ShouldBindJSON(&settings); err != nil {
		validationError(c, err)
		return
	}

	if err := uh.svc.UpdateAllSettings(settings); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusAccepted)
}
