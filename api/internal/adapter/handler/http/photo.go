package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

const maxPageLimit = 100

// PhotoHandler represents the HTTP handler for photo-related requests
type PhotoHandler struct {
	photoSvc        port.PhotoService
	jobSvc          port.JobService
	jobCancellers   map[string]context.CancelFunc
	jobCancellersMu sync.Mutex
}

// NewPhotoHandler creates a new PhotoHandler instance
func NewPhotoHandler(photoSvc port.PhotoService, jobSvc port.JobService) *PhotoHandler {
	return &PhotoHandler{
		photoSvc:      photoSvc,
		jobSvc:        jobSvc,
		jobCancellers: make(map[string]context.CancelFunc),
	}
}

func photosPaginationParams(resp []*domain.Photo) (string, string, string) {
	if len(resp) > 0 {
		fromId := strconv.FormatInt(resp[0].CreatedEpoch, 10)
		toId := strconv.FormatInt(resp[len(resp)-1].CreatedEpoch, 10)
		return fromId, toId, "/api/photos?limit=10&fromId=%s"
	}
	return "", "", ""
}

func (ph *PhotoHandler) ListPhotos(ctx *gin.Context) {
	albumId := ctx.Query("albumId")
	fromId := ctx.Query("fromId")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	} else if limit > maxPageLimit {
		limit = maxPageLimit
	}
	includeThumbnail, err := strconv.ParseBool(ctx.DefaultQuery("thumbnail", "false"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "thumbnail query parameter must be either true or false"})
		return
	}

	favorites, _ := strconv.ParseBool(ctx.DefaultQuery("favorites", "false"))
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var photoResp []*domain.Photo

	mediaType := ctx.Query("mediaType")

	tagsQuery := ctx.Query("tags")
	if tagsQuery != "" {
		tags := strings.Split(tagsQuery, ",")
		photoResp, err = ph.photoSvc.ListPhotosByTags(tags, fromId, limit, includeThumbnail)
	} else if favorites {
		photoResp, err = ph.photoSvc.ListFavoritePhotos(fromId, limit, includeThumbnail, startDate, endDate)
	} else if albumId != "" {
		photoResp, err = ph.photoSvc.ListPhotosInAlbum(albumId, fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	} else {
		photoResp, err = ph.photoSvc.ListPhotos(fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	}

	if err != nil {
		handleError(ctx, err)
		return
	}

	// Post-DB mediaType filter is only needed for tags and favorites routes
	// since ListPhotos and ListPhotosInAlbum now filter at the DB level
	if mediaType != "" && (tagsQuery != "" || favorites) {
		filtered := make([]*domain.Photo, 0, len(photoResp))
		for _, p := range photoResp {
			if p.MediaType == mediaType {
				filtered = append(filtered, p)
			}
		}
		photoResp = filtered
	}

	fromId, toId, nextPage := photosPaginationParams(photoResp)
	handlePaginatedSuccess(ctx, photoResp, fromId, toId, len(photoResp), nextPage)
}

func (ph *PhotoHandler) GetPhoto(ctx *gin.Context) {
	photoId := ctx.Param("id")[1:] // Strip leading slash from catch-all parameter
	includeThumbnail, err := strconv.ParseBool(ctx.DefaultQuery("thumbnail", "false"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "thumbnail query parameter must be either true or false"})
		return
	}

	photoResp, err := ph.photoSvc.GetPhoto(photoId, includeThumbnail)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photoResp)
}

func (ph *PhotoHandler) GetPhotoCount(ctx *gin.Context) {
	albumId := ctx.Query("albumId")
	if albumId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "albumId query parameter is required"})
		return
	}

	photoCount, err := ph.photoSvc.PhotoCount(albumId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, photoCount)
}

