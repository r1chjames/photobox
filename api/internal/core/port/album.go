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
	// CreateAlbumWithMetadata creates an album with the given metadata
	CreateAlbumWithMetadata(album *domain.Album) error
	// CreateAlbumIfNotExists creates an album if it doesn't exist
	CreateAlbumIfNotExists(name string) (*domain.Album, error)
	// ListAllAlbums returns all stored albums
	ListAllAlbums(fromId string, pageSize int) ([]*domain.Album, error)
	// AlbumCount returns a count of all albums
	AlbumCount() (int64, error)
	// UpdateAlbum updates an album
	UpdateAlbum(album *domain.Album) error
	// DeleteAlbum deletes an album by id
	DeleteAlbum(id string) error
	// ReassignPhotosToAlbum moves all photos from one album to another
	ReassignPhotosToAlbum(fromAlbumId, toAlbumId string) error
	// SearchAlbums searches albums by query
	SearchAlbums(query string, limit int) ([]*domain.Album, error)
}

// AlbumService is an interface for interacting with Album-related business logic
type AlbumService interface {
	// GetAlbumById returns an album using its ID
	GetAlbumById(id string) (*domain.Album, error)
	// GetAlbumByName returns an album using its name
	GetAlbumByName(name string) (*domain.Album, error)
	// ListAlbums returns all albums
	ListAlbums(fromId string, limit int) ([]*domain.Album, error)
	// AlbumCount returns a count of all albums
	AlbumCount() (int64, error)
	// CreateAlbum creates an album
	CreateAlbum(name string) (*domain.Album, error)
	// CreateSmartAlbum creates a rule-based smart album
	CreateSmartAlbum(name string, rules domain.SmartAlbumRules) (*domain.Album, error)
	// UpdateSmartAlbum replaces a smart album's rules
	UpdateSmartAlbum(id string, rules domain.SmartAlbumRules) (*domain.Album, error)
	// UpdateAlbum updates an album
	UpdateAlbum(id string, updates map[string]any) (*domain.Album, error)
	// DeleteAlbum deletes an album
	DeleteAlbum(id string, deletePhotos bool) error
		// SearchAlbums searches albums by query
	SearchAlbums(query string, limit int) ([]*domain.Album, error)
}
