package service

import (
	archivezip "archive/zip"
	b64 "encoding/base64"
	"bytes"
	"context"
	"fmt"
	json "github.com/goccy/go-json"
	"hash/fnv"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buckket/go-blurhash"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	ws "gitlab.com/r1chjames/photobox/api/internal/components/websocket"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

// getPhotoEpoch extracts the photo's date from EXIF DateTimeOriginal and returns
// it as epoch milliseconds. Falls back to the Unix timestamp (seconds) if EXIF
// date is unavailable, and to time.Now() as a last resort.
func getPhotoEpoch(exifData exif.Exif, unixFallback int64) int64 {
	tag, err := exifData.Get(exif.DateTimeOriginal)
	if err == nil {
		dateStr, err := tag.StringVal()
		if err == nil {
			// EXIF DateTimeOriginal format: "2024:05:15 14:30:00"
			t, err := time.Parse("2006:01:02 15:04:05", dateStr)
			if err == nil {
				return t.UnixMilli()
			}
		}
	}
	if unixFallback > 0 {
		return unixFallback * 1000
	}
	return time.Now().UnixMilli()
}

// getPhotoYearMonth extracts the year and month from EXIF DateTimeOriginal,
// falling back to the Unix timestamp (seconds). Returns (0, 0) if no valid
// date can be determined.
func getPhotoYearMonth(exifData exif.Exif, unixFallback int64) (int, int) {
	tag, err := exifData.Get(exif.DateTimeOriginal)
	if err == nil {
		dateStr, err := tag.StringVal()
		if err == nil {
			t, err := time.Parse("2006:01:02 15:04:05", dateStr)
			if err == nil {
				return t.Year(), int(t.Month())
			}
		}
	}
	if unixFallback > 0 {
		t := time.Unix(unixFallback, 0)
		return t.Year(), int(t.Month())
	}
	return 0, 0
}

type PhotoService struct {
	photoRepo        port.PhotoRepository
	albumSvc         port.AlbumService
	filesystemSvc    port.FilesystemService
	cacheSvc         port.CacheService
	aiSvc            port.AIService
	jobSvc           port.JobService
	config           appconfig.AppConfig
	thumbnailStorage port.ThumbnailStorage
	wsHub            *ws.Hub
	geocoder         *Geocoder
}

// NewPhotoService creates a new Photo service instance
func NewPhotoService(photoRepo port.PhotoRepository, albumRepo port.AlbumService, filesystemSvc port.FilesystemService, cacheSvc port.CacheService, aiSvc port.AIService, config appconfig.AppConfig, thumbnailStorage port.ThumbnailStorage, wsHub *ws.Hub, jobSvc port.JobService) *PhotoService {
	return &PhotoService{
		photoRepo,
		albumRepo,
		filesystemSvc,
		cacheSvc,
		aiSvc,
		jobSvc,
		config,
		thumbnailStorage,
		wsHub,
		nil,
	}
}

// SetGeocoder installs the reverse-geocoder used by the metadata location
// enrichment. Optional: when nil, location lookups return empty.
func (ps *PhotoService) SetGeocoder(g *Geocoder) {
	ps.geocoder = g
}

// ReverseGeocode resolves a human-readable location for the photo's GPS
// coordinates. Returns empty when the photo has no coordinates or the
// geocoder is not configured.
func (ps *PhotoService) ReverseGeocode(ctx context.Context, photoId string) (string, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return "", err
	}
	if ps.geocoder == nil || (photo.Latitude == 0 && photo.Longitude == 0) {
		return "", nil
	}
	return ps.geocoder.Reverse(ctx, photo.Latitude, photo.Longitude)
}

func (ps *PhotoService) ListPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListAllPhotosInAlbum(albumId, fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}

func (ps *PhotoService) ListPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListAllPhotos(fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}

func (ps *PhotoService) GetPhoto(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	resp, err := ps.photoRepo.GetPhotoById(photoId, includeThumbnail)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	ps.setPhotoSourcePath(resp)
	return resp, nil
}

func (ps *PhotoService) PerformPhotoIndex(ctx context.Context) {
	// Notify clients that indexing has started
	if ps.wsHub != nil {
		ps.wsHub.BroadcastEvent(ws.Event{
			Type: ws.EventIndexProgress,
			Payload: ws.IndexProgressPayload{
				Phase: "started",
			},
		})
	}

	cache, _ := ps.photoRepo.GetPhotoIndexCache()
	ps.filesystemSvc.PerformPhotoIndex(ctx, ps.SavePhotos, cache)

	// Notify clients that indexing has completed
	if ps.wsHub != nil {
		ps.wsHub.BroadcastEvent(ws.Event{
			Type: ws.EventIndexComplete,
			Payload: ws.IndexCompletePayload{},
		})
	}
}

func (ps *PhotoService) PhotoCount(albumId string) (int64, error) {
	return ps.photoRepo.GetPhotosInAlbumCount(albumId)
}

func (ps *PhotoService) PhotoBinary(photoId string) (string, error) {
	photoInfo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return "", err
	}
	// Validate path is within photo directory
	absBase, _ := filepath.Abs(ps.config.PhotoDir)
	unescapedPath := utils.UnescapeInvalidCharacters(photoInfo.FilesystemPath)
	absReq, _ := filepath.Abs(unescapedPath)
	if !strings.HasPrefix(absReq, absBase) {
		return "", domain.ErrForbidden
	}
	return unescapedPath, nil
}

func (ps *PhotoService) PhotoThumbnail(photoId string) ([]byte, error) {
	photoInfo, err := ps.photoRepo.GetPhotoById(photoId, true)
	if err != nil {
		return nil, err
	}
	return photoInfo.Thumbnail, nil
}

func (ps *PhotoService) PhotoThumbnailBytes(photoId string) ([]byte, error) {
	return ps.photoRepo.GetThumbnailBytes(photoId)
}

func (ps *PhotoService) PhotoThumbnailPath(photoId string) (string, error) {
	return ps.photoRepo.GetThumbnailPath(photoId)
}

func (ps *PhotoService) PhotoThumbnails(photoIds []string) (map[string][]byte, error) {
	return ps.photoRepo.GetPhotoThumbnails(photoIds)
}

func (ps *PhotoService) setPhotoSourcePath(photo *domain.Photo) {
	photo.SourcePath = fmt.Sprintf("photo/%s/bin", photo.ID)
	photo.ThumbnailUrl = fmt.Sprintf("photo/%s/thumbnail", photo.ID)
}

func (ps *PhotoService) setPhotosSourcePath(photos []*domain.Photo) {
	for _, photo := range photos {
		ps.setPhotoSourcePath(photo)
	}
}

func computeFileHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	hash := fnv.New64a()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum64())
}

