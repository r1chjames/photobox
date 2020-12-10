package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
)

func GetPhotoInfoById(photoId string) (Photo, error) {
	var photo Photo
	result := dbConn.First(&photo, photoId)
	return photo, result.Error
}

func GetAllPhotos(pageNumber int, pageSize int) ([]Photo, error) {
	var photos []Photo
	result := dbConn.Scopes(Paginate(pageNumber, pageSize)).Find(&photos)
	return photos, result.Error
}

func GetAllPhotosInfoInAlbum(albumId string, pageNumber int, pageSize int) ([]Photo, error) {
	var photos []Photo
	result := dbConn.Scopes(Paginate(pageNumber, pageSize)).Find(&photos, "AlbumID = ?", albumId)
	return photos, result.Error
}

func GetPhotosInAlbumCount(albumId string) (int64, error) {
	var photos []Photo
	result := dbConn.Find(&photos, "AlbumID = ?", albumId)
	return result.RowsAffected, result.Error
}

func CreatePhoto(photo Photo) error {
	result := dbConn.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&photo)
	return result.Error
}

//result := dbConn.Clauses(clause.OnConflict{
//Columns:   []clause.Column{{Name: "id"}},
//DoUpdates: clause.Assignments(map[string]interface{}{
//"name": photo.Name,
//"filesystem_path": photo.FilesystemPath,
//"album_id": photo.AlbumId,
//"tags": photo.Tags,
//"metadata": photo.Metadata,
//}),