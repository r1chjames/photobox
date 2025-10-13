package repository

import (
	b64 "encoding/base64"

	. "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
)

type PhotoRepository struct {
	dbEnv *Env
}

func NewPhotoRepository(dbEnv *Env) *PhotoRepository {
	return &PhotoRepository{
		dbEnv,
	}
}

func (pr *PhotoRepository) GetPhotoById(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	var photo domain.Photo
	photo.ID = photoId
	result := pr.dbEnv.Db
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.First(&photo)
	err := HandleError(result)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (pr *PhotoRepository) ListAllPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Limit(limit)
	if fromId != "" {
		fromPhoto, _ := pr.GetPhotoById(fromId, false)
		result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, nil
}

func (pr *PhotoRepository) ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Limit(limit)
	if fromId != "" {
		fromEpoch, _ := b64.StdEncoding.DecodeString(fromId)
		result.Where("created_epoch > ?", fromEpoch)
	}
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.Find(&photos, "album_id = ?", albumId)
	err := HandleError(result)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func (pr *PhotoRepository) GetPhotosInAlbumCount(albumId string) (int64, error) {
	var count int64
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("album_id = ?", albumId).Count(&count)
	err := HandleError(result)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pr *PhotoRepository) CreatePhotoInfo(photo domain.Photo) error {
	result := pr.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&photo)
	return result.Error
}
