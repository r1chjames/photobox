package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
)

func GetAlbumById(albumId string) (Album, error) {
	var album Album
	result := dbConn.First(&album, albumId)
	return album, result.Error
}

func GetAlbumByName(albumName string) (Album, error) {
	var album Album
	result := dbConn.First(&album, "name = ?", albumName)
	return album, result.Error
}

func GetAllAlbums() ([]Album, error) {
	var albums []Album
	result := dbConn.Find(&albums)
	return albums, result.Error
}

func GetAlbumCount() (int64, error) {
	var albums []Album
	result := dbConn.Find(&albums)
	return result.RowsAffected, result.Error
}

func CreateAlbum(name string) error {
	var album Album
	album.Name = name
	result := dbConn.Save(album)
	return result.Error
}