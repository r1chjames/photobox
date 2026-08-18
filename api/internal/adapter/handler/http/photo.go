package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

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
		return fromId, toId, "/api/photos?limit=60&fromId=%s"
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
	lowQuality, _ := strconv.ParseBool(ctx.DefaultQuery("lowQuality", "false"))
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var photoResp []*domain.Photo

	mediaType := ctx.Query("mediaType")

	tagsQuery := ctx.Query("tags")

	// Advanced search filters (issue 101): camera, hasGps, orientation.
	// When any is present, route through the combined filter query so the
	// filters compose with each other and with date/mediaType.
	camera := ctx.Query("camera")
	hasGps, _ := strconv.ParseBool(ctx.DefaultQuery("hasGps", "false"))
	orientation := ctx.Query("orientation")
	advanced := camera != "" || hasGps || orientation != ""

	if advanced {
		filters := domain.PhotoSearchFilters{
			StartDate:   startDate,
			EndDate:     endDate,
			MediaType:   mediaType,
			Camera:      camera,
			HasGPS:      hasGps,
			Orientation: orientation,
			Tags:        nil,
			Favorite:    favorites,
			LowQuality:  lowQuality,
		}
		photoResp, err = ph.photoSvc.SearchPhotosWithFilters(filters, fromId, limit, includeThumbnail)
	} else if tagsQuery != "" {
		tags := strings.Split(tagsQuery, ",")
		photoResp, err = ph.photoSvc.ListPhotosByTags(tags, fromId, limit, includeThumbnail)
	} else if lowQuality {
		photoResp, err = ph.photoSvc.ListLowQualityPhotos(fromId, limit, includeThumbnail)
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
		ctx.Header("Content-Type", "image/webp")
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		ctx.Header("ETag", fmt.Sprintf(`"%s:%s"`, photoId, size))
	}

	// Try unified thumbnail retrieval (handles valkey or filesystem)
	thumbnailBytes, err := ph.photoSvc.PhotoThumbnailBytesForSize(photoId, size)
	if err == nil && len(thumbnailBytes) > 0 {
		setThumbnailHeaders()
		ctx.Data(http.StatusOK, "image/webp", thumbnailBytes)
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
			ctx.Data(http.StatusOK, "image/webp", thumbnailBytes)
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
		ctx.Header("Content-Type", "image/webp")
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		ctx.Header("ETag", fmt.Sprintf(`"%s:%s"`, photoId, size))
		ctx.Data(http.StatusOK, "image/webp", photoBinary)
		return
	}

	handleError(ctx, domain.ErrDataNotFound)
}

