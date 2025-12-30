package port

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

//go:generate mockgen -source=photo.go -destination=mock/photo.go -package=mock

// PhotoRepository is an interface for interacting with photo-related data
type PhotoRepository interface {
	// GetPhotoById inserts a new user into the database
	GetPhotoById(photoId string, includeThumbnail bool) (*domain.Photo, error)
	// ListAllPhotos selects a list of users with pagination
	ListAllPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// ListAllPhotosInAlbum selects a user by id
	ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// GetPhotosInAlbumCount selects a user by id
	GetPhotosInAlbumCount(albumId string) (int64, error)
	// CreatePhotoInfo selects a user by email
	CreatePhotoInfo(photo domain.Photo) error
	// CreatePhotosInfo creates multiple photo records in a single batch operation
	CreatePhotosInfo(photos []domain.Photo) error
}

// PhotoService is an interface for interacting with photo-related business logic
type PhotoService interface {
	//GetPhoto returns a photo
	GetPhoto(photoId string, includeThumbnail bool) (*domain.Photo, error)
	// ListPhotos registers a new user
	ListPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// ListPhotosInAlbum registers a new user
	ListPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	//PhotoCount returns a count of photos in an album
	PhotoCount(albumId string) (int64, error)
	//PhotoBinary returns the binary photo from disk
	PhotoBinary(photoId string) (string, error)
	//PhotoThumbnail returns the binary photo from disk
	PhotoThumbnail(photoId string) ([]byte, error)
	//SavePhoto saves photo to database
	SavePhoto(photo domain.PhotoFile) error
	//SavePhotos saves photos to database
	SavePhotos(photos []domain.PhotoFile) error
	// PerformPhotoIndex initiates an index of image files on filesystem
	PerformPhotoIndex()
}
