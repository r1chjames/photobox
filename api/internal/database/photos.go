package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
)

func (dbEnv *Env) GetPhotoInfoById(photoId string, includeThumbnail bool) (Photo, error) {
	var photo Photo
	photo.ID = photoId
	result := dbEnv.Db
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.First(&photo)
	return photo, result.Error
}

func (dbEnv *Env) GetAllPhotos(pageNumber int, pageSize int, includeThumbnail bool) ([]Photo, error) {
	var photos []Photo
	result := dbEnv.Db.Scopes(Paginate(pageNumber, pageSize))
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.Find(&photos)
	return photos, result.Error
}

func (dbEnv *Env) GetAllPhotosInfoInAlbum(albumId string, pageNumber int, pageSize int, includeThumbnail bool) ([]Photo, error) {
	var photos []Photo
	result := dbEnv.Db.Scopes(Paginate(pageNumber, pageSize))
	if !includeThumbnail {
		result.Omit("thumbnail")
	}
	result.Find(&photos, "album_id = ?", albumId)
	return photos, result.Error
}

func (dbEnv *Env) GetPhotosInAlbumCount(albumId string) (int64, error) {
	var count int64
	result := dbEnv.Db.Model(&[]Photo{}).Where("album_id = ?", albumId).Count(&count)
	return count, result.Error
}

func (dbEnv *Env) CreatePhotoInfo(photo Photo) error {
	result := dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&photo)
	return result.Error
}
