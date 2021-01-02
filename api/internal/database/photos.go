package database

import (
	"encoding/json"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
)

func (dbEnv *Env) GetPhotoInfoById(photoId string) (Photo, error) {
	var photo Photo
	photo.ID = photoId
	result := dbEnv.Db.First(&photo)
	return photo, result.Error
}

func (dbEnv *Env) GetAllPhotos(pageNumber int, pageSize int, includeThumbnails bool) ([]Photo, error) {
	var photos []Photo
	result := dbEnv.Db.Scopes(Paginate(pageNumber, pageSize)).Find(&photos)
	if !includeThumbnails {
		return removeThumbnails(photos), result.Error
	}
	return photos, result.Error
}

func (dbEnv *Env) GetAllPhotosInfoInAlbum(albumId string, pageNumber int, pageSize int, includeThumbnails bool) ([]Photo, error) {
	var photos []Photo
	result := dbEnv.Db.Scopes(Paginate(pageNumber, pageSize)).Find(&photos, "album_id = ?", albumId)
	if !includeThumbnails {
		return removeThumbnails(photos), result.Error
	}
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

func removeThumbnail(photo Photo) Photo {
	var photoMetadata PhotoFile
	retrievedPhotoMetadata := photo.Metadata
	_ = json.Unmarshal(retrievedPhotoMetadata, &photoMetadata)
	photoMetadata.Thumbnail = nil
	parsedMetadata, _ := json.Marshal(&photoMetadata)
	photo.Metadata = parsedMetadata
	return photo
}

func removeThumbnails(photos []Photo) []Photo {
	var processedPhotos []Photo
	for _, photo := range photos {
		processedPhotos = append(processedPhotos, removeThumbnail(photo))
	}
	return processedPhotos
}
