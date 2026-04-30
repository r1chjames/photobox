package repository

import (
	b64 "encoding/base64"
	"time"

	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
)

type PhotoRepository struct {
	dbEnv *db.Env
}

func NewPhotoRepository(dbEnv *db.Env) *PhotoRepository {
	return &PhotoRepository{
		dbEnv,
	}
}

func (pr *PhotoRepository) GetPhotoById(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	var photo domain.Photo
	photo.ID = photoId
	result := pr.dbEnv.Db
	if !includeThumbnail {
		result = result.Omit("thumbnail")
	}
	result = result.First(&photo)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (pr *PhotoRepository) ListAllPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("deleted_at IS NULL").Order("created_epoch ASC").Omit("thumbnail")
	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}
	if startDate != "" {
		result = result.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		result = result.Where("created_at <= ?", endDate)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, nil
}

func (pr *PhotoRepository) ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("deleted_at IS NULL AND album_id = ?", albumId).Order("created_epoch ASC").Omit("thumbnail")
	if fromId != "" {
		fromEpoch, _ := b64.StdEncoding.DecodeString(fromId)
		result = result.Where("created_epoch > ?", fromEpoch)
	}
	if startDate != "" {
		result = result.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		result = result.Where("created_at <= ?", endDate)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, nil
}

func (pr *PhotoRepository) GetPhotosInAlbumCount(albumId string) (int64, error) {
	var count int64
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("album_id = ? AND deleted_at IS NULL", albumId).Count(&count)
	err := db.HandleError(result)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pr *PhotoRepository) CreatePhotoInfo(photo domain.Photo) error {
	result := pr.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&photo)
	return result.Error
}

// CreatePhotosInfo creates multiple photo records in a single batch operation
func (pr *PhotoRepository) CreatePhotosInfo(photos []domain.Photo) error {
	if len(photos) == 0 {
		return nil
	}
	result := pr.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&photos)
	return result.Error
}

func (pr *PhotoRepository) SoftDeletePhoto(photoId string) (*domain.Photo, error) {
	photo, err := pr.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	photo.DeletedAt = &now
	result := pr.dbEnv.Db.Save(photo)
	if result.Error != nil {
		return nil, result.Error
	}
	return photo, nil
}

func (pr *PhotoRepository) RestorePhoto(photoId string) (*domain.Photo, error) {
	photo, err := pr.GetPhotoById(photoId, false)
	if err != nil {
		return nil, err
	}
	photo.DeletedAt = nil
	result := pr.dbEnv.Db.Save(photo)
	if result.Error != nil {
		return nil, result.Error
	}
	return photo, nil
}

func (pr *PhotoRepository) ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("deleted_at IS NOT NULL").Order("created_epoch ASC").Omit("thumbnail")
	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, nil
}

func (pr *PhotoRepository) EmptyTrash() error {
	result := pr.dbEnv.Db.Where("deleted_at IS NOT NULL").Delete(&domain.Photo{})
	return result.Error
}

func (pr *PhotoRepository) UpdatePhoto(photo domain.Photo) error {
	result := pr.dbEnv.Db.Save(&photo)
	return result.Error
}

func (pr *PhotoRepository) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("favorite = ? AND deleted_at IS NULL", true).Order("created_epoch ASC").Omit("thumbnail")
	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}
	if startDate != "" {
		result = result.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		result = result.Where("created_at <= ?", endDate)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, nil
}

func (pr *PhotoRepository) SearchPhotos(query string, limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).
		Limit(limit).
		Where("deleted_at IS NULL").
		Where("to_tsvector('english', coalesce(name, '')) @@ plainto_tsquery('english', ?)", query).
		Order("created_epoch DESC").
		Find(&photos)
	if result.RowsAffected == 0 {
		return nil, domain.ErrDataNotFound
	}
	return photos, result.Error
}

func (pr *PhotoRepository) GetTimeline() ([]domain.TimelineEntry, error) {
	var entries []domain.TimelineEntry
	result := pr.dbEnv.Db.Raw(`
		SELECT 
			EXTRACT(YEAR FROM to_timestamp(created_epoch / 1000))::int AS year,
			EXTRACT(MONTH FROM to_timestamp(created_epoch / 1000))::int AS month,
			COUNT(*) AS count
		FROM photobox.photos
		WHERE deleted_at IS NULL
		GROUP BY year, month
		ORDER BY year DESC, month DESC
	`).Scan(&entries)
	if result.Error != nil {
		return nil, result.Error
	}
	return entries, nil
}

func (pr *PhotoRepository) GetPhotosWithGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error) {
	var photos []domain.Photo
	result := pr.dbEnv.Db.Model(&domain.Photo{}).
		Where("deleted_at IS NULL").
		Where("metadata::jsonb -> 'Exif' ->> 'GPSLatitude' IS NOT NULL").
		Omit("thumbnail").
		Find(&photos)
	if result.Error != nil {
		return nil, result.Error
	}

	var results []domain.PhotoGeoData
	for _, p := range photos {
		results = append(results, domain.PhotoGeoData{
			ID:        p.ID,
			Thumbnail: p.SourcePath,
			DateTaken: p.CreatedAt.Format(time.RFC3339),
		})
	}
	return results, nil
}

func (pr *PhotoRepository) GetPhotoThumbnails(photoIds []string) (map[string][]byte, error) {
	var photos []domain.Photo
	result := pr.dbEnv.Db.Model(&domain.Photo{}).
		Select("id", "thumbnail").
		Where("id IN ?", photoIds).
		Find(&photos)
	if result.Error != nil {
		return nil, result.Error
	}

	thumbnails := make(map[string][]byte, len(photos))
	for _, p := range photos {
		thumbnails[p.ID] = p.Thumbnail
	}
	return thumbnails, nil
}