func (ph *PhotoHandler) StartJob(c *gin.Context) {
	jobType := c.Param("type")
	
	var friendlyName string
	var runFunc func()
	
	switch jobType {
	case "Photo_index":
		friendlyName = "Photo indexing"
		runFunc = func() {
			ctx, cancel := context.WithCancel(context.Background())
			ph.jobCancellersMu.Lock()
			ph.jobCancellers[jobType] = cancel
			ph.jobCancellersMu.Unlock()
			defer func() {
				ph.jobCancellersMu.Lock()
				delete(ph.jobCancellers, jobType)
				ph.jobCancellersMu.Unlock()
				cancel()
			}()
			ph.photoSvc.PerformPhotoIndex(ctx)
		}
	case "Thumbnail_regenerate":
		friendlyName = "Thumbnail regeneration"
		runFunc = func() {
			ctx, cancel := context.WithCancel(context.Background())
			ph.jobCancellersMu.Lock()
			ph.jobCancellers[jobType] = cancel
			ph.jobCancellersMu.Unlock()
			defer func() {
				ph.jobCancellersMu.Lock()
				delete(ph.jobCancellers, jobType)
				ph.jobCancellersMu.Unlock()
				cancel()
			}()
			ph.photoSvc.RegenerateThumbnails(ctx)
		}
	case "AI_analysis":
		friendlyName = "AI photo analysis"
		runFunc = func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("AI analysis job panicked", "recover", r)
				}
				if err := ph.jobSvc.JobComplete(jobType); err != nil {
					slog.Error("Unable to complete AI analysis job", "error", err)
				}
			}()
			if err := ph.photoSvc.AnalyzeExistingPhotos(); err != nil {
				slog.Error("AI analysis job failed", "error", err)
			}
		}
	case "Quality_analysis":
		friendlyName = "Quality analysis"
		runFunc = func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Quality analysis job panicked", "recover", r)
				}
				if err := ph.jobSvc.JobComplete(jobType); err != nil {
					slog.Error("Unable to complete quality analysis job", "error", err)
				}
			}()
			scored, err := ph.photoSvc.ScoreAllPhotoQuality(200)
			if err != nil {
				slog.Error("Quality analysis job failed", "error", err)
				return
			}
			slog.Info("Quality analysis complete", "scored", scored)
		}
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unknown job type: " + jobType})
		return
	}
	
	// Atomically start the job - returns error if already running
	err := ph.jobSvc.StartJobIfNotRunning(jobType)
	if err != nil {
		if errors.Is(err, domain.ErrJobAlreadyRunning) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": friendlyName + " is already in progress"})
			return
		}
		handleError(c, err)
		return
	}
	
	c.Status(http.StatusAccepted)
	go runFunc()
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

func (ph *PhotoHandler) GetJobStatuses(c *gin.Context) {
	jobs, err := ph.jobSvc.GetAllJobs()
	if err != nil {
		handleError(c, err)
		return
	}
	handleSuccess(c, gin.H{"jobs": jobs})
}

