package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type ShareHandler struct {
	shareSvc port.ShareService
}

func NewShareHandler(shareSvc port.ShareService) *ShareHandler {
	return &ShareHandler{
		shareSvc,
	}
}

type createShareRequest struct {
	ResourceType string  `json:"resourceType" binding:"required,oneof=photo album"`
	ResourceId   string  `json:"resourceId" binding:"required"`
	Expiry       *string `json:"expiry"`
	Password     *string `json:"password"`
}

func (sh *ShareHandler) CreateShare(ctx *gin.Context) {
	var req createShareRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payload := GetAuthPayload(ctx)
	if payload == nil {
		handleAbort(ctx, domain.ErrUnauthorized)
		return
	}

	share, err := sh.shareSvc.CreateShare(req.ResourceType, req.ResourceId, payload.ID.String(), req.Expiry, req.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}

	baseURL := ctx.Request.Host + "/api/shared/"
	resp := gin.H{
		"token": share.Token,
		"url":   baseURL + share.Token,
	}
	if share.Expiry != nil {
		resp["expiry"] = share.Expiry.Format("2006-01-02T15:04:05Z07:00")
	}

	handleSuccess(ctx, resp)
}

func (sh *ShareHandler) GetShared(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	var password *string
	if p := ctx.Query("password"); p != "" {
		password = &p
	}

	share, err := sh.shareSvc.GetSharedResource(token, password)
	if err != nil {
		if err == domain.ErrSharedLinkExpired {
			ctx.AbortWithStatusJSON(http.StatusGone, gin.H{"error": "shared link expired"})
			return
		}
		if err == domain.ErrSharedLinkPasswordRequired {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "password required"})
			return
		}
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, share)
}

func (sh *ShareHandler) ListShares(ctx *gin.Context) {
	shares, err := sh.shareSvc.ListShares()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, shares)
}

func (sh *ShareHandler) RevokeShare(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	err := sh.shareSvc.RevokeShare(token)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Share revoked"})
}
