package port

import (
	"context"
	"io"
	"time"

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
	// SetTrashPath records the trash file location for a soft-deleted photo
	SetTrashPath(photoId, trashPath string) error
	// SetEditParams stores the non-destructive edit parameters for a photo
	SetEditParams(photoId string, params []byte) error
	// ListTrashPhotos returns all photos in trash
	ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// ListExpiredTrashPhotos returns trashed photos deleted before the given cutoff
	ListExpiredTrashPhotos(cutoff time.Time) ([]*domain.Photo, error)
	// PermanentlyDeletePhotos hard-deletes the given photo rows
	PermanentlyDeletePhotos(photoIds []string) error
	// EmptyTrash permanently deletes all trashed photos
	EmptyTrash() error
	// UpdatePhoto updates a photo record
	UpdatePhoto(photo domain.Photo) error
	// UpdatePhotoMetadata persists user-editable metadata overrides
	UpdatePhotoMetadata(photoId string, updates domain.Photo) error
	// SetFavorite updates only the favorite field of a photo
	SetFavorite(photoId string, favorite bool) error
	// SetFavoriteForPhotos updates the favorite flag for multiple photos
	SetFavoriteForPhotos(photoIds []string, favorite bool) error
	// AssignPhotosToAlbum moves multiple photos into an album
	AssignPhotosToAlbum(photoIds []string, albumId string) error
	// ListFavoritePhotos returns only favorited photos
	ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error)
	// ListLowQualityPhotos returns photos flagged as low quality
	ListLowQualityPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// ListMemories returns photos taken on a given month/day in prior years
	ListMemories(month, day int, limit int) ([]*domain.Photo, error)
	// SearchPhotosWithFilters returns photos matching the combined filters
	SearchPhotosWithFilters(filters domain.PhotoSearchFilters, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// ListUnscoredPhotos returns photos without a quality score yet
	ListUnscoredPhotos(limit int) ([]*domain.Photo, error)
	// UpdatePhotoQuality persists quality metrics for a photo
	UpdatePhotoQuality(photoId string, qualityScore int, blurScore float64, isLowQuality bool) error
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
	// GetPhotoTags returns all tags for a specific photo
	GetPhotoTags(photoId string) ([]string, error)
	// ListPhotosByTags returns photos matching all specified tags
	ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// UpdatePhotoTags updates the tags for a single photo (preserves AI tags)
	UpdatePhotoTags(photoId string, tags string) error
	// AddAITags adds AI-generated tags for a photo
	AddAITags(photoId string, tags []string) error
	// ListPhotosWithoutAITags returns photos that have no AI-generated tags
	ListPhotosWithoutAITags(limit int) ([]*domain.Photo, error)
	// ListPhotosPendingAnalysis returns photos that need AI analysis (not completed, not failed-with-backoff)
	ListPhotosPendingAnalysis(limit int) ([]*domain.Photo, error)
	// SavePhotoAnalysis saves or updates the AI analysis result for a photo
	SavePhotoAnalysis(analysis domain.PhotoAnalysis) error
	// GetPhotoAnalysis retrieves the AI analysis for a photo
	GetPhotoAnalysis(photoId string) (*domain.PhotoAnalysis, error)
	// GetDuplicatePhotos returns photos that have duplicate file hashes
	GetDuplicatePhotos() ([]*domain.Photo, error)
	// GetPhotoIndexCache returns a map of photo ID to file hash and modified time for skip-unchanged optimization
	GetPhotoIndexCache() (map[string]struct{FileHash string; FileModifiedTime int64}, error)
	// HidePhoto marks a photo as hidden (corrupted/unusable)
	HidePhoto(photoId string) error
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
	// UploadPhoto writes an uploaded photo and indexes it
	UploadPhoto(upload domain.PhotoUpload) (*domain.Photo, error)
	//SavePhotos saves photos to database
	SavePhotos(photos []domain.PhotoFile) error
	// PerformPhotoIndex initiates an index of image files on filesystem
	PerformPhotoIndex(ctx context.Context)
	// DeletePhoto soft-deletes a photo (moves to trash)
	DeletePhoto(photoId string) error
	// RestorePhoto restores a photo from trash
	RestorePhoto(photoId string) error
	// ListTrashPhotos returns paginated trash photos
	ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// EmptyTrash permanently deletes all trashed photos
	EmptyTrash() error
	// PurgeExpiredTrash permanently deletes trashed photos older than the cutoff
	PurgeExpiredTrash(cutoff time.Time) (int, error)
	// SetFavorite toggles favorite status
	SetFavorite(photoId string, favorite bool) (*domain.Photo, error)
	// UpdatePhotoMetadata persists user-editable metadata overrides
	UpdatePhotoMetadata(photoId string, description string, latitude *float64, longitude *float64, dateTaken *string) (*domain.Photo, error)
	// ReverseGeocode resolves a human-readable location for a photo's GPS coords
	ReverseGeocode(ctx context.Context, photoId string) (string, error)
	// BatchSetFavorite sets favorite for multiple photos
	BatchSetFavorite(photoIds []string, favorite bool) error
	// BatchAddToAlbum assigns multiple photos to an album
	BatchAddToAlbum(photoIds []string, albumId string) error
	// BatchDeletePhotos soft-deletes (moves to trash) multiple photos
	BatchDeletePhotos(photoIds []string) error
	// ListFavoritePhotos returns favorited photos
	ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error)
	// ListLowQualityPhotos returns photos flagged as low quality
	ListLowQualityPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// ListMemories returns "On This Day" photos grouped by year
	ListMemories(month, day int, maxPerYear int) ([]domain.MemoryGroup, error)
	// SearchPhotosWithFilters returns photos matching the combined filters
	SearchPhotosWithFilters(filters domain.PhotoSearchFilters, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error)
	// ScorePhotoQuality runs the quality analyzer on a photo and stores results
	ScorePhotoQuality(photoId string) (bool, error)
	// ScoreAllPhotoQuality backfills quality scores for all unscored photos
	ScoreAllPhotoQuality(batchSize int) (int, error)
	// Search searches photos by query
	Search(query string, limit int) ([]*domain.Photo, error)
	// GetTimeline returns photo timeline data
	GetTimeline() ([]domain.TimelineEntry, error)
	// GetGeodata returns photos with GPS coordinates
	GetGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error)
	// RotatePhoto rotates a photo
	RotatePhoto(photoId string, direction string) (*domain.Photo, error)
	// EditPhoto applies a non-destructive edit (rotate/crop/adjust)
	EditPhoto(photoId string, params domain.EditParams) (*domain.Photo, error)
	// ClearEdits removes any edited copy and stored edit params
	ClearEdits(photoId string) (*domain.Photo, error)
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
	// TriggerAIAnalysis starts AI analysis manually (returns immediately, runs in background)
	TriggerAIAnalysis() error
	// RegenerateThumbnails regenerates thumbnails for all photos (skips photos with existing thumbnails)
	RegenerateThumbnails(ctx context.Context)
}
