package database

import (
	"encoding/json"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
)

func GetPhotoInfoById(photoId string) (Photo, error) {
	var photo Photo
	photo.ID = photoId
	result := dbConn.First(&photo)
	return photo, result.Error
}

func GetAllPhotos(pageNumber int, pageSize int, includeThumbnails bool) ([]Photo, error) {
	var photos []Photo
	result := dbConn.Scopes(Paginate(pageNumber, pageSize)).Find(&photos)
	if !includeThumbnails {
		return removeThumbnails(photos), result.Error
	}
	return photos, result.Error
}

func GetAllPhotosInfoInAlbum(albumId string, pageNumber int, pageSize int, includeThumbnails bool) ([]Photo, error) {
	var photos []Photo
	result := dbConn.Scopes(Paginate(pageNumber, pageSize)).Find(&photos, "album_id = ?", albumId)
	if !includeThumbnails {
		return removeThumbnails(photos), result.Error
	}
	return photos, result.Error
}

func GetPhotosInAlbumCount(albumId string) (int64, error) {
	var photos []Photo
	result := dbConn.Find(&photos, "album_id = ?", albumId)
	return result.RowsAffected, result.Error
}

func CreatePhoto(photo Photo) error {
	result := dbConn.Clauses(clause.OnConflict{
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