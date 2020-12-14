package database

import (
	"github.com/google/uuid"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
)

func GetAlbumById(albumId string) (Album, error) {
	var album Album
	album.ID = albumId
	result := dbConn.First(&album)
	return album, result.Error
}

func GetAlbumByName(albumName string) (Album, error) {
	var album Album
	result := dbConn.First(&album, "name = ?", albumName)
	return album, result.Error
}

func GetAllAlbums(page int, limit int) ([]Album, error) {
	var albums []Album
	result := dbConn.Scopes(Paginate(page, limit)).Find(&albums)
	return albums, result.Error
}

func GetAlbumCount() (int64, error) {
	var albums []Album
	result := dbConn.Find(&albums)
	return result.RowsAffected, result.Error
}

func CreateAlbum(name string) (string, error) {
	var album Album
	album.ID = uuid.New().String()
	album.Name = name
	result := dbConn.Create(&album)
	return album.ID, result.Error
}