// computeDominantColor resizes the image to 1x1 to extract the average color,
// returning a hex string like "#aabbcc". Returns empty string on failure.
func computeDominantColor(path string) string {
	img, err := imaging.Open(path)
	if err != nil {
		return ""
	}
	// Resize to 1x1 to get the average color of the entire image
	onePixel := imaging.Resize(img, 1, 1, imaging.Box)
	r, g, b, _ := onePixel.At(0, 0).RGBA()
	// RGBA() returns values in [0, 65535]; shift down to 8-bit
	return fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

func (ps *PhotoService) SavePhotos(photos []domain.PhotoFile) error {
	if len(photos) == 0 {
		return nil
	}

	// Cache album lookups to avoid repeated queries
	albumCache := make(map[string]string) // map[albumName]albumId

	// Build all photo records
	photoRecords := make([]domain.Photo, 0, len(photos))

	for _, photo := range photos {
		// Check cache first
		albumId, found := albumCache[photo.Directory]

		if !found {
			// Look up or create album
			result, err := ps.albumSvc.GetAlbumByName(photo.Directory)
			if result == nil || err != nil {
				newAlbum, albErr := ps.albumSvc.CreateAlbum(photo.Directory)
				if albErr != nil {
					slog.Error("Unable to insert album record", "error", albErr)
					albumId = "" // Use empty album ID if creation fails
				} else {
					albumId = newAlbum.ID
					slog.Info("Created album", "name", photo.Directory, "id", albumId)
				}
			} else {
				albumId = result.ID
			}
			// Cache the result
			albumCache[photo.Directory] = albumId
		}

		photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
		photoMetadata, _ := json.Marshal(&photo)
		photoInfo := domain.Photo{
			ID:             photoHash,
			Name:           utils.EscapeInvalidCharacters(photo.Name),
			FilesystemPath: photo.Path,
			AlbumId:        albumId,
			Metadata:       photoMetadata,
			Thumbnail:      photo.Thumbnail,
			CreatedEpoch:   getPhotoEpoch(photo.Exif, photo.ModifiedTime),
			Year:           0,
			Month:          0,
			FileHash:       computeFileHash(photo.Path),
			FileModifiedTime: photo.ModifiedTime,
			MediaType:      photo.MediaType,
			Duration:       photo.Duration,
			Width:          photo.Width,
			Height:         photo.Height,
			Latitude:       photo.Latitude,
			Longitude:      photo.Longitude,
			DominantColor:  computeDominantColor(photo.Path),
		}
		photoInfo.Year, photoInfo.Month = getPhotoYearMonth(photo.Exif, photo.ModifiedTime)

		// Compute deterministic quality metrics (blur/exposure) for images.
		if photo.MediaType != "video" {
			if quality, qErr := AnalyzeQuality(photo.Path); qErr == nil && quality != nil {
				photoInfo.QualityScore = quality.Overall
				photoInfo.BlurScore = quality.BlurScore
				photoInfo.IsLowQuality = quality.IsLowQuality
			}
		}

		slog.Info("Adding photo", "photo", photo.Name, "album", photo.Directory)

		// Generate and store medium thumbnail at index time so thumbnails
		// are available immediately for both the batch endpoint (DB column)
		// and the individual endpoint (thumbnail storage).
		thumbnailBytes := ps.filesystemSvc.GenerateThumbnail(photo.Path, photo.Exif, 600, 600)
		if len(thumbnailBytes) > 0 {
			photoInfo.Thumbnail = thumbnailBytes
			if err := ps.thumbnailStorage.Put(context.Background(), photoHash, "m", thumbnailBytes); err != nil {
				slog.Error("Failed to store thumbnail", "photo", photo.Name, "error", err)
			}
			// Encode blurhash placeholder from thumbnail for instant grid rendering
			if img, _, err := image.Decode(bytes.NewReader(thumbnailBytes)); err == nil {
				if blurhashStr, err := blurhash.Encode(4, 3, img); err == nil {
					photoInfo.Blurhash = blurhashStr
				} else {
					slog.Warn("Failed to encode blurhash", "photo", photo.Name, "error", err)
				}
			} else {
				slog.Warn("Failed to decode thumbnail for blurhash", "photo", photo.Name, "error", err)
			}
			// Notify connected clients that a thumbnail is ready
			if ps.wsHub != nil {
				ps.wsHub.BroadcastEvent(ws.Event{
					Type: ws.EventThumbnailReady,
					Payload: ws.ThumbnailReadyPayload{
						PhotoID: photoHash,
						Size:    "m",
					},
				})
			}
		}

		photoRecords = append(photoRecords, photoInfo)
	}

	// Batch insert all photos
	err := ps.photoRepo.CreatePhotosInfo(photoRecords)
	if err != nil {
		return err
	}

	// Run AI analysis asynchronously for each photo if enabled
	if ps.config.AIEnabled && ps.aiSvc != nil {
		for _, photo := range photoRecords {
			go ps.analyzeAndTagPhoto(photo.ID, utils.UnescapeInvalidCharacters(photo.FilesystemPath))
		}
	}

	return nil
}

func (ps *PhotoService) SavePhoto(photo domain.PhotoFile) error {
	result, err := ps.albumSvc.GetAlbumByName(photo.Directory)
	var albumId string
	if result == nil || err != nil {
		newAlbumId, albErr := ps.albumSvc.CreateAlbum(photo.Directory)
		if albErr != nil {
			slog.Error("Unable to insert album record", "error", albErr)
			albumId = "" // Use empty album ID if creation fails
		} else {
			albumId = newAlbumId.ID
			slog.Info("Created album", "name", photo.Directory, "id", albumId)
		}
	} else {
		albumId = result.ID
	}
	photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
	photoMetadata, _ := json.Marshal(&photo)
	photoInfo := domain.Photo{
		ID:             photoHash,
		Name:           utils.EscapeInvalidCharacters(photo.Name),
		FilesystemPath: photo.Path,
		AlbumId:        albumId,
		Metadata:       photoMetadata,
		Thumbnail:      photo.Thumbnail,
		CreatedEpoch:   getPhotoEpoch(photo.Exif, photo.ModifiedTime),
		Year:           0,
		Month:          0,
		FileHash:       computeFileHash(photo.Path),
		FileModifiedTime: photo.ModifiedTime,
		MediaType:      photo.MediaType,
		Duration:       photo.Duration,
		Width:          photo.Width,
		Height:         photo.Height,
		Latitude:       photo.Latitude,
		Longitude:      photo.Longitude,
		DominantColor:  computeDominantColor(photo.Path),
	}
	photoInfo.Year, photoInfo.Month = getPhotoYearMonth(photo.Exif, photo.ModifiedTime)

	slog.Info("Adding photo", "photo", photo.Name, "album", photo.Directory)

	// Generate and store medium thumbnail at index time so thumbnails
	// are available immediately for both the batch endpoint (DB column)
	// and the individual endpoint (thumbnail storage).
	thumbnailBytes := ps.filesystemSvc.GenerateThumbnail(photo.Path, photo.Exif, 600, 600)
	if len(thumbnailBytes) > 0 {
		photoInfo.Thumbnail = thumbnailBytes
		if err := ps.thumbnailStorage.Put(context.Background(), photoHash, "m", thumbnailBytes); err != nil {
			slog.Error("Failed to store thumbnail", "photo", photo.Name, "error", err)
		}
		// Notify connected clients that a thumbnail is ready
		if ps.wsHub != nil {
			ps.wsHub.BroadcastEvent(ws.Event{
				Type: ws.EventThumbnailReady,
				Payload: ws.ThumbnailReadyPayload{
					PhotoID: photoHash,
					Size:    "m",
				},
			})
		}
	}

	err = ps.photoRepo.CreatePhotoInfo(photoInfo)
	if err != nil {
		return err
	}

	if ps.config.AIEnabled && ps.aiSvc != nil {
		go ps.analyzeAndTagPhoto(photoInfo.ID, utils.UnescapeInvalidCharacters(photoInfo.FilesystemPath))
	}

	return nil
}

func (ps *PhotoService) analyzeAndTagPhoto(photoId string, imagePath string) {
	startedAt := time.Now()
	analysis, err := ps.aiSvc.AnalyzeImage(imagePath)
	now := time.Now() // after LLM call — reflects actual completion time

	if err != nil {
		slog.Warn("AI analysis failed", "photo", photoId, "error", err)
		ps.recordAnalysisFailure(photoId, &startedAt, now, err)
		return
	}

	// Combine tags and objects for storage
	allTags := make([]string, 0, len(analysis.Tags)+len(analysis.Objects))
	allTags = append(allTags, analysis.Tags...)
	allTags = append(allTags, analysis.Objects...)

	if len(allTags) == 0 && analysis.Caption == "" {
		return
	}

	// Store tags in photo_tags for search
	if err := ps.photoRepo.AddAITags(photoId, allTags); err != nil {
		slog.Warn("Failed to store AI tags", "photo", photoId, "error", err)
	}

	// Store full analysis in photo_analysis
	tagsJSON, _ := json.Marshal(analysis.Tags)
	objectsJSON, _ := json.Marshal(analysis.Objects)
	photoAnalysis := domain.PhotoAnalysis{
		PhotoID:       photoId,
		Model:         ps.config.OllamaModel,
		Caption:       analysis.Caption,
		Tags:          tagsJSON,
		Objects:       objectsJSON,
		IsNSFW:        analysis.IsNSFW,
		IsPortrait:    analysis.IsPortrait,
		Status:        "completed",
		Attempts:      1,
		StartedAt:     &startedAt,
		LastAttemptAt: &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := ps.photoRepo.SavePhotoAnalysis(photoAnalysis); err != nil {
		slog.Warn("Failed to store photo analysis", "photo", photoId, "error", err)
	}

	slog.Info("AI analysis complete", "photo", photoId, "tags", len(allTags), "caption", analysis.Caption)
}

func (ps *PhotoService) recordAnalysisFailure(photoId string, startedAt *time.Time, attemptTime time.Time, err error) {
	existing, _ := ps.photoRepo.GetPhotoAnalysis(photoId)
	attempts := 1
	if existing != nil {
		attempts = existing.Attempts + 1
	}
	errMsg := err.Error()

	// Permanent failures — not worth retrying
	retryable := true
	if strings.Contains(errMsg, "unknown format") {
		retryable = false
	}

	analysis := domain.PhotoAnalysis{
		PhotoID:       photoId,
		Status:        "failed",
		Attempts:      attempts,
		ErrorMessage:  errMsg,
		Retryable:     retryable,
		StartedAt:     startedAt,
		LastAttemptAt: &attemptTime,
		UpdatedAt:     attemptTime,
	}
	if existing == nil {
		analysis.CreatedAt = attemptTime
	}
	if err := ps.photoRepo.SavePhotoAnalysis(analysis); err != nil {
		slog.Warn("Failed to record analysis failure", "photo", photoId, "error", err)
	}
}

func (ps *PhotoService) AnalyzeExistingPhotos() error {
	if !ps.config.AIEnabled {
		slog.Info("AI analysis skipped: AI_ENABLED is false")
		return nil
	}
	if ps.aiSvc == nil {
		slog.Warn("AI analysis skipped: aiSvc is nil")
		return nil
	}

	photos, err := ps.photoRepo.ListPhotosPendingAnalysis(50)
	if err != nil {
		slog.Error("Failed to list photos pending analysis", "error", err)
		return err
	}

	if len(photos) == 0 {
		slog.Info("AI analysis: no photos pending analysis")
		return nil
	}

	slog.Info("AI analysis job started", "batch", len(photos))

	// Process photos concurrently — matches OLLAMA_NUM_PARALLEL
	const workers = 2
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup
	var processed int32
	total := int32(len(photos))

	for _, photo := range photos {
		wg.Add(1)
		sem <- struct{}{} // acquire
		go func(p *domain.Photo) {
			defer func() {
				<-sem // release
				n := atomic.AddInt32(&processed, 1)
				if n%10 == 0 || n == total {
					slog.Info("AI analysis progress", "done", n, "total", total)
				}
				wg.Done()
			}()
			ps.analyzeAndTagPhoto(p.ID, utils.UnescapeInvalidCharacters(p.FilesystemPath))
		}(photo)
	}

	wg.Wait()
	slog.Info("AI analysis job finished", "total", total)
	return nil
}

func (ps *PhotoService) TriggerAIAnalysis() error {
	if !ps.config.AIEnabled || ps.aiSvc == nil {
		return fmt.Errorf("AI is not enabled")
	}
	go func() {
		if err := ps.AnalyzeExistingPhotos(); err != nil {
			slog.Error("Manual AI analysis trigger failed", "error", err)
		}
	}()
	return nil
}

func (ps *PhotoService) RegenerateThumbnails(ctx context.Context) {
	slog.Info("Starting thumbnail regeneration for all photos")

	total := 0
	skipped := 0
	fromId := ""

	for {
		select {
		case <-ctx.Done():
			slog.Info("Thumbnail regeneration cancelled", "generated", total, "skipped", skipped)
			return
		default:
		}

		photos, err := ps.photoRepo.ListAllPhotos(fromId, 100, false, "", "", "")
		if err != nil {
			slog.Error("Failed to list photos for thumbnail regeneration", "error", err)
			return
		}
		if len(photos) == 0 {
			break
		}

		for _, photo := range photos {
			if photo.FilesystemPath == "" {
				continue
			}

			// Skip if all 3 sizes already exist AND blurhash is present
			allExist := true
			for _, size := range []string{"s", "m", "l"} {
				exists, err := ps.thumbnailStorage.Exists(context.Background(), photo.ID, size)
				if err != nil || !exists {
					allExist = false
					break
				}
			}
			if allExist && photo.Blurhash != "" {
				skipped++
				fromId = photo.ID
				continue
			}

			if _, err := ps.GenerateThumbnailForPhoto(photo.ID); err != nil {
				slog.Warn("Failed to regenerate thumbnail", "photo", photo.ID, "error", err)
			}
			total++
			fromId = photo.ID
		}
	}

	slog.Info("Thumbnail regeneration complete", "generated", total, "skipped", skipped)
}

func (ps *PhotoService) GenerateThumbnailForPhoto(photoId string) (string, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return "", err
	}
	if photo.FilesystemPath == "" {
		return "", domain.ErrDataNotFound
	}

	unescapedPath := utils.UnescapeInvalidCharacters(photo.FilesystemPath)
	if _, statErr := os.Stat(unescapedPath); os.IsNotExist(statErr) {
		return "", domain.ErrDataNotFound
	}

	sizes := map[string][2]int{
		"s": {200, 200},
		"m": {600, 600},
		"l": {1200, 1200},
	}

	anyGenerated := false
	for sizeCode, dims := range sizes {
		thumbnail := ps.filesystemSvc.GenerateThumbnail(unescapedPath, exif.Exif{}, dims[0], dims[1])
		if len(thumbnail) == 0 {
			if !utils.IsVideoFile(unescapedPath) {
				_ = ps.photoRepo.HidePhoto(photoId)
			}
			break
		}

		anyGenerated = true
		if err := ps.thumbnailStorage.Put(context.Background(), photoId, sizeCode, thumbnail); err != nil {
			slog.Error("Failed to store thumbnail", "photo", photoId, "size", sizeCode, "error", err)
		}
		// Cache the path for filesystem access (used by response handlers that serve via file)
		cacheKey := fmt.Sprintf("thumbnail:%s:%s", photoId, sizeCode)
		_ = ps.cacheSvc.Set(cacheKey, "", 24*time.Hour)
		if sizeCode == "m" {
			// Encode blurhash from medium thumbnail for grid placeholders
			if img, _, err := image.Decode(bytes.NewReader(thumbnail)); err == nil {
				if blurhashStr, err := blurhash.Encode(4, 3, img); err == nil {
					photo.Blurhash = blurhashStr
				} else {
					slog.Warn("Failed to encode blurhash", "photo", photoId, "error", err)
				}
			} else {
				slog.Warn("Failed to decode thumbnail for blurhash", "photo", photoId, "error", err)
			}
		}
	}

	if !anyGenerated {
		return "", domain.ErrDataNotFound
	}

	// Update DB
	_ = ps.photoRepo.UpdatePhoto(*photo)

	return "", nil
}

func (ps *PhotoService) PhotoThumbnailPathForSize(photoId string, size string) (string, error) {
	if size == "" {
		size = "m"
	}

	// Check cache first
	cacheKey := fmt.Sprintf("thumbnail:%s:%s", photoId, size)
	if cachedPath, err := ps.cacheSvc.Get(cacheKey); err == nil && cachedPath != "" {
		if _, err := os.Stat(cachedPath); err == nil {
			return cachedPath, nil
		}
	}

	// Check filesystem
	safeId := strings.ReplaceAll(photoId, "/", "_")
	safeId = strings.ReplaceAll(safeId, "+", "-")
	safeId = strings.ReplaceAll(safeId, "=", "")
	path := filepath.Join(ps.config.PhotoDir, ".thumbnails", size, safeId+".webp")
	if _, err := os.Stat(path); err == nil {
		_ = ps.cacheSvc.Set(cacheKey, path, 24*time.Hour)
		return path, nil
	}

	// Check DB (legacy ThumbnailPath for medium size)
	if size == "m" {
		dbPath, err := ps.photoRepo.GetThumbnailPath(photoId)
		if err == nil && dbPath != "" {
			if _, err := os.Stat(dbPath); err == nil {
				_ = ps.cacheSvc.Set(cacheKey, dbPath, 24*time.Hour)
				return dbPath, nil
			}
		}
	}

	return "", domain.ErrDataNotFound
}

func (ps *PhotoService) PhotoThumbnailBytesForSize(photoId string, size string) ([]byte, error) {
	if size == "" {
		size = "m"
	}

	// Try thumbnail storage adapter first
	data, err := ps.thumbnailStorage.Get(context.Background(), photoId, size)
	if err == nil && len(data) > 0 {
		return data, nil
	}
	// Fallback: try filesystem path (legacy or direct file access)
	path, err := ps.PhotoThumbnailPathForSize(photoId, size)
	if err == nil && path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
	}

	return nil, domain.ErrDataNotFound
}

func (ps *PhotoService) DeletePhoto(photoId string) error {
	photo, err := ps.photoRepo.SoftDeletePhoto(photoId)
	if err != nil {
		return err
	}
	if photo.FilesystemPath == "" {
		return nil
	}
	unescapedPath := utils.UnescapeInvalidCharacters(photo.FilesystemPath)
	trashPath, err := ps.filesystemSvc.MoveToTrash(unescapedPath)
	if err != nil {
		return err
	}
	// Remember where the file went so RestorePhoto can move it back.
	return ps.photoRepo.SetTrashPath(photoId, trashPath)
}

// BatchDeletePhotos soft-deletes (moves to trash) multiple photos in one pass.
// A single photo failing does not abort the batch; the error is logged.
func (ps *PhotoService) BatchDeletePhotos(photoIds []string) error {
	for _, photoId := range photoIds {
		if err := ps.DeletePhoto(photoId); err != nil {
			slog.Warn("Failed to delete photo in batch", "id", photoId, "error", err)
		}
	}
	return nil
}

// BatchSetFavorite sets the favorite flag for multiple photos in a single
// database query.
func (ps *PhotoService) BatchSetFavorite(photoIds []string, favorite bool) error {
	return ps.photoRepo.SetFavoriteForPhotos(photoIds, favorite)
}

