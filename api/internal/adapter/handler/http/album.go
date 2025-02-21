package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"net/http"
	"strconv"
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
		ctx.IndentedJSON(http.StatusBadRequest, &apiError{http.StatusNotFound, missingQueryParam("album ID")})
	}

	resp, err := ah.svc.GetAlbumById(albumId)
	handleError(ctx, err)
	handleSuccess(ctx, resp)
}

func (ah *AlbumHandler) ListAlbums(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "30"))

	resp, err := ah.svc.ListAlbums(page, limit)
	handleError(ctx, err)
	handlePaginatedSuccess(ctx, resp, resp[0].ID, resp[len(resp)-1].ID, len(resp))
}

func (ah *AlbumHandler) AlbumCount(ctx *gin.Context) {
	resp, err := ah.svc.AlbumCount()
	handleError(ctx, err)
	handleSuccess(ctx, resp)
}
