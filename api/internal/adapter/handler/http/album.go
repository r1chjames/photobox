package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type AlbumHandler struct {
	svc port.AlbumService
}

// NewAlbumHandler creates a new AlbumHandler instance
func NewAlbumHandler(svc port.AlbumService) *AlbumHandler {
	return &AlbumHandler{
		svc,
	}
}

func (ah *AlbumHandler) GetAlbum(ctx *gin.Context) {
	albumId := ctx.Param("id")
	if albumId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "album ID is required"})
		return
	}

	resp, err := ah.svc.GetAlbumById(albumId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// Cross-tenant guard (issue #74).
	if !resourceInWorkspace(ctx, resp.WorkspaceID) {
		return
	}
	handleSuccess(ctx, resp)
}

func albumPaginationParams(resp []*domain.Album) (string, string, string) {
	if len(resp) > 0 {
		fromId := strconv.FormatInt(resp[0].CreatedEpoch, 10)
		toId := strconv.FormatInt(resp[len(resp)-1].CreatedEpoch, 10)
		return fromId, toId, "/api/albums?limit=10&fromId=%s"
	}
	return "", "", ""
}

func (ah *AlbumHandler) ListAlbums(ctx *gin.Context) {
	fromId := ctx.Query("fromId")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "30"))
	if limit <= 0 {
		limit = 30
	} else if limit > maxPageLimit {
		limit = maxPageLimit
	}

	resp, err := ah.svc.ListAlbums(fromId, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Scope to the caller's workspace (issue #74). Post-query filter is a
	// stopgap until the Phase 2 query-scoping sweep moves it into the query.
	wc := GetWorkspaceContext(ctx)
	if wc == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	resp = filterByWorkspace(resp, wc.WorkspaceID, func(a *domain.Album) string { return a.WorkspaceID })

	fromId, toId, nextPage := albumPaginationParams(resp)
	handlePaginatedSuccess(ctx, resp, fromId, toId, len(resp), nextPage)
}

func (ah *AlbumHandler) AlbumCount(ctx *gin.Context) {
	resp, err := ah.svc.AlbumCount()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, resp)
}

type createAlbumRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (ah *AlbumHandler) CreateAlbum(ctx *gin.Context) {
	var req createAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	album, err := ah.svc.CreateAlbum(req.Name)
	if err != nil {
		handleError(ctx, err)
		return
	}

	if req.Description != "" {
		album, err = ah.svc.UpdateAlbum(album.ID, map[string]any{"description": req.Description})
		if err != nil {
			handleError(ctx, err)
			return
		}
	}

	handleSuccess(ctx, album)
}

// createSmartAlbumRequest is the body for POST /albums/smart.
type createSmartAlbumRequest struct {
	Name  string                 `json:"name" binding:"required"`
	Rules domain.SmartAlbumRules `json:"rules"`
}

// CreateSmartAlbum creates a rule-based album (contents derived from filters).
func (ah *AlbumHandler) CreateSmartAlbum(ctx *gin.Context) {
	var req createSmartAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	album, err := ah.svc.CreateSmartAlbum(req.Name, req.Rules)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, album)
}

// updateSmartAlbumRequest is the body for PATCH /albums/:id/smart.
type updateSmartAlbumRequest struct {
	Rules domain.SmartAlbumRules `json:"rules"`
}

// UpdateSmartAlbum replaces a smart album's filter rules.
func (ah *AlbumHandler) UpdateSmartAlbum(ctx *gin.Context) {
	albumId := ctx.Param("id")
	if albumId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "album ID is required"})
		return
	}

	var req updateSmartAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	album, err := ah.svc.UpdateSmartAlbum(albumId, req.Rules)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, album)
}

type updateAlbumRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	CoverPhotoId string `json:"coverPhotoId"`
	Tags         string `json:"tags"`
}

func (ah *AlbumHandler) UpdateAlbum(ctx *gin.Context) {
	albumId := ctx.Param("id")
	if albumId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "album ID is required"})
		return
	}

	var req updateAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CoverPhotoId != "" {
		updates["coverPhotoId"] = req.CoverPhotoId
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}

	album, err := ah.svc.UpdateAlbum(albumId, updates)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, album)
}

func (ah *AlbumHandler) DeleteAlbum(ctx *gin.Context) {
	albumId := ctx.Param("id")
	if albumId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "album ID is required"})
		return
	}

	deletePhotos, _ := strconv.ParseBool(ctx.DefaultQuery("deletePhotos", "false"))

	// Cross-tenant guard (issue #74): verify ownership before deleting. The
	// album's workspace is immutable, so the check cannot be raced.
	existing, err := ah.svc.GetAlbumById(albumId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	if !resourceInWorkspace(ctx, existing.WorkspaceID) {
		return
	}

	err = ah.svc.DeleteAlbum(albumId, deletePhotos)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Album deleted"})
}
