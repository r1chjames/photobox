package repository

import (
	"github.com/google/uuid"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
)

type AlbumRepository struct {
	dbEnv *database.Env
}

func NewAlbumRepository(dbEnv *database.Env) *AlbumRepository {
	return &AlbumRepository{
		dbEnv,
	}
}

func (ar *AlbumRepository) GetAlbumById(id string) (*domain.Album, error) {
	var album *domain.Album
	album.ID = id
	result := ar.dbEnv.Db.First(&album)
	return album, result.Error
}

func (ar *AlbumRepository) GetAlbumByName(name string) (*domain.Album, error) {
	var album *domain.Album
	result := ar.dbEnv.Db.First(&album, "name = ?", name)
	return album, result.Error
}

func (ar *AlbumRepository) CreateAlbum(name string) (*domain.Album, error) {
	var album *domain.Album
	album.ID = uuid.New().String()
	album.Name = name
	result := ar.dbEnv.Db.Create(&album)
	return album, result.Error
}

func (ar *AlbumRepository) CreateAlbumIfNotExists(name string) (*domain.Album, error) {
	var album *domain.Album
	album.ID = uuid.New().String()
	album.Name = name

	result := ar.dbEnv.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&album)
	return album, result.Error
}

func (ar *AlbumRepository) ListAllAlbums(page int, limit int) ([]*domain.Album, error) {
	var albums []*domain.Album
	result := ar.dbEnv.Db.Scopes(database.Paginate(page, limit)).Find(&albums)
	return albums, result.Error
}

func (ar *AlbumRepository) AlbumCount() (int64, error) {
	var albums []domain.Album
	result := ar.dbEnv.Db.Find(&albums)
	return result.RowsAffected, result.Error
}
