package repository

import (
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

func (pr *PhotoRepository) ListAllPhotos(pageNumber, pageSize int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Scopes(Paginate(pageNumber, pageSize))
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, nil
}

func (pr *PhotoRepository) ListAllPhotosInAlbum(albumId string, pageNumber, pageSize int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Scopes(Paginate(pageNumber, pageSize))
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