// BatchAddToAlbum assigns multiple photos to an album in a single query.
func (ps *PhotoService) BatchAddToAlbum(photoIds []string, albumId string) error {
	return ps.photoRepo.AssignPhotosToAlbum(photoIds, albumId)
}

func (ps *PhotoService) RestorePhoto(photoId string) error {
	photo, err := ps.photoRepo.RestorePhoto(photoId)
	if err != nil {
		return err
	}
	if photo.FilesystemPath == "" {
		return nil
	}
	originalPath := utils.UnescapeInvalidCharacters(photo.FilesystemPath)
	trashPath := utils.UnescapeInvalidCharacters(photo.TrashPath)
	if trashPath == "" {
		// Legacy photos trashed before TrashPath was tracked: the file's
		// original path is the best guess (it may already be in .trash).
		return nil
	}
	if err := ps.filesystemSvc.RestoreFromTrash(trashPath, originalPath); err != nil {
		return err
	}
	// File is back home; clear the recorded trash location.
	return ps.photoRepo.SetTrashPath(photoId, "")
}

func (ps *PhotoService) ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListTrashPhotos(fromId, limit, includeThumbnail)
	if err != nil {
		return nil, err
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}

func (ps *PhotoService) EmptyTrash() error {
	// Fetch all trashed photos so we can remove their files, thumbnails and
	// cache entries before hard-deleting the DB rows.
	trashed, err := ps.photoRepo.ListTrashPhotos("", 10000, false)
	if err != nil {
		return err
	}
	return ps.purgeTrashPhotos(trashed)
}

