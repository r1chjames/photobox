package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
)

func (dbEnv *Env) GetPhotoInfoById(photoId string) (Photo, error) {
	var photo Photo
	photo.ID = photoId
	result := dbEnv.Db.First(&photo)
	return photo, result.Error
}

func (dbEnv *Env) GetAllPhotos(pageNumber int, pageSize int) ([]Photo, error) {
	var photos []Photo
	result := dbEnv.Db.Scopes(Paginate(pageNumber, pageSize)).Find(&photos)
	return photos, result.Error
}

func (dbEnv *Env) GetAllPhotosInfoInAlbum(albumId string, pageNumber int, pageSize int) ([]Photo, error) {
	var photos []Photo
	result := dbEnv.Db.Scopes(Paginate(pageNumber, pageSize)).Find(&photos, "album_id = ?", albumId)
	return photos, result.Error
}

func (dbEnv *Env) GetPhotosInAlbumCount(albumId string) (int64, error) {
	var photos []Photo
	result := dbEnv.Db.Find(&photos, "album_id = ?", albumId)
	return result.RowsAffected, result.Error
}

func (dbEnv *Env) CreatePhotoInfo(photo Photo) error {
	result := dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&photo)
	return result.Error
}