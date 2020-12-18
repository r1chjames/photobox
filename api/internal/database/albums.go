package database

import (
	"github.com/google/uuid"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
)

func (dbEnv *Env) GetAlbumById(albumId string) (Album, error) {
	var album Album
	album.ID = albumId
	result := dbEnv.db.First(&album)
	return album, result.Error
}

func (dbEnv *Env) GetAlbumByName(albumName string) (Album, error) {
	var album Album
	result := dbEnv.db.First(&album, "name = ?", albumName)
	return album, result.Error
}

func (dbEnv *Env) GetAllAlbums(page int, limit int) ([]Album, error) {
	var albums []Album
	result := dbEnv.db.Scopes(Paginate(page, limit)).Find(&albums)
	return albums, result.Error
}

func (dbEnv *Env) GetAlbumCount() (int64, error) {
	var albums []Album
	result := dbEnv.db.Find(&albums)
	return result.RowsAffected, result.Error
}

func (dbEnv *Env) CreateAlbum(name string) (string, error) {
	var album Album
	album.ID = uuid.New().String()
	album.Name = name
	result := dbEnv.db.Create(&album)
	return album.ID, result.Error
}
