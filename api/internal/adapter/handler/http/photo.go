package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/components"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"strconv"
)
import "net/http"

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
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	includeThumbnail, err := strconv.ParseBool(ctx.DefaultQuery("thumbnail", "false"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "thumbnail query parameter must be either true or false")
	}

	var photoResp []*domain.Photo

	if albumId != "" {
		photoResp, err = ph.photoSvc.ListPhotosInAlbum(albumId, page, limit, includeThumbnail)
		handleError(ctx, err)
	} else {
		photoResp, err = ph.photoSvc.ListPhotos(page, limit, includeThumbnail)
		handleError(ctx, err)
	}

	handleSuccess(ctx, photoResp)
}

func (ph *PhotoHandler) GetPhoto(ctx *gin.Context) {
	photoId := ctx.Param("id")
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
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("photo ID")})
	}

	photoBinary, err := ph.photoSvc.PhotoBinary(photoId)
	handleError(ctx, err)
	ctx.File(photoBinary)
}

func (ph *PhotoHandler) GetPhotoThumbnail(ctx *gin.Context) {
	photoId := ctx.Param("id")
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
			components.PerformPhotoIndex()
		}()
	}
}
