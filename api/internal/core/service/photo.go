package service

import (
	"fmt"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type PhotoService struct {
	repo   port.PhotoRepository
	config appconfig.AppConfig
}

// NewPhotoService creates a new Photo service instance
func NewPhotoService(repo port.PhotoRepository, config appconfig.AppConfig) *PhotoService {
	return &PhotoService{
		repo,
		config,
	}
}

func (ps *PhotoService) ListPhotosInAlbum(albumId string, page, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.repo.ListAllPhotosInAlbum(albumId, page, limit, includeThumbnail)
	ps.setPhotosSourcePath(resp)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	return resp, nil
}

func (ps *PhotoService) ListPhotos(page, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.repo.ListAllPhotos(page, limit, includeThumbnail)
	ps.setPhotosSourcePath(resp)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	return resp, nil
}

func (ps *PhotoService) GetPhoto(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	resp, err := ps.repo.GetPhotoById(photoId, includeThumbnail)
	ps.setPhotoSourcePath(resp)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	return resp, nil
}

//func addPhoto(c *gin.Context) {
//	var photo PhotoUpload
//	err := c.BindJSON(&photo)
//
//	photoFile := components.WriteFileToFilesystem(dbEnv, photo)
//	dbEnv.SavePhotoRecordsToDatabase([]PhotoFile{photoFile})
//
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, invalidRequest()})
//	} else {
//		c.Status(http.StatusCreated)
//	}
//}

func (ps *PhotoService) PhotoCount(albumId string) (int64, error) {
	return ps.repo.GetPhotosInAlbumCount(albumId)
}

func (ps *PhotoService) PhotoBinary(photoId string) (string, error) {

	photoInfo, err := ps.repo.GetPhotoById(photoId, false)
	if err != nil {
		return "", err
	}

	return photoInfo.FilesystemPath, nil
}

func (ps *PhotoService) PhotoThumbnail(photoId string) ([]byte, error) {

	photoInfo, err := ps.repo.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}
	return photoInfo.Thumbnail, nil

}

func (ps *PhotoService) setPhotoSourcePath(photo *domain.Photo) {
	photo.SourcePath = fmt.Sprintf("%s/photo/%s}/bin", ps.config.ApiBasePath, photo.ID)
}

func (ps *PhotoService) setPhotosSourcePath(photos []*domain.Photo) {
	for _, photo := range photos {
		ps.setPhotoSourcePath(photo)
	}
}