func (ph *PhotoHandler) GetPhotoBin(ctx *gin.Context) {
	photoId := ctx.Param("id")[1:] // Strip leading slash from catch-all parameter
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	photoBinary, err := ph.photoSvc.PhotoBinary(photoId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Unescape single quotes in filesystem path
	unescapedPath := strings.ReplaceAll(photoBinary, `\'`, `'`) 
	ctx.File(unescapedPath)
}

func (ph *PhotoHandler) GetPhotoThumbnail(ctx *gin.Context) {
	photoId := ctx.Param("id")[1:] // Strip leading slash from catch-all parameter
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	size := ctx.DefaultQuery("size", "m")
	if size != "s" && size != "m" && size != "l" {
		size = "m"
	}

	setThumbnailHeaders := func() {
		ctx.Header("Content-Type", "image/jpeg")
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		ctx.Header("ETag", fmt.Sprintf(`"%s:%s"`, photoId, size))
	}

	// Try unified thumbnail retrieval (handles valkey or filesystem)
	thumbnailBytes, err := ph.photoSvc.PhotoThumbnailBytesForSize(photoId, size)
	if err == nil && len(thumbnailBytes) > 0 {
		setThumbnailHeaders()
		ctx.Data(http.StatusOK, "image/jpeg", thumbnailBytes)
		return
	}

	// Generate on-demand if missing (always retry retrieval even if path
	// is empty — valkey-stored thumbnails don't have a filesystem path)
	_, err = ph.photoSvc.GenerateThumbnailForPhoto(photoId)
	if err == nil {
		// After generation, try again
		thumbnailBytes, err = ph.photoSvc.PhotoThumbnailBytesForSize(photoId, size)
		if err == nil && len(thumbnailBytes) > 0 {
			setThumbnailHeaders()
			ctx.Data(http.StatusOK, "image/jpeg", thumbnailBytes)
			return
		}
	}

	// Fallback to DB bytes (only for medium size)
	if size == "m" {
		photoBinary, err := ph.photoSvc.PhotoThumbnailBytes(photoId)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ctx.Header("Content-Type", "image/jpeg")
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		ctx.Header("ETag", fmt.Sprintf(`"%s:%s"`, photoId, size))
		ctx.Data(http.StatusOK, "image/jpeg", photoBinary)
		return
	}

	handleError(ctx, domain.ErrDataNotFound)
}

func (ph *PhotoHandler) IndexPhotos(c *gin.Context) {
	// Atomically start the job - returns error if already running
	err := ph.jobSvc.StartJobIfNotRunning("Photo_index")
	if err != nil {
		if errors.Is(err, domain.ErrJobAlreadyRunning) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Photo indexing is already in progress"})
			return
		}
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	ph.jobCancellersMu.Lock()
	ph.jobCancellers["Photo_index"] = cancel
	ph.jobCancellersMu.Unlock()

	c.Status(http.StatusAccepted)
	go func() {
		defer func() {
			ph.jobCancellersMu.Lock()
			delete(ph.jobCancellers, "Photo_index")
			ph.jobCancellersMu.Unlock()
			cancel()
		}()
		ph.photoSvc.PerformPhotoIndex(ctx)
	}()
}

func (ph *PhotoHandler) RegenerateThumbnails(c *gin.Context) {
	err := ph.jobSvc.StartJobIfNotRunning("Thumbnail_regenerate")
	if err != nil {
		if errors.Is(err, domain.ErrJobAlreadyRunning) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Thumbnail regeneration is already in progress"})
			return
		}
		handleError(c, err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	ph.jobCancellersMu.Lock()
	ph.jobCancellers["Thumbnail_regenerate"] = cancel
	ph.jobCancellersMu.Unlock()

	c.Status(http.StatusAccepted)
	go func() {
		defer func() {
			ph.jobCancellersMu.Lock()
			delete(ph.jobCancellers, "Thumbnail_regenerate")
			ph.jobCancellersMu.Unlock()
			cancel()
		}()
		ph.photoSvc.RegenerateThumbnails(ctx)
	}()
}

func (ph *PhotoHandler) StopJob(c *gin.Context) {
	jobType := c.Param("type")

	ph.jobCancellersMu.Lock()
	cancel, exists := ph.jobCancellers[jobType]
	ph.jobCancellersMu.Unlock()

	if !exists || cancel == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No running job of type: " + jobType})
		return
	}

	cancel()

	// Update job status in DB
	_ = ph.jobSvc.JobComplete(jobType)

	c.JSON(http.StatusOK, gin.H{"message": "Job stop requested"})
}

func (ph *PhotoHandler) DeletePhoto(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	err := ph.photoSvc.DeletePhoto(photoId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Photo moved to trash"})
}

func (ph *PhotoHandler) RestorePhoto(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	err := ph.photoSvc.RestorePhoto(photoId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Photo restored"})
}

func (ph *PhotoHandler) ListTrashPhotos(ctx *gin.Context) {
	fromId := ctx.Query("fromId")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	} else if limit > maxPageLimit {
		limit = maxPageLimit
	}
	includeThumbnail, _ := strconv.ParseBool(ctx.DefaultQuery("thumbnail", "false"))

	photoResp, err := ph.photoSvc.ListTrashPhotos(fromId, limit, includeThumbnail)
	if err != nil {
		handleError(ctx, err)
		return
	}

	fromId, toId, nextPage := photosPaginationParams(photoResp)
	handlePaginatedSuccess(ctx, photoResp, fromId, toId, len(photoResp), nextPage)
}

func (ph *PhotoHandler) EmptyTrash(ctx *gin.Context) {
	err := ph.photoSvc.EmptyTrash()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, gin.H{"message": "Trash emptied"})
}

type setFavoriteRequest struct {
	Favorite bool `json:"favorite"`
}

func (ph *PhotoHandler) SetFavorite(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	var req setFavoriteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	photo, err := ph.photoSvc.SetFavorite(photoId, req.Favorite)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

type downloadPhotosRequest struct {
	PhotoIds []string `json:"photoIds" binding:"required,min=1"`
}

type updateTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

type batchTagsRequest struct {
	PhotoIds  []string `json:"photoIds" binding:"required,min=1"`
	Tags      []string `json:"tags" binding:"required,min=1"`
	Operation string   `json:"operation" binding:"required,oneof=add remove set"`
}

func (ph *PhotoHandler) DownloadPhotos(ctx *gin.Context) {
	var req downloadPhotosRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", `attachment; filename="photobox-download.zip"`)

	_ = ph.photoSvc.DownloadPhotos(req.PhotoIds, ctx.Writer)
}

func (ph *PhotoHandler) RotatePhoto(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	direction := ctx.DefaultQuery("direction", "cw")
	if direction != "cw" && direction != "ccw" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "direction must be cw or ccw"})
		return
	}

	photo, err := ph.photoSvc.RotatePhoto(photoId, direction)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

