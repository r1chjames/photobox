package repository

import (
	. "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/datatypes"
	"time"
)

type AlbumRepositoryMock struct {
	dbEnv *Env
}

func NewAlbumRepositoryMock() *AlbumRepositoryMock {
	return &AlbumRepositoryMock{}
}

func (ar *AlbumRepositoryMock) GetAlbumById(id string) (*domain.Album, error) {
	var album = domain.Album{
		ID:          id,
		Name:        "Name",
		Description: "Description",
		Tags:        "a",
		Metadata:    make(datatypes.JSON, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Thumnail:    "tn",
	}
	return &album, nil
}

func (ar *AlbumRepositoryMock) GetAlbumByName(name string) (*domain.Album, error) {
	var album = domain.Album{
		ID:          "1",
		Name:        name,
		Description: "Description",
		Tags:        "a",
		Metadata:    make(datatypes.JSON, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Thumnail:    "tn",
	}
	return &album, nil
}

func (ar *AlbumRepositoryMock) CreateAlbum(name string) (*domain.Album, error) {
	var album = domain.Album{
		ID:          "1",
		Name:        name,
		Description: "Description",
		Tags:        "a",
		Metadata:    make(datatypes.JSON, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Thumnail:    "tn",
	}
	return &album, nil
}

func (ar *AlbumRepositoryMock) CreateAlbumIfNotExists(name string) (*domain.Album, error) {
	var album = domain.Album{
		ID:          "1",
		Name:        name,
		Description: "Description",
		Tags:        "a",
		Metadata:    make(datatypes.JSON, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Thumnail:    "tn",
	}
	return &album, nil
}

func (ar *AlbumRepositoryMock) ListAllAlbums(page int, limit int) ([]*domain.Album, error) {
	return []*domain.Album{
		{
			ID:          "id1",
			Name:        "Name1",
			Description: "Description",
			Tags:        "a",
			Metadata:    make(datatypes.JSON, 0),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Thumnail:    "tn",
		},
		{
			ID:          "id2",
			Name:        "Name2",
			Description: "Description",
			Tags:        "b",
			Metadata:    make(datatypes.JSON, 0),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Thumnail:    "tn",
		},
	}, nil
}

func (ar *AlbumRepositoryMock) AlbumCount() (int64, error) {
	return 2, nil
}
