package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
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

func albumPaginationParams(resp []*domain.Album) (string, string, string) {
	if (resp == nil) || (len(resp) == 0) {
		fromId := strconv.FormatInt(resp[0].CreatedEpoch, 10)
		toId := strconv.FormatInt(resp[len(resp)-1].CreatedEpoch, 10)
		return fromId, toId, "/api/albums?limit=10&fromId=%s"
	}
	return "", "", ""
}

func (ah *AlbumHandler) ListAlbums(ctx *gin.Context) {
	fromId := ctx.Query("fromId")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "30"))

	resp, err := ah.svc.ListAlbums(fromId, limit)
	handleError(ctx, err)
	fromId, toId, nextPage := albumPaginationParams(resp)

	handlePaginatedSuccess(ctx, resp, fromId, toId, len(resp), nextPage)
}

func (ah *AlbumHandler) AlbumCount(ctx *gin.Context) {
	resp, err := ah.svc.AlbumCount()
	handleError(ctx, err)
	handleSuccess(ctx, resp)
}
