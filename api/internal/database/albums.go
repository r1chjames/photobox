package database

import (
	"github.com/google/uuid"
)

func (dbEnv *Env) GetAlbumById(albumId string) (Album, error) {
	var album Album
	album.ID = albumId
	result := dbEnv.Db.First(&album)
	return album, result.Error
}

func (dbEnv *Env) GetAlbumByName(albumName string) (Album, error) {
	var album Album
	result := dbEnv.Db.First(&album, "name = ?", albumName)
	return album, result.Error
}

func (dbEnv *Env) CreateAlbumIfNotExists(albumName string) (Album, error) {
	album, err := dbEnv.GetAlbumByName(albumName)
	if checkNotFoundError(err) {
		_, err = dbEnv.CreateAlbum(albumName)
	}
	return album, nil
}

func (dbEnv *Env) GetAllAlbums(page int, limit int) ([]Album, error) {
	var albums []Album
	result := dbEnv.Db.Scopes(Paginate(page, limit)).Find(&albums)
	return albums, result.Error
}

func (dbEnv *Env) GetAlbumCount() (int64, error) {
	var albums []Album
	result := dbEnv.Db.Find(&albums)
	return result.RowsAffected, result.Error
}

func (dbEnv *Env) CreateAlbum(name string) (string, error) {
	var album Album
	album.ID = uuid.New().String()
	album.Name = name
	result := dbEnv.Db.Create(&album)
	return album.ID, result.Error
}