// PurgeExpiredTrash permanently deletes photos that were soft-deleted before
// the given cutoff (used by the trash retention job). Returns the number of
// photos permanently deleted.
func (ps *PhotoService) PurgeExpiredTrash(cutoff time.Time) (int, error) {
	expired, err := ps.photoRepo.ListExpiredTrashPhotos(cutoff)
	if err != nil {
		return 0, err
	}
	if len(expired) == 0 {
		return 0, nil
	}
	if err := ps.purgeTrashPhotos(expired); err != nil {
		return 0, err
	}
	return len(expired), nil
}

// purgeTrashPhotos permanently deletes the given trashed photos: original
// file, thumbnails and cache entries, then the DB row itself.
func (ps *PhotoService) purgeTrashPhotos(photos []*domain.Photo) error {
	photoIds := make([]string, 0, len(photos))
	for _, photo := range photos {
		if photo == nil {
			continue
		}
		photoIds = append(photoIds, photo.ID)

		// Original file (may live in .trash or at its original path).
		if photo.FilesystemPath != "" {
			unescapedPath := utils.UnescapeInvalidCharacters(photo.FilesystemPath)
			_ = ps.filesystemSvc.PermanentlyDelete(unescapedPath)
			_ = ps.filesystemSvc.PermanentlyDeleteTrashFile(unescapedPath)
		}

		// Thumbnails from storage (filesystem or S3).
		_ = ps.thumbnailStorage.Delete(context.Background(), photo.ID)

		// Cached thumbnail paths.
		for _, size := range []string{"s", "m", "l"} {
			_ = ps.cacheSvc.Delete(fmt.Sprintf("thumbnail:%s:%s", photo.ID, size))
		}
	}

	if len(photoIds) == 0 {
		return nil
	}
	return ps.photoRepo.PermanentlyDeletePhotos(photoIds)
}

