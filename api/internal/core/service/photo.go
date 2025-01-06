package service

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
	"log"
)

type PhotoService struct {
	photoRepo port.PhotoRepository
	albumRepo port.AlbumRepository
	config    appconfig.AppConfig
}

// NewPhotoService creates a new Photo service instance
func NewPhotoService(photoRepo port.PhotoRepository, albumRepo port.AlbumRepository, config appconfig.AppConfig) *PhotoService {
	return &PhotoService{
		photoRepo,
		albumRepo,
		config,
	}
}

func (ps *PhotoService) ListPhotosInAlbum(albumId string, page, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListAllPhotosInAlbum(albumId, page, limit, includeThumbnail)
	ps.setPhotosSourcePath(resp)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	return resp, nil
}

func (ps *PhotoService) ListPhotos(page, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	resp, err := ps.photoRepo.ListAllPhotos(page, limit, includeThumbnail)
	ps.setPhotosSourcePath(resp)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}
	return resp, nil
}

func (ps *PhotoService) GetPhoto(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	resp, err := ps.photoRepo.GetPhotoById(photoId, includeThumbnail)
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
	photo.SourcePath = fmt.Sprintf("%s/photo/%s}/bin", ps.config.ApiBasePath, photo.ID)
}

func (ps *PhotoService) setPhotosSourcePath(photos []*domain.Photo) {
	for _, photo := range photos {
		ps.setPhotoSourcePath(photo)
	}
}

func (ps *PhotoService) SavePhotos(photos []domain.PhotoFile) error {
	for _, photo := range photos {
		err := ps.SavePhoto(photo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PhotoService) SavePhoto(photo domain.PhotoFile) error {
	result, err := ps.albumRepo.GetAlbumByName(photo.Directory)
	if err != nil {
		albumId, albErr := ps.albumRepo.CreateAlbum(photo.Directory)
		if albErr != nil {
			log.Printf("unable to insert album record, %s", albErr)
		}
		log.Printf("created album name: %s, id: %s", photo.Directory, albumId)
	}
	photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
	photoMetadata, _ := json.Marshal(&photo)
	photoInfo := domain.Photo{
		ID:             photoHash,
		Name:           utils.EscapeInvalidCharacters(photo.Name),
		FilesystemPath: utils.EscapeInvalidCharacters(photo.Path),
		AlbumId:        result.ID,
		Tags:           "",
		Metadata:       photoMetadata,
		Thumbnail:      photo.Thumbnail,
	}

	log.Printf("Adding photo: %s to album: %s", photo.Name, photo.Directory)

	return ps.photoRepo.CreatePhotoInfo(photoInfo)
}