func (ph *PhotoHandler) StopAllJobs(c *gin.Context) {
	ph.jobCancellersMu.Lock()
	for jobType, cancel := range ph.jobCancellers {
		cancel()
		delete(ph.jobCancellers, jobType)
	}
	ph.jobCancellersMu.Unlock()

	err := ph.jobSvc.UpdateAllJobsStatus("NOT_RUNNING")
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All jobs stopped"})
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

// updateMetadataRequest carries user-editable metadata overrides. All fields
// are optional; only provided fields are persisted.
type updateMetadataRequest struct {
	Description *string   `json:"description"`
	Latitude    *float64  `json:"latitude"`
	Longitude   *float64  `json:"longitude"`
	DateTaken   *string   `json:"dateTaken"`
}

// UpdatePhotoMetadata persists user-editable metadata overrides for a photo.
func (ph *PhotoHandler) UpdatePhotoMetadata(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	var req updateMetadataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	photo, err := ph.photoSvc.UpdatePhotoMetadata(photoId, description, req.Latitude, req.Longitude, req.DateTaken)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

// PhotoLocationResponse is the reverse-geocoding result.
type PhotoLocationResponse struct {
	Location string `json:"location"`
}

// GetPhotoLocation reverse-geocodes a photo's GPS coordinates.
func (ph *PhotoHandler) GetPhotoLocation(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	location, err := ph.photoSvc.ReverseGeocode(ctx, photoId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, PhotoLocationResponse{Location: location})
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

// batchPhotosRequest is the unified bulk-operation body:
//
//	{ "photoIds": [...], "action": "add_to_album"|"delete"|"favorite"|"tag", "albumId": "...", "favorite": true, "tags": [...], "tagOperation": "add" }
type batchPhotosRequest struct {
	PhotoIds     []string `json:"photoIds" binding:"required,min=1"`
	Action       string   `json:"action" binding:"required,oneof=add_to_album delete favorite tag"`
	AlbumId      string   `json:"albumId"`
	Favorite     *bool    `json:"favorite"`
	Tags         []string `json:"tags"`
	TagOperation string   `json:"tagOperation" binding:"omitempty,oneof=add remove set"`
}

// BatchPhotos applies a single bulk action to many photos in one request.
func (ph *PhotoHandler) BatchPhotos(ctx *gin.Context) {
	var req batchPhotosRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	switch req.Action {
	case "add_to_album":
		if req.AlbumId == "" {
			validationError(ctx, fmt.Errorf("albumId is required for add_to_album action"))
			return
		}
		if err := ph.photoSvc.BatchAddToAlbum(req.PhotoIds, req.AlbumId); err != nil {
			handleError(ctx, err)
			return
		}
	case "delete":
		if err := ph.photoSvc.BatchDeletePhotos(req.PhotoIds); err != nil {
			handleError(ctx, err)
			return
		}
	case "favorite":
		if req.Favorite == nil {
			validationError(ctx, fmt.Errorf("favorite is required for favorite action"))
			return
		}
		if err := ph.photoSvc.BatchSetFavorite(req.PhotoIds, *req.Favorite); err != nil {
			handleError(ctx, err)
			return
		}
	case "tag":
		operation := req.TagOperation
		if operation == "" {
			operation = "add"
		}
		if err := ph.photoSvc.BatchUpdatePhotoTags(req.PhotoIds, req.Tags, operation); err != nil {
			handleError(ctx, err)
			return
		}
	default:
		validationError(ctx, fmt.Errorf("unsupported batch action: %s", req.Action))
		return
	}

	handleSuccess(ctx, gin.H{"message": fmt.Sprintf("Batch %s completed for %d photos", req.Action, len(req.PhotoIds))})
}

// addPhotosToAlbumRequest is the body for POST /albums/:id/photos.
type addPhotosToAlbumRequest struct {
	PhotoIds []string `json:"photoIds" binding:"required,min=1"`
}

// AddPhotosToAlbum assigns photos to the album named by the URL :id param.
func (ph *PhotoHandler) AddPhotosToAlbum(ctx *gin.Context) {
	albumId := ctx.Param("id")
	if albumId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "album ID is required"})
		return
	}
	var req addPhotosToAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := ph.photoSvc.BatchAddToAlbum(req.PhotoIds, albumId); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, gin.H{"message": fmt.Sprintf("Added %d photos to album", len(req.PhotoIds))})
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

// editPhotoRequest is the body for POST /photos/:id/edit. All fields are
// optional; the edit is applied to a copy, never the original.
type editPhotoRequest struct {
	Rotate      *int                    `json:"rotate"`
	Crop        *domain.CropParams      `json:"crop"`
	Brightness  *float64                `json:"brightness"`
	Contrast    *float64                `json:"contrast"`
	Saturation  *float64                `json:"saturation"`
	AutoEnhance *bool                   `json:"autoEnhance"`
}

// EditPhoto applies a non-destructive edit to a photo.
func (ph *PhotoHandler) EditPhoto(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	var req editPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	params := domain.EditParams{
		Rotate:      req.Rotate,
		Crop:        req.Crop,
		Brightness:  req.Brightness,
		Contrast:    req.Contrast,
		Saturation:  req.Saturation,
		AutoEnhance: req.AutoEnhance,
	}

	photo, err := ph.photoSvc.EditPhoto(photoId, params)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

// ClearEdits discards any edited copy and stored edit params.
func (ph *PhotoHandler) ClearEdits(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	photo, err := ph.photoSvc.ClearEdits(photoId)
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

// GetMemories returns "On This Day" photos grouped by year. Accepts an
// optional date=YYYY-MM-DD (defaults to today).
func (ph *PhotoHandler) GetMemories(ctx *gin.Context) {
	month, day := time.Now().Month(), time.Now().Day()
	if dateStr := ctx.Query("date"); dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			validationError(ctx, err)
			return
		}
		month, day = t.Month(), t.Day()
	}

	groups, err := ph.photoSvc.ListMemories(int(month), day, 10)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, groups)
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