func (ps *PhotoService) SetFavorite(photoId string, favorite bool) (*domain.Photo, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}
	photo.Favorite = favorite
	err = ps.photoRepo.SetFavorite(photoId, favorite)
	if err != nil {
		return nil, err
	}
	ps.setPhotoSourcePath(photo)
	return photo, nil
}

// UpdatePhotoMetadata persists user-editable metadata overrides
// (description, GPS location, date taken). The date override is converted
// to created_epoch + year/month so the photo re-sorts in the timeline.
// Original EXIF in the metadata JSONB is preserved (overrides are stored
// on dedicated columns).
func (ps *PhotoService) UpdatePhotoMetadata(photoId string, description string, latitude *float64, longitude *float64, dateTaken *string) (*domain.Photo, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}

	updates := domain.Photo{}

	if description != "" {
		updates.Description = description
		photo.Description = description
	}
	if latitude != nil {
		updates.Latitude = *latitude
		photo.Latitude = *latitude
	}
	if longitude != nil {
		updates.Longitude = *longitude
		photo.Longitude = *longitude
	}
	if dateTaken != nil {
		t, err := time.Parse("2006-01-02T15:04:05Z07:00", *dateTaken)
		if err != nil {
			t, err = time.Parse(time.RFC3339, *dateTaken)
			if err != nil {
				return nil, domain.ErrInvalidRequest
			}
		}
		updates.CreatedEpoch = t.UnixMilli()
		updates.Year = t.Year()
		updates.Month = int(t.Month())
		photo.CreatedEpoch = updates.CreatedEpoch
		photo.Year = updates.Year
		photo.Month = updates.Month
	}

	err = ps.photoRepo.UpdatePhotoMetadata(photoId, updates)
	if err != nil {
		return nil, err
	}
	ps.setPhotoSourcePath(photo)
	return photo, nil
}

func (ps *PhotoService) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListFavoritePhotos(fromId, limit, includeThumbnail, startDate, endDate)
	if err != nil {
		return nil, err
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}

func (ps *PhotoService) ListLowQualityPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListLowQualityPhotos(fromId, limit, includeThumbnail)
	if err != nil {
		return nil, err
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}

// ScorePhotoQuality runs the deterministic quality analyzer on a photo and
// stores the result. Returns (false, nil) when the photo is unscorable
// (unreadable file) or already scored.
func (ps *PhotoService) ScorePhotoQuality(photoId string) (bool, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return false, err
	}
	if photo.QualityScore != 0 || photo.FilesystemPath == "" {
		return false, nil
	}
	analysis, err := AnalyzeQuality(utils.UnescapeInvalidCharacters(photo.FilesystemPath))
	if err != nil || analysis == nil {
		return false, nil
	}
	if err := ps.photoRepo.UpdatePhotoQuality(photoId, analysis.Overall, analysis.BlurScore, analysis.IsLowQuality); err != nil {
		return false, err
	}
	return true, nil
}

