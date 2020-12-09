package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
)

func GetPhotoInfoById(photoId string) (Photo, error) {
	var photo Photo
	result := dbConn.First(&photo, photoId)
	return photo, result.Error
}

func GetAllPhotos() ([]Photo, error) {
	var photos []Photo
	result := dbConn.Find(&photos)
	return photos, result.Error
}

func GetAllPhotosInfoInAlbum(albumId string) ([]Photo, error) {
	var photos []Photo
	result := dbConn.Find(&photos, "AlbumID = ?", albumId)
	return photos, result.Error
}

func GetPhotosInAlbumCount(albumId string) (int64, error) {
	var photos []Photo
	result := dbConn.Find(&photos, "AlbumID = ?", albumId)
	return result.RowsAffected, result.Error
}

func CreatePhoto(photo Photo) error {
	result := dbConn.Save(photo)
	return result.Error
}