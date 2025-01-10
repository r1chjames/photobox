package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// getAuthPayload is a helper function to get the auth payload from the context
func getAuthPayload(ctx *gin.Context, key string) *domain.TokenPayload {
	return ctx.MustGet(key).(*domain.TokenPayload)
}

// toMap is a helper function to add meta and data to a map
func toMap(m meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}