// ScoreAllPhotoQuality backfills quality scores for all unscored photos
// (the index only scores newly-added files, so existing libraries need a
// one-off pass). Returns the number of photos scored.
func (ps *PhotoService) ScoreAllPhotoQuality(batchSize int) (int, error) {
	if batchSize <= 0 {
		batchSize = 200
	}
	scored := 0
	for {
		unscored, err := ps.photoRepo.ListUnscoredPhotos(batchSize)
		if err != nil {
			return scored, err
		}
		if len(unscored) == 0 {
			break
		}
		for _, photo := range unscored {
			ok, err := ps.ScorePhotoQuality(photo.ID)
			if err != nil {
				slog.Warn("Quality scoring failed", "photo", photo.ID, "error", err)
				continue
			}
			if ok {
				scored++
			}
		}
		if len(unscored) < batchSize {
			break
		}
	}
	return scored, nil
}

func (ps *PhotoService) Search(query string, limit int) ([]*domain.Photo, error) {
	return ps.photoRepo.SearchPhotos(query, limit)
}

func (ps *PhotoService) GetTimeline() ([]domain.TimelineEntry, error) {
	return ps.photoRepo.GetTimeline()
}

func (ps *PhotoService) GetGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error) {
	return ps.photoRepo.GetPhotosWithGeodata(north, south, east, west)
}

func (ps *PhotoService) RotatePhoto(photoId string, direction string) (*domain.Photo, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}
	// TODO: implement actual EXIF rotation or re-encode
	_ = direction
	return photo, nil
}

func (ps *PhotoService) DownloadPhotos(photoIds []string, writer io.Writer) error {
	zipWriter := archivezip.NewWriter(writer)
	defer zipWriter.Close()

	for _, id := range photoIds {
		photo, err := ps.photoRepo.GetPhotoById(id, false)
		if err != nil {
			slog.Warn("Skipping missing photo in download", "id", id, "error", err)
			continue
		}
		file, err := os.Open(utils.UnescapeInvalidCharacters(photo.FilesystemPath))
		if err != nil {
			slog.Warn("Unable to open photo for download", "path", photo.FilesystemPath, "error", err)
			continue
		}
		w, err := zipWriter.Create(photo.Name)
		if err != nil {
			file.Close()
			continue
		}
		_, err = io.Copy(w, file)
		file.Close()
		if err != nil {
			slog.Warn("Error copying photo to zip", "path", photo.FilesystemPath, "error", err)
		}
	}
	return nil
}

func (ps *PhotoService) GetAllTags() ([]string, error) {
	return ps.photoRepo.GetAllTags()
}

func (ps *PhotoService) ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListPhotosByTags(tags, fromId, limit, includeThumbnail)
	if err != nil {
		return nil, err
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}

func (ps *PhotoService) UpdatePhotoTags(photoId string, tags []string) (*domain.Photo, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}
	tagsStr := strings.Join(tags, ",")
	err = ps.photoRepo.UpdatePhotoTags(photoId, tagsStr)
	if err != nil {
		return nil, err
	}
	ps.setPhotoSourcePath(photo)
	return photo, nil
}

func (ps *PhotoService) BatchUpdatePhotoTags(photoIds []string, tags []string, operation string) error {
	for _, photoId := range photoIds {
		_, err := ps.photoRepo.GetPhotoById(photoId, false)
		if err != nil {
			continue
		}
		existingTagsList, err := ps.photoRepo.GetPhotoTags(photoId)
		if err != nil {
			continue
		}
		existingTags := make(map[string]struct{})
		for _, t := range existingTagsList {
			t = strings.TrimSpace(t)
			if t != "" {
				existingTags[t] = struct{}{}
			}
		}
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			switch operation {
			case "add":
				existingTags[tag] = struct{}{}
			case "remove":
				delete(existingTags, tag)
			case "set":
				existingTags = map[string]struct{}{tag: {}}
			}
		}
		if operation == "set" && len(tags) > 1 {
			existingTags = make(map[string]struct{})
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					existingTags[tag] = struct{}{}
				}
			}
		}
		newTags := make([]string, 0, len(existingTags))
		for t := range existingTags {
			newTags = append(newTags, t)
		}
		tagsStr := strings.Join(newTags, ",")
		_ = ps.photoRepo.UpdatePhotoTags(photoId, tagsStr)
	}
	return nil
}

func (ps *PhotoService) GetDuplicatePhotos() ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.GetDuplicatePhotos()
	if err != nil {
		return nil, err
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}
