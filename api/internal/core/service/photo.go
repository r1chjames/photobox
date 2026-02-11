package service

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

type PhotoService struct {
	photoRepo     port.PhotoRepository
	albumSvc      port.AlbumService
	filesystemSvc port.FilesystemService
	config        appconfig.AppConfig
}

// NewPhotoService creates a new Photo service instance
func NewPhotoService(photoRepo port.PhotoRepository, albumRepo port.AlbumService, filesystemSvc port.FilesystemService, config appconfig.AppConfig) *PhotoService {
	return &PhotoService{
		photoRepo,
		albumRepo,
		filesystemSvc,
		config,
	}
}

func (ps *PhotoService) ListPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListAllPhotosInAlbum(albumId, fromId, limit, includeThumbnail)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	ps.setPhotosSourcePath(resp)
	return resp, nil
}

func (ps *PhotoService) ListPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListAllPhotos(fromId, limit, includeThumbnail)
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
	ps.filesystemSvc.PerformPhotoIndex(ps.SavePhoto)
}

func (ps *PhotoService) PhotoCount(albumId string) (int64, error) {
	return ps.photoRepo.GetPhotosInAlbumCount(albumId)
}

func (ps *PhotoService) PhotoBinary(photoId string) (string, error) {
	photoInfo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return "", err
	}
	return photoInfo.FilesystemPath, nil
}

func (ps *PhotoService) PhotoThumbnail(photoId string) ([]byte, error) {
	photoInfo, err := ps.photoRepo.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}
	return photoInfo.Thumbnail, nil
}

func (ps *PhotoService) setPhotoSourcePath(photo *domain.Photo) {
	photo.SourcePath = fmt.Sprintf("photo/%s/bin", photo.ID)
}

func (ps *PhotoService) setPhotosSourcePath(photos []*domain.Photo) {
	for _, photo := range photos {
		ps.setPhotoSourcePath(photo)
	}
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
		}

		slog.Info("Adding photo", "photo", photo.Name, "album", photo.Directory)
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
	}

	slog.Info("Adding photo", "photo", photo.Name, "album", photo.Directory)

	return ps.photoRepo.CreatePhotoInfo(photoInfo)
}
