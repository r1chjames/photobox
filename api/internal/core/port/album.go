package port

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

//go:generate mockgen -source=album.go -destination=mock/album.go -package=mock

// AlbumRepository is an interface for interacting with Album-related data
type AlbumRepository interface {
	// GetAlbumById returns an album using its ID
	GetAlbumById(id string) (*domain.Album, error)
	// GetAlbumByName returns an album using its name
	GetAlbumByName(name string) (*domain.Album, error)
	// CreateAlbum creates an album
	CreateAlbum(name string) (*domain.Album, error)
	// CreateAlbumIfNotExists creates an album if it doesn't exist
	CreateAlbumIfNotExists(name string) (*domain.Album, error)
	// ListAllAlbums returns all stored albums
	ListAllAlbums(pageNumber int, pageSize int) ([]*domain.Album, error)
	// AlbumCount returns a count of all albums
	AlbumCount() (int64, error)
}

// AlbumService is an interface for interacting with Album-related business logic
type AlbumService interface {
	//GetAlbum returns a Album
	GetAlbum(id string) (*domain.Album, error)
	// ListAlbums returns all albums
	ListAlbums(page, limit int) ([]*domain.Album, error)
	// AlbumCount returns a count of all albums
	AlbumCount() (int64, error)
}
