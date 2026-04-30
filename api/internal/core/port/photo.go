package port

import (
	"io"

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
	// SoftDeletePhoto marks a photo as deleted and returns its filesystem path
	SoftDeletePhoto(photoId string) (*domain.Photo, error)
	// RestorePhoto restores a photo from trash
	RestorePhoto(photoId string) (*domain.Photo, error)
	// ListTrashPhotos returns all photos in trash
	ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// EmptyTrash permanently deletes all trashed photos
	EmptyTrash() error
	// UpdatePhoto updates a photo record
	UpdatePhoto(photo domain.Photo) error
	// ListFavoritePhotos returns only favorited photos
	ListFavoritePhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// SearchPhotos searches photos by query
	SearchPhotos(query string, limit int) ([]*domain.Photo, error)
	// GetTimeline returns photo counts grouped by year/month
	GetTimeline() ([]domain.TimelineEntry, error)
	// GetPhotosWithGeodata returns photos that have GPS coordinates
	GetPhotosWithGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error)
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
	// DeletePhoto soft-deletes a photo (moves to trash)
	DeletePhoto(photoId string) error
	// RestorePhoto restores a photo from trash
	RestorePhoto(photoId string) error
	// ListTrashPhotos returns paginated trash photos
	ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// EmptyTrash permanently deletes all trashed photos
	EmptyTrash() error
	// SetFavorite toggles favorite status
	SetFavorite(photoId string, favorite bool) (*domain.Photo, error)
	// ListFavoritePhotos returns favorited photos
	ListFavoritePhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// Search searches photos by query
	Search(query string, limit int) ([]*domain.Photo, error)
	// GetTimeline returns photo timeline data
	GetTimeline() ([]domain.TimelineEntry, error)
	// GetGeodata returns photos with GPS coordinates
	GetGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error)
	// RotatePhoto rotates a photo
	RotatePhoto(photoId string, direction string) (*domain.Photo, error)
	// DownloadPhotos streams a zip of photos
	DownloadPhotos(photoIds []string, writer io.Writer) error
}
