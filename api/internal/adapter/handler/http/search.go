package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type SearchHandler struct {
	photoSvc port.PhotoService
	albumSvc port.AlbumService
}

func NewSearchHandler(photoSvc port.PhotoService, albumSvc port.AlbumService) *SearchHandler {
	return &SearchHandler{
		photoSvc,
		albumSvc,
	}
}

type searchResult struct {
	ID   string      `json:"id"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (sh *SearchHandler) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "q query parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "30"))
	if limit <= 0 {
		limit = 30
	} else if limit > maxPageLimit {
		limit = maxPageLimit
	}

	// Scope results to the caller's workspace (issue #74). Search queries are
	// not yet workspace-filtered (Phase 2 sweep), so filter here to avoid
	// returning other tenants' photos and albums. NOTE: the limit is applied
	// before filtering, so a tenant can receive fewer than `limit` results —
	// the sweep removes both the filtering and the caveat.
	wc := GetWorkspaceContext(ctx)
	if wc == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	photos, err := sh.photoSvc.Search(query, limit)
	var results []searchResult
	if err == nil {
		for _, p := range filterByWorkspace(photos, wc.WorkspaceID, func(p *domain.Photo) string { return p.WorkspaceID }) {
			results = append(results, searchResult{ID: p.ID, Type: "photo", Data: p})
		}
	}

	albums, err := sh.albumSvc.ListAlbums("", limit)
	if err == nil {
		for _, a := range filterByWorkspace(albums, wc.WorkspaceID, func(a *domain.Album) string { return a.WorkspaceID }) {
			results = append(results, searchResult{ID: a.ID, Type: "album", Data: a})
		}
	}

	handleSuccess(ctx, results)
}