func (ph *PhotoHandler) GetTimeline(ctx *gin.Context) {
	entries, err := ph.photoSvc.GetTimeline()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, entries)
}

func (ph *PhotoHandler) GetGeodata(ctx *gin.Context) {
	north, _ := strconv.ParseFloat(ctx.DefaultQuery("north", "90"), 64)
	south, _ := strconv.ParseFloat(ctx.DefaultQuery("south", "-90"), 64)
	east, _ := strconv.ParseFloat(ctx.DefaultQuery("east", "180"), 64)
	west, _ := strconv.ParseFloat(ctx.DefaultQuery("west", "-180"), 64)

	entries, err := ph.photoSvc.GetGeodata(north, south, east, west)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, entries)
}

func (ph *PhotoHandler) GetDuplicatePhotos(ctx *gin.Context) {
	photos, err := ph.photoSvc.GetDuplicatePhotos()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, photos)
}

func (ph *PhotoHandler) GetAllTags(ctx *gin.Context) {
	tags, err := ph.photoSvc.GetAllTags()
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, tags)
}

func (ph *PhotoHandler) UpdatePhotoTags(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	var req updateTagsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	photo, err := ph.photoSvc.UpdatePhotoTags(photoId, req.Tags)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

func (ph *PhotoHandler) BatchUpdatePhotoTags(ctx *gin.Context) {
	var req batchTagsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := ph.photoSvc.BatchUpdatePhotoTags(req.PhotoIds, req.Tags, req.Operation)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Tags updated"})
}
