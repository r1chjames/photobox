package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)
import "net/http"

type UtilityHandler struct {
	svc port.UtilityService
}

// NewUtilityHandler creates a new UtilityHandler instance
func NewUtilityHandler(svc port.UtilityService) *UtilityHandler {
	return &UtilityHandler{
		svc,
	}
}

func (uh *UtilityHandler) HealthCheck(ctx *gin.Context) {
	resp, err := uh.svc.Healthcheck()
	handleError(ctx, err)
	handleSuccess(ctx, resp)
}

func (uh *UtilityHandler) ListAllSettings(ctx *gin.Context) {
	resp, err := uh.svc.ListAllSettings()
	handleError(ctx, err)
	handleSuccess(ctx, resp)
}

func (uh *UtilityHandler) UpdateSettings(c *gin.Context) {
	var settings []*domain.Setting
	err := c.BindJSON(&settings)

	err = uh.svc.UpdateAllSettings(settings)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		c.Status(http.StatusAccepted)
	}
}
