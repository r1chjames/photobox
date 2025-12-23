package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

// PhotoHandler represents the HTTP handler for photo-related requests
type PhotoHandler struct {
	photoSvc port.PhotoService
	jobSvc   port.JobService
}

// NewPhotoHandler creates a new PhotoHandler instance
func NewPhotoHandler(photoSvc port.PhotoService, jobSvc port.JobService) *PhotoHandler {
	return &PhotoHandler{
		photoSvc,
		jobSvc,
	}
}

func photosPaginationParams(resp []*domain.Photo) (string, string, string) {
	if resp != nil || len(resp) != 0 {
		fromId := strconv.FormatInt(resp[0].CreatedEpoch, 10)
		toId := strconv.FormatInt(resp[len(resp)-1].CreatedEpoch, 10)
		return fromId, toId, "/api/photos?limit=10&fromId=%s"
	}
	return "", "", ""
}

// ListPhotos godoc
//
//	@Summary		Get photos
//	@Description	Get photos
//	@Tags			Photos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"User ID"
//	@Success		200	{object}	userResponse	"User displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users/{id} [get]
//	@Security		BearerAuth
func (ph *PhotoHandler) ListPhotos(ctx *gin.Context) {
	albumId := ctx.Query("albumId")
	fromId := ctx.Query("fromId")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	includeThumbnail, err := strconv.ParseBool(ctx.DefaultQuery("thumbnail", "false"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "thumbnail query parameter must be either true or false")
	}

	var photoResp []*domain.Photo

	if albumId != "" {
		photoResp, err = ph.photoSvc.ListPhotosInAlbum(albumId, fromId, limit, includeThumbnail)
		handleError(ctx, err)
	} else {
		photoResp, err = ph.photoSvc.ListPhotos(fromId, limit, includeThumbnail)
		handleError(ctx, err)
	}

	fromId, toId, nextPage := photosPaginationParams(photoResp)
	handlePaginatedSuccess(ctx, photoResp, fromId, toId, len(photoResp), nextPage)
}

func (ph *PhotoHandler) GetPhoto(ctx *gin.Context) {
	photoId := ctx.Param("id")[1:] // Strip leading slash from catch-all parameter
	includeThumbnail, err := strconv.ParseBool(ctx.DefaultQuery("thumbnail", "false"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "thumbnail query parameter must be either true or false")
	}

	photoResp, err := ph.photoSvc.GetPhoto(photoId, includeThumbnail)
	handleError(ctx, err)

	handleSuccess(ctx, photoResp)
}

func (ph *PhotoHandler) GetPhotoCount(ctx *gin.Context) {
	albumId := ctx.Query("albumId")
	if albumId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("album ID")})
	}

	photoCount, err := ph.photoSvc.PhotoCount(albumId)
	handleError(ctx, err)
	handleSuccess(ctx, photoCount)
}

func (ph *PhotoHandler) GetPhotoBin(ctx *gin.Context) {
	photoId := ctx.Param("id")[1:] // Strip leading slash from catch-all parameter
	if photoId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("photo ID")})
		return
	}

	photoBinary, err := ph.photoSvc.PhotoBinary(photoId)
	handleError(ctx, err)

	// Unescape single quotes in filesystem path
	unescapedPath := strings.ReplaceAll(photoBinary, `\'`, `'`)
	ctx.File(unescapedPath)
}

func (ph *PhotoHandler) GetPhotoThumbnail(ctx *gin.Context) {
	photoId := ctx.Param("id")[1:] // Strip leading slash from catch-all parameter
	if photoId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("photo ID")})
	}

	photoBinary, err := ph.photoSvc.PhotoThumbnail(photoId)
	handleError(ctx, err)
	ctx.Data(http.StatusOK, "application/octet-stream", photoBinary)
}

func (ph *PhotoHandler) IndexPhotos(c *gin.Context) {
	isRunning, _ := ph.jobSvc.IsJobRunning("Photo_index")

	if isRunning {
		c.IndentedJSON(http.StatusConflict, "Photo Index already running")
	} else {
		c.Status(http.StatusAccepted)
		go func() {
			ph.photoSvc.PerformPhotoIndex()
		}()
	}
}
