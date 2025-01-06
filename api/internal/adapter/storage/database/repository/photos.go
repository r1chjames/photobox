package repository

import (
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
)

type PhotoRepository struct {
	dbEnv *database.Env
}

func NewPhotoRepository(dbEnv *database.Env) *PhotoRepository {
	return &PhotoRepository{
		dbEnv,
	}
}

func (pr *PhotoRepository) GetPhotoById(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	var photo *domain.Photo
	photo.ID = photoId
	result := pr.dbEnv.Db
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.First(&photo)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photo, result.Error
}

func (pr *PhotoRepository) ListAllPhotos(pageNumber int, pageSize int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Scopes(database.Paginate(pageNumber, pageSize))
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, result.Error
}

func (pr *PhotoRepository) ListAllPhotosInAlbum(albumId string, pageNumber int, pageSize int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Scopes(database.Paginate(pageNumber, pageSize))
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.Find(&photos, "album_id = ?", albumId)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, result.Error
}

func (pr *PhotoRepository) GetPhotosInAlbumCount(albumId string) (int64, error) {
	var count int64
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("album_id = ?", albumId).Count(&count)
	if result.RowsAffected == 0 {
		return 0, nil
	}
	return count, result.Error
}

func (pr *PhotoRepository) CreatePhotoInfo(photo domain.Photo) error {
	result := pr.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&photo)
	return result.Error
}
