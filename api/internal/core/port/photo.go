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
	ListAllPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error)
	// ListAllPhotosInAlbum selects a user by id
	ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error)
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
	// SetFavorite updates only the favorite field of a photo
	SetFavorite(photoId string, favorite bool) error
	// ListFavoritePhotos returns only favorited photos
	ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error)
	// SearchPhotos searches photos by query
	SearchPhotos(query string, limit int) ([]*domain.Photo, error)
	// GetTimeline returns photo counts grouped by year/month
	GetTimeline() ([]domain.TimelineEntry, error)
	// GetPhotosWithGeodata returns photos that have GPS coordinates
	GetPhotosWithGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error)
	// GetPhotoThumbnails returns thumbnails for multiple photos
	GetPhotoThumbnails(photoIds []string) (map[string][]byte, error)
	// GetThumbnailBytes returns thumbnail bytes for a single photo
	GetThumbnailBytes(photoId string) ([]byte, error)
	// GetThumbnailPath returns the filesystem path for a photo's thumbnail
	GetThumbnailPath(photoId string) (string, error)
	// GetAllTags returns all distinct tags across photos
	GetAllTags() ([]string, error)
	// ListPhotosByTags returns photos matching all specified tags
	ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// UpdatePhotoTags updates the tags for a single photo (preserves AI tags)
	UpdatePhotoTags(photoId string, tags string) error
	// AddAITags adds AI-generated tags for a photo
	AddAITags(photoId string, tags []string) error
	// ListPhotosWithoutAITags returns photos that have no AI-generated tags
	ListPhotosWithoutAITags(limit int) ([]*domain.Photo, error)
	// GetDuplicatePhotos returns photos that have duplicate file hashes
	GetDuplicatePhotos() ([]*domain.Photo, error)
	// GetPhotoIndexCache returns a map of photo ID to file hash and modified time for skip-unchanged optimization
	GetPhotoIndexCache() (map[string]struct{FileHash string; FileModifiedTime int64}, error)
}

// PhotoService is an interface for interacting with photo-related business logic
type PhotoService interface {
	//GetPhoto returns a photo
	GetPhoto(photoId string, includeThumbnail bool) (*domain.Photo, error)
	// ListPhotos registers a new user
	ListPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error)
	// ListPhotosInAlbum registers a new user
	ListPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error)
	//PhotoCount returns a count of photos in an album
	PhotoCount(albumId string) (int64, error)
	//PhotoBinary returns the binary photo from disk
	PhotoBinary(photoId string) (string, error)
	//PhotoThumbnail returns the binary photo from disk
	PhotoThumbnail(photoId string) ([]byte, error)
	//PhotoThumbnailBytes returns thumbnail bytes directly from DB
	PhotoThumbnailBytes(photoId string) ([]byte, error)
	//PhotoThumbnailPath returns the filesystem path for a photo's thumbnail
	PhotoThumbnailPath(photoId string) (string, error)
	//PhotoThumbnails returns thumbnails for multiple photos
	PhotoThumbnails(photoIds []string) (map[string][]byte, error)
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
	ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error)
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
	// GetAllTags returns all distinct tags across photos
	GetAllTags() ([]string, error)
	// ListPhotosByTags returns photos matching all specified tags
	ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// UpdatePhotoTags updates tags for a single photo
	UpdatePhotoTags(photoId string, tags []string) (*domain.Photo, error)
	// BatchUpdatePhotoTags adds/removes/sets tags for multiple photos
	BatchUpdatePhotoTags(photoIds []string, tags []string, operation string) error
	// GetDuplicatePhotos returns photos with duplicate file hashes
	GetDuplicatePhotos() ([]*domain.Photo, error)
	// GenerateThumbnailForPhoto generates a thumbnail on-demand for a photo
	GenerateThumbnailForPhoto(photoId string) (string, error)
	// PhotoThumbnailPathForSize returns the filesystem path for a photo's thumbnail at a given size
	PhotoThumbnailPathForSize(photoId string, size string) (string, error)
	// PhotoThumbnailBytesForSize returns the raw thumbnail bytes for a photo at a given size
	PhotoThumbnailBytesForSize(photoId string, size string) ([]byte, error)
	// AnalyzeExistingPhotos runs AI analysis on photos that haven't been analyzed yet
	AnalyzeExistingPhotos() error
	// RegenerateThumbnails regenerates thumbnails for all photos (skips photos with existing thumbnails)
	RegenerateThumbnails()
}
