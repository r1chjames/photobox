package service

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/datatypes"
	"time"
)

type AlbumServiceMock struct{}

// NewAlbumServiceMock creates a new Album service instance
func NewAlbumServiceMock() *AlbumServiceMock {
	return &AlbumServiceMock{}
}

func (as *AlbumServiceMock) GetAlbum(id string) (*domain.Album, error) {
	return &domain.Album{
		ID:          id,
		Name:        "Name",
		Description: "Description",
		Tags:        "a",
		Metadata:    make(datatypes.JSON, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Thumbnail:   "tn",
	}, nil
}

func (as *AlbumServiceMock) ListAlbums(page, limit int) ([]*domain.Album, error) {
	return []*domain.Album{
		{
			ID:          "id1",
			Name:        "Name1",
			Description: "Description",
			Tags:        "a",
			Metadata:    make(datatypes.JSON, 0),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Thumbnail:   "tn",
		},
		{
			ID:          "id2",
			Name:        "Name2",
			Description: "Description",
			Tags:        "b",
			Metadata:    make(datatypes.JSON, 0),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Thumbnail:   "tn",
		},
	}, nil
}

func (as *AlbumServiceMock) AlbumCount() (int64, error) {
	return 2, nil
}
