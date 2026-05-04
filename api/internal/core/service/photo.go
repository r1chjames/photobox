package service

import (
	archivezip "archive/zip"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

type PhotoService struct {
	photoRepo     port.PhotoRepository
	albumSvc      port.AlbumService
	filesystemSvc port.FilesystemService
	cacheSvc      port.CacheService
	config        appconfig.AppConfig
}

// NewPhotoService creates a new Photo service instance
func NewPhotoService(photoRepo port.PhotoRepository, albumRepo port.AlbumService, filesystemSvc port.FilesystemService, cacheSvc port.CacheService, config appconfig.AppConfig) *PhotoService {
	return &PhotoService{
		photoRepo,
		albumRepo,
		filesystemSvc,
		cacheSvc,
		config,
	}
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

func (ps *PhotoService) PerformPhotoIndex() {
	cache, _ := ps.photoRepo.GetPhotoIndexCache()
	ps.filesystemSvc.PerformPhotoIndex(ps.SavePhotos, cache)
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
	absReq, _ := filepath.Abs(photoInfo.FilesystemPath)
	if !strings.HasPrefix(absReq, absBase) {
		return "", domain.ErrForbidden
	}
	return photoInfo.FilesystemPath, nil
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

func (ps *PhotoService) thumbnailPathForSize(photoId string, size string) string {
	safeId := strings.ReplaceAll(photoId, "/", "_")
	safeId = strings.ReplaceAll(safeId, "+", "-")
	safeId = strings.ReplaceAll(safeId, "=", "")
	return filepath.Join(ps.config.PhotoDir, ".thumbnails", size, safeId+".jpg")
}

func (ps *PhotoService) writeThumbnailToDisk(photoId string, data []byte, size string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	path := ps.thumbnailPathForSize(photoId, size)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
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
			FilesystemPath: utils.EscapeInvalidCharacters(photo.Path),
			AlbumId:        albumId,
			Tags:           "",
			Metadata:       photoMetadata,
			Thumbnail:      photo.Thumbnail,
			CreatedEpoch:   time.Now().UnixMilli(),
			FileHash:       computeFileHash(photo.Path),
			FileModifiedTime: photo.ModifiedTime,
			MediaType:      photo.MediaType,
			Duration:       photo.Duration,
			Width:          photo.Width,
			Height:         photo.Height,
			Latitude:       photo.Latitude,
			Longitude:      photo.Longitude,
		}

		slog.Info("Adding photo", "photo", photo.Name, "album", photo.Directory)

		if len(photo.Thumbnail) > 0 {
			thumbPath, err := ps.writeThumbnailToDisk(photoHash, photo.Thumbnail, "m")
			if err == nil {
				photoInfo.ThumbnailPath = thumbPath
			} else {
				slog.Error("Failed to write thumbnail to disk", "photo", photo.Name, "error", err)
			}
		}

		photoRecords = append(photoRecords, photoInfo)
	}

	// Batch insert all photos
	return ps.photoRepo.CreatePhotosInfo(photoRecords)
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
		FilesystemPath: utils.EscapeInvalidCharacters(photo.Path),
		AlbumId:        albumId,
		Tags:           "",
		Metadata:       photoMetadata,
		Thumbnail:      photo.Thumbnail,
		CreatedEpoch:   time.Now().UnixMilli(),
		FileHash:       computeFileHash(photo.Path),
		FileModifiedTime: photo.ModifiedTime,
		MediaType:      photo.MediaType,
		Duration:       photo.Duration,
		Width:          photo.Width,
		Height:         photo.Height,
		Latitude:       photo.Latitude,
		Longitude:      photo.Longitude,
	}

	slog.Info("Adding photo", "photo", photo.Name, "album", photo.Directory)

	if len(photo.Thumbnail) > 0 {
		thumbPath, err := ps.writeThumbnailToDisk(photoHash, photo.Thumbnail, "m")
		if err == nil {
			photoInfo.ThumbnailPath = thumbPath
		} else {
			slog.Error("Failed to write thumbnail to disk", "photo", photo.Name, "error", err)
		}
	}

	return ps.photoRepo.CreatePhotoInfo(photoInfo)
}

func (ps *PhotoService) GenerateThumbnailForPhoto(photoId string) (string, error) {
	photo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return "", err
	}
	if photo.FilesystemPath == "" {
		return "", domain.ErrDataNotFound
	}

	sizes := map[string][2]int{
		"s": {200, 200},
		"m": {600, 600},
		"l": {1200, 1200},
	}

	var mediumPath string
	for sizeCode, dims := range sizes {
		thumbnail := ps.filesystemSvc.GenerateThumbnail(photo.FilesystemPath, exif.Exif{}, dims[0], dims[1])
		if len(thumbnail) == 0 {
			continue
		}
		path, err := ps.writeThumbnailToDisk(photoId, thumbnail, sizeCode)
		if err != nil {
			slog.Error("Failed to write thumbnail", "size", sizeCode, "error", err)
			continue
		}
		if sizeCode == "m" {
			mediumPath = path
		}

		// Cache the path
		cacheKey := fmt.Sprintf("thumbnail:%s:%s", photoId, sizeCode)
		_ = ps.cacheSvc.Set(cacheKey, path, 24*time.Hour)
	}

	// Update DB
	photo.ThumbnailPath = mediumPath
	_ = ps.photoRepo.UpdatePhoto(*photo)

	return mediumPath, nil
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
	path := ps.thumbnailPathForSize(photoId, size)
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

func (ps *PhotoService) DeletePhoto(photoId string) error {
	photo, err := ps.photoRepo.SoftDeletePhoto(photoId)
	if err != nil {
		return err
	}
	if photo.FilesystemPath == "" {
		return nil
	}
	unescapedPath := utils.EscapeInvalidCharacters(photo.FilesystemPath)
	_, err = ps.filesystemSvc.MoveToTrash(unescapedPath)
	return err
}

func (ps *PhotoService) RestorePhoto(photoId string) error {
	photo, err := ps.photoRepo.RestorePhoto(photoId)
	if err != nil {
		return err
	}
	if photo.FilesystemPath == "" {
		return nil
	}
	trashPath := utils.EscapeInvalidCharacters(photo.FilesystemPath)
	// The original path is the same as the stored path before trash
	// The filesystem repo will move it back
	// For now, we assume the path in DB is the original path
	// and the trash path is computed by the repo
	_ = trashPath
	return nil
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
	return ps.photoRepo.EmptyTrash()
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

func (ps *PhotoService) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListFavoritePhotos(fromId, limit, includeThumbnail, startDate, endDate)
	if err != nil {
		return nil, err
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
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
		file, err := os.Open(photo.FilesystemPath)
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
	photo.Tags = tagsStr
	ps.setPhotoSourcePath(photo)
	return photo, nil
}

func (ps *PhotoService) BatchUpdatePhotoTags(photoIds []string, tags []string, operation string) error {
	for _, photoId := range photoIds {
		photo, err := ps.photoRepo.GetPhotoById(photoId, false)
		if err != nil {
			continue
		}
		existingTags := make(map[string]struct{})
		for _, t := range strings.Split(photo.Tags, ",") {
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
