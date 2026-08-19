package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type ApiKeyHandler struct {
	svc port.ApiKeyService
}

func NewApiKeyHandler(svc port.ApiKeyService) *ApiKeyHandler {
	return &ApiKeyHandler{svc}
}

type createApiKeyRequest struct {
	Name  string           `json:"name" binding:"required"`
	Scope domain.ApiKeyScope `json:"scope" binding:"required,oneof=read_only read_write admin"`
}

// CreateApiKey generates a new API key. The plaintext key is returned once.
func (h *ApiKeyHandler) CreateApiKey(ctx *gin.Context) {
	var req createApiKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payload := GetAuthPayload(ctx)
	createdBy := ""
	if payload != nil {
		createdBy = payload.ID.String()
	}

	key, apiKey, err := h.svc.CreateKey(req.Name, req.Scope, createdBy)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{
		"key":       key,
		"id":        apiKey.ID,
		"name":      apiKey.Name,
		"scope":     apiKey.Scope,
		"createdAt": apiKey.CreatedAt,
	})
}

// ListApiKeys returns all API keys (hashes only, never plaintext).
func (h *ApiKeyHandler) ListApiKeys(ctx *gin.Context) {
	keys, err := h.svc.ListKeys()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, keys)
}

// RevokeApiKey revokes an API key by ID.
func (h *ApiKeyHandler) RevokeApiKey(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "key ID is required"})
		return
	}
	if err := h.svc.RevokeKey(id); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, gin.H{"message": "API key revoked"})
}
