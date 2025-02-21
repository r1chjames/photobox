package repository

import (
	b64 "encoding/base64"
	"github.com/google/uuid"
	. "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
	"time"
)

type AlbumRepository struct {
	dbEnv *Env
}

func NewAlbumRepository(dbEnv *Env) *AlbumRepository {
	return &AlbumRepository{
		dbEnv,
	}
}

func (ar *AlbumRepository) GetAlbumById(id string) (*domain.Album, error) {
	var album domain.Album
	album.ID = id
	result := ar.dbEnv.Db.First(&album)
	err := HandleError(result)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (ar *AlbumRepository) GetAlbumByName(name string) (*domain.Album, error) {
	var album *domain.Album
	result := ar.dbEnv.Db.First(&album, "name = ?", name)
	err := HandleError(result)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (ar *AlbumRepository) CreateAlbum(name string) (*domain.Album, error) {
	var album domain.Album
	album.ID = uuid.New().String()
	album.Name = name
	album.CreatedEpoch = time.Now().UnixMilli()
	ar.dbEnv.Db.Create(&album)
	return &album, nil
}

func (ar *AlbumRepository) CreateAlbumIfNotExists(name string) (*domain.Album, error) {
	var album domain.Album
	album.ID = uuid.New().String()
	album.Name = name
	album.CreatedEpoch = time.Now().UnixMilli()

	ar.dbEnv.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&album)
	return &album, nil
}

func (ar *AlbumRepository) ListAllAlbums(fromId string, limit int) ([]*domain.Album, error) {
	var albums []*domain.Album
	result := ar.dbEnv.Db.Model(&[]domain.Album{}).Limit(limit)
	if fromId != "" {
		fromEpoch, _ := b64.StdEncoding.DecodeString(fromId)
		result.Where("created_epoch > ?", fromEpoch)
	}
	result.Find(&albums)
	err := HandleError(result)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (ar *AlbumRepository) AlbumCount() (int64, error) {
	var albums []domain.Album
	result := ar.dbEnv.Db.Find(&albums)
	err := HandleError(result)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}
