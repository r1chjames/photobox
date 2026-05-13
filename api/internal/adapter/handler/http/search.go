package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	photos, err := sh.photoSvc.Search(query, limit)
	var results []searchResult
	if err == nil {
		for _, p := range photos {
			results = append(results, searchResult{ID: p.ID, Type: "photo", Data: p})
		}
	}

	albums, err := sh.albumSvc.ListAlbums("", limit)
	if err == nil {
		for _, a := range albums {
			results = append(results, searchResult{ID: a.ID, Type: "album", Data: a})
		}
	}

	handleSuccess(ctx, results)
}
