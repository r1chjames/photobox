package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
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

// svc returns the photo service scoped to the request's workspace (issue
// #74). Fail-closed: a request without a resolved workspace gets an empty
// (scoped) workspace, which matches nothing, rather than an unscoped service.
func (ph *PhotoHandler) svc(ctx *gin.Context) port.PhotoService {
	if wc := GetWorkspaceContext(ctx); wc != nil {
		return ph.photoSvc.WithWorkspace(wc.WorkspaceID)
	}
	return ph.photoSvc.WithWorkspace("")
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
		photoResp, err = ph.svc(ctx).SearchPhotosWithFilters(filters, fromId, limit, includeThumbnail)
	} else if tagsQuery != "" {
		tags := strings.Split(tagsQuery, ",")
		photoResp, err = ph.svc(ctx).ListPhotosByTags(tags, fromId, limit, includeThumbnail)
	} else if lowQuality {
		photoResp, err = ph.svc(ctx).ListLowQualityPhotos(fromId, limit, includeThumbnail)
	} else if favorites {
		photoResp, err = ph.svc(ctx).ListFavoritePhotos(fromId, limit, includeThumbnail, startDate, endDate)
	} else if albumId != "" {
		photoResp, err = ph.svc(ctx).ListPhotosInAlbum(albumId, fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	} else {
		photoResp, err = ph.svc(ctx).ListPhotos(fromId, limit, includeThumbnail, startDate, endDate, mediaType)
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

	// Scope to the caller's workspace (issue #74). NOTE: this filters after
	// the repository query, so a limit applied pre-filter can yield fewer rows
	// than requested — the Phase 2 query-scoping sweep moves this into the
	// query and removes both the filtering and the caveat.
	wc := GetWorkspaceContext(ctx)
	if wc == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	photoResp = filterByWorkspace(photoResp, wc.WorkspaceID, func(p *domain.Photo) string { return p.WorkspaceID })

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

	photoResp, err := ph.svc(ctx).GetPhoto(photoId, includeThumbnail)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// Cross-tenant guard (issue #74): 404 if the photo is not in the
	// caller's workspace.
	if !resourceInWorkspace(ctx, photoResp.WorkspaceID) {
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

	photoCount, err := ph.svc(ctx).PhotoCount(albumId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, photoCount)
}

// uploadPhotoRequest is the body for POST /photo (base64-encoded file).
type uploadPhotoRequest struct {
	Name          string `json:"name" binding:"required"`
	AlbumName     string `json:"albumName" binding:"required"`
	BinaryContent string `json:"binaryContent" binding:"required"`
}

// UploadPhoto accepts a base64-encoded photo, writes it to the filesystem
// and indexes it (issue #98).
func (ph *PhotoHandler) UploadPhoto(ctx *gin.Context) {
	var req uploadPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	photo, err := ph.svc(ctx).UploadPhoto(domain.PhotoUpload{
		Name:          req.Name,
		AlbumName:     req.AlbumName,
		BinaryContent: req.BinaryContent,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, photo)
}

func (ph *PhotoHandler) GetPhotoBin(ctx *gin.Context) {
	photoId := ctx.Param("id")[1:] // Strip leading slash from catch-all parameter
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	photoBinary, err := ph.svc(ctx).PhotoBinary(photoId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Unescape single quotes in filesystem path
	unescapedPath := strings.ReplaceAll(photoBinary, `\'`, `'`)

	// Enable HTTP range requests (seeking) for all binary content. gin's
	// ctx.File uses http.ServeFile which honours Range headers, but we set
	// Accept-Ranges explicitly so video players/clients know seeking is
	// supported (issue #103).
	ctx.Header("Accept-Ranges", "bytes")

	// Set a correct Content-Type for videos so browsers stream/seek rather
	// than download. http.ServeFile would infer from extension, but MOV/WebM
	// sometimes fall back to application/octet-stream.
	if strings.HasSuffix(strings.ToLower(unescapedPath), ".mp4") {
		ctx.Header("Content-Type", "video/mp4")
	} else if strings.HasSuffix(strings.ToLower(unescapedPath), ".mov") {
		ctx.Header("Content-Type", "video/quicktime")
	} else if strings.HasSuffix(strings.ToLower(unescapedPath), ".webm") {
		ctx.Header("Content-Type", "video/webm")
	} else if strings.HasSuffix(strings.ToLower(unescapedPath), ".mkv") {
		ctx.Header("Content-Type", "video/x-matroska")
	}

	ctx.File(unescapedPath)
}

// GetPhotoLiveVideo serves the paired Live Photo video for a given photo.
func (ph *PhotoHandler) GetPhotoLiveVideo(ctx *gin.Context) {
	photoId := ctx.Param("id")
	if photoId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "photo ID is required"})
		return
	}

	liveVideoPath, err := ph.svc(ctx).PhotoLiveVideoPath(photoId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Header("Accept-Ranges", "bytes")
	ctx.Header("Content-Type", "video/quicktime")

	unescapedPath := strings.ReplaceAll(liveVideoPath, `\'`, `'`)
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

	// Cross-tenant guard (issue #74): the thumbnail is addressable by photo
	// ID, so verify the photo is in the caller's workspace before serving.
	owner, err := ph.svc(ctx).GetPhoto(photoId, false)
	if err != nil {
		handleError(ctx, err)
		return
	}
	if !resourceInWorkspace(ctx, owner.WorkspaceID) {
		return
	}

	setThumbnailHeaders := func() {
		ctx.Header("Content-Type", "image/webp")
		ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
		ctx.Header("ETag", fmt.Sprintf(`"%s:%s"`, photoId, size))
	}

	// Try unified thumbnail retrieval (handles valkey or filesystem)
	thumbnailBytes, err := ph.svc(ctx).PhotoThumbnailBytesForSize(photoId, size)
	if err == nil && len(thumbnailBytes) > 0 {
		setThumbnailHeaders()
		ctx.Data(http.StatusOK, "image/webp", thumbnailBytes)
		return
	}

	// Generate on-demand if missing (always retry retrieval even if path
	// is empty — valkey-stored thumbnails don't have a filesystem path)
	_, err = ph.svc(ctx).GenerateThumbnailForPhoto(photoId)
	if err == nil {
		// After generation, try again
		thumbnailBytes, err = ph.svc(ctx).PhotoThumbnailBytesForSize(photoId, size)
		if err == nil && len(thumbnailBytes) > 0 {
			setThumbnailHeaders()
			ctx.Data(http.StatusOK, "image/webp", thumbnailBytes)
			return
		}
	}

	// Fallback to DB bytes (only for medium size)
	if size == "m" {
		photoBinary, err := ph.svc(ctx).PhotoThumbnailBytes(photoId)
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

// thumbCapRegex is a strict UUIDv4 pattern for capability URLs. Rejecting
// anything else up front keeps the DB lookup and any logging clean.
var thumbCapRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// GetThumbnailByCap serves thumbnails by capability URL:
// GET /t/{thumb_cap}/{size}.webp (issue #74 D1).
//
// The route is unauthenticated: the cap itself is the credential (random
// 128-bit UUID, immutable, unenumerable). Unknown or malformed caps return
// 404 — identical to a missing thumbnail — so the route leaks nothing.
// Responses are immutable-cacheable so a CDN edge can absorb reads.
func (ph *PhotoHandler) GetThumbnailByCap(ctx *gin.Context) {
	cap := ctx.Param("cap")
	if !thumbCapRegex.MatchString(cap) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	size := strings.TrimSuffix(ctx.Param("size"), ".webp")
	if size != "s" && size != "m" && size != "l" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	photo, err := ph.svc(ctx).GetPhotoByThumbCap(cap)
	if err != nil {
		// Unknown cap = 404 (same as missing thumbnail; no existence oracle).
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Serve from thumbnail storage. If the thumbnail is missing (e.g. not yet
	// generated for a just-indexed photo), generate on demand once.
	thumbnailBytes, err := ph.svc(ctx).PhotoThumbnailBytesForSize(photo.ID, size)
	if err != nil || len(thumbnailBytes) == 0 {
		if _, genErr := ph.svc(ctx).GenerateThumbnailForPhoto(photo.ID); genErr != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		thumbnailBytes, err = ph.svc(ctx).PhotoThumbnailBytesForSize(photo.ID, size)
		if err != nil || len(thumbnailBytes) == 0 {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	ctx.Header("Content-Type", "image/webp")
	// Capability URLs never change per photo: cache forever at the edge.
	ctx.Header("Cache-Control", "public, max-age=31536000, immutable")
	ctx.Header("ETag", fmt.Sprintf(`"%s:%s"`, photo.ThumbCap, size))
	ctx.Data(http.StatusOK, "image/webp", thumbnailBytes)
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
			ph.svc(c).PerformPhotoIndex(ctx)
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
			ph.svc(c).RegenerateThumbnails(ctx)
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
			if err := ph.svc(c).AnalyzeExistingPhotos(); err != nil {
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
			scored, err := ph.svc(c).ScoreAllPhotoQuality(200)
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

	err := ph.svc(ctx).DeletePhoto(photoId)
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

	err := ph.svc(ctx).RestorePhoto(photoId)
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

	photoResp, err := ph.svc(ctx).ListTrashPhotos(fromId, limit, includeThumbnail)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// Scope the listing to the caller's workspace (issue #74). The query
	// itself is not yet workspace-filtered (Phase 2 sweep), so filter here to
	// avoid returning other tenants' trashed photos.
	wc := GetWorkspaceContext(ctx)
	if wc == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	photoResp = filterByWorkspace(photoResp, wc.WorkspaceID, func(ph *domain.Photo) string { return ph.WorkspaceID })

	fromId, toId, nextPage := photosPaginationParams(photoResp)
	handlePaginatedSuccess(ctx, photoResp, fromId, toId, len(photoResp), nextPage)
}

func (ph *PhotoHandler) EmptyTrash(ctx *gin.Context) {
	err := ph.svc(ctx).EmptyTrash()
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

	// Cross-tenant guard (issue #74): verify ownership before mutating. The
	// resource's workspace is immutable, so the check cannot be raced.
	existing, err := ph.svc(ctx).GetPhoto(photoId, false)
	if err != nil {
		handleError(ctx, err)
		return
	}
	if !resourceInWorkspace(ctx, existing.WorkspaceID) {
		return
	}

	photo, err := ph.svc(ctx).SetFavorite(photoId, req.Favorite)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

// updateMetadataRequest carries user-editable metadata overrides. All fields
// are optional; only provided fields are persisted.
type updateMetadataRequest struct {
	Description *string  `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	DateTaken   *string  `json:"dateTaken"`
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

	photo, err := ph.svc(ctx).UpdatePhotoMetadata(photoId, description, req.Latitude, req.Longitude, req.DateTaken)
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

	location, err := ph.svc(ctx).ReverseGeocode(ctx, photoId)
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
		if err := ph.svc(ctx).BatchAddToAlbum(req.PhotoIds, req.AlbumId); err != nil {
			handleError(ctx, err)
			return
		}
	case "delete":
		if err := ph.svc(ctx).BatchDeletePhotos(req.PhotoIds); err != nil {
			handleError(ctx, err)
			return
		}
	case "favorite":
		if req.Favorite == nil {
			validationError(ctx, fmt.Errorf("favorite is required for favorite action"))
			return
		}
		if err := ph.svc(ctx).BatchSetFavorite(req.PhotoIds, *req.Favorite); err != nil {
			handleError(ctx, err)
			return
		}
	case "tag":
		operation := req.TagOperation
		if operation == "" {
			operation = "add"
		}
		if err := ph.svc(ctx).BatchUpdatePhotoTags(req.PhotoIds, req.Tags, operation); err != nil {
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
	if err := ph.svc(ctx).BatchAddToAlbum(req.PhotoIds, albumId); err != nil {
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

	_ = ph.svc(ctx).DownloadPhotos(req.PhotoIds, ctx.Writer)
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

	photo, err := ph.svc(ctx).RotatePhoto(photoId, direction)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

// editPhotoRequest is the body for POST /photos/:id/edit. All fields are
// optional; the edit is applied to a copy, never the original.
type editPhotoRequest struct {
	Rotate      *int               `json:"rotate"`
	Crop        *domain.CropParams `json:"crop"`
	Brightness  *float64           `json:"brightness"`
	Contrast    *float64           `json:"contrast"`
	Saturation  *float64           `json:"saturation"`
	AutoEnhance *bool              `json:"autoEnhance"`
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

	photo, err := ph.svc(ctx).EditPhoto(photoId, params)
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

	photo, err := ph.svc(ctx).ClearEdits(photoId)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, photo)
}

func (ph *PhotoHandler) GetTimeline(ctx *gin.Context) {
	entries, err := ph.svc(ctx).GetTimeline()
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

	groups, err := ph.svc(ctx).ListMemories(int(month), day, 10)
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

	entries, err := ph.svc(ctx).GetGeodata(north, south, east, west)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, entries)
}

func (ph *PhotoHandler) GetDuplicatePhotos(ctx *gin.Context) {
	photos, err := ph.svc(ctx).GetDuplicatePhotos()
	if err != nil {
		handleError(ctx, err)
		return
	}
	// Scope to the caller's workspace (issue #74).
	wc := GetWorkspaceContext(ctx)
	if wc == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	photos = filterByWorkspace(photos, wc.WorkspaceID, func(p *domain.Photo) string { return p.WorkspaceID })
	handleSuccess(ctx, photos)
}

func (ph *PhotoHandler) GetAllTags(ctx *gin.Context) {
	tags, err := ph.svc(ctx).GetAllTags()
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

	photo, err := ph.svc(ctx).UpdatePhotoTags(photoId, req.Tags)
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

	err := ph.svc(ctx).BatchUpdatePhotoTags(req.PhotoIds, req.Tags, req.Operation)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Tags updated"})
}
