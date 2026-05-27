package repository

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
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

func (pr *PhotoRepository) ListAllPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("deleted_at IS NULL").Where("hidden = ?", false).Order("created_epoch ASC").Omit("thumbnail")
	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}
	var startEpoch, endEpoch int64
	if startDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05", startDate)
		if err == nil {
			startEpoch = t.UnixMilli()
		}
	}
	if endDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05", endDate)
		if err == nil {
			endEpoch = t.UnixMilli()
		}
	}
	if startDate != "" {
		result = result.Where("created_epoch >= ?", startEpoch)
	}
	if endDate != "" {
		result = result.Where("created_epoch <= ?", endEpoch)
	}
	if mediaType != "" {
		result = result.Where("media_type = ?", mediaType)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("deleted_at IS NULL AND album_id = ?", albumId).Where("hidden = ?", false).Order("created_epoch ASC").Omit("thumbnail")
	if fromId != "" {
		fromEpoch, _ := b64.StdEncoding.DecodeString(fromId)
		result = result.Where("created_epoch > ?", fromEpoch)
	}
	var startEpoch, endEpoch int64
	if startDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05", startDate)
		if err == nil {
			startEpoch = t.UnixMilli()
		}
	}
	if endDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05", endDate)
		if err == nil {
			endEpoch = t.UnixMilli()
		}
	}
	if startDate != "" {
		result = result.Where("created_epoch >= ?", startEpoch)
	}
	if endDate != "" {
		result = result.Where("created_epoch <= ?", endEpoch)
	}
	if mediaType != "" {
		result = result.Where("media_type = ?", mediaType)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	return photos, result.Error
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
	result := pr.dbEnv.Db.Model(&domain.Photo{}).Where("id = ?", photoId).Update("deleted_at", now)
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
	result := pr.dbEnv.Db.Model(&domain.Photo{}).Where("id = ?", photoId).Update("deleted_at", nil)
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
	return photos, result.Error
}

func (pr *PhotoRepository) EmptyTrash() error {
	result := pr.dbEnv.Db.Where("deleted_at IS NOT NULL").Delete(&domain.Photo{})
	return result.Error
}

func (pr *PhotoRepository) UpdatePhoto(photo domain.Photo) error {
	result := pr.dbEnv.Db.Model(&domain.Photo{}).Where("id = ?", photo.ID).Updates(map[string]interface{}{
		"thumbnail_path": photo.ThumbnailPath,
		"blurhash":       photo.Blurhash,
		"updated_at":     time.Now(),
	})
	return result.Error
}

func (pr *PhotoRepository) SetFavorite(photoId string, favorite bool) error {
	result := pr.dbEnv.Db.Model(&domain.Photo{}).Where("id = ?", photoId).Update("favorite", favorite)
	return result.Error
}

func (pr *PhotoRepository) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).Where("favorite = ? AND deleted_at IS NULL", true).Where("hidden = ?", false).Order("created_epoch ASC").Omit("thumbnail")
	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}
	var startEpoch, endEpoch int64
	if startDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05", startDate)
		if err == nil {
			startEpoch = t.UnixMilli()
		}
	}
	if endDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05", endDate)
		if err == nil {
			endEpoch = t.UnixMilli()
		}
	}
	if startDate != "" {
		result = result.Where("created_epoch >= ?", startEpoch)
	}
	if endDate != "" {
		result = result.Where("created_epoch <= ?", endEpoch)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) SearchPhotos(query string, limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&[]domain.Photo{}).
		Limit(limit).
		Where("deleted_at IS NULL").
		Where("hidden = ?", false).
		Where("to_tsvector('english', coalesce(name, '') || ' ' || coalesce(tags, '')) @@ plainto_tsquery('english', ?)", query).
		Order("created_epoch DESC").
		Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) GetTimeline() ([]domain.TimelineEntry, error) {
	var entries []domain.TimelineEntry
	result := pr.dbEnv.Db.Model(&domain.Photo{}).
		Select("year", "month", "COUNT(*) as count").
		Where("deleted_at IS NULL AND year > 0 AND month > 0").
		Group("year, month").
		Order("year DESC, month DESC").
		Scan(&entries)
	if result.Error != nil {
		return nil, result.Error
	}
	return entries, nil
}

func (pr *PhotoRepository) GetPhotosWithGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error) {
	var photos []domain.Photo
	result := pr.dbEnv.Db.Model(&domain.Photo{}).
		Where("deleted_at IS NULL").
		Where("hidden = ?", false).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL").
		Where("latitude BETWEEN ? AND ?", south, north).
		Where("longitude BETWEEN ? AND ?", west, east).
		Omit("thumbnail").
		Limit(1000).
		Find(&photos)
	if result.Error != nil {
		return nil, result.Error
	}

	var results []domain.PhotoGeoData
	for _, p := range photos {
		results = append(results, domain.PhotoGeoData{
			ID:        p.ID,
			Lat:       p.Latitude,
			Lng:       p.Longitude,
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

func (pr *PhotoRepository) GetThumbnailBytes(photoId string) ([]byte, error) {
	var result struct{ Thumbnail []byte }
	err := pr.dbEnv.Db.Model(&domain.Photo{}).Select("thumbnail").Where("id = ?", photoId).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result.Thumbnail, nil
}

func (pr *PhotoRepository) GetThumbnailPath(photoId string) (string, error) {
	var result struct{ ThumbnailPath string }
	err := pr.dbEnv.Db.Model(&domain.Photo{}).Select("thumbnail_path").Where("id = ?", photoId).Scan(&result).Error
	if err != nil {
		return "", err
	}
	return result.ThumbnailPath, nil
}

func (pr *PhotoRepository) GetAllTags() ([]string, error) {
	var tags []string
	result := pr.dbEnv.Db.Model(&domain.PhotoTag{}).
		Distinct("tag").
		Where("tag <> ''").
		Order("tag ASC").
		Pluck("tag", &tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (pr *PhotoRepository) ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo

	// Build subquery for photos matching ALL tags
	result := pr.dbEnv.Db.Model(&domain.Photo{}).
		Where("deleted_at IS NULL").
		Where("hidden = ?", false).
		Order("created_epoch ASC").
		Omit("thumbnail")

	// Join with photo_tags and filter
	for i, tag := range tags {
		alias := fmt.Sprintf("pt%d", i)
		result = result.Joins(fmt.Sprintf("JOIN photobox.photo_tags %s ON %s.photo_id = photobox.photos.id AND %s.tag = ?", alias, alias, alias), tag)
	}

	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) GetPhotoTags(photoId string) ([]string, error) {
	var tags []string
	result := pr.dbEnv.Db.Model(&domain.PhotoTag{}).
		Where("photo_id = ?", photoId).
		Where("tag <> ''").
		Order("tag ASC").
		Pluck("tag", &tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (pr *PhotoRepository) UpdatePhotoTags(photoId string, tags string) error {
	// Sync junction table: remove only non-AI tags
	pr.dbEnv.Db.Where("photo_id = ? AND (source = '' OR source IS NULL)", photoId).Delete(&domain.PhotoTag{})

	tagList := []domain.PhotoTag{}
	for _, tag := range strings.Split(tags, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tagList = append(tagList, domain.PhotoTag{PhotoID: photoId, Tag: tag, Source: ""})
		}
	}
	if len(tagList) > 0 {
		return pr.dbEnv.Db.Create(&tagList).Error
	}
	return nil
}

func (pr *PhotoRepository) AddAITags(photoId string, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	// Remove existing AI tags for this photo
	pr.dbEnv.Db.Where("photo_id = ? AND source = 'ai'", photoId).Delete(&domain.PhotoTag{})

	tagList := []domain.PhotoTag{}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tagList = append(tagList, domain.PhotoTag{PhotoID: photoId, Tag: tag, Source: "ai"})
		}
	}
	if len(tagList) == 0 {
		return nil
	}

	return pr.dbEnv.Db.Create(&tagList).Error
}

func (pr *PhotoRepository) ListPhotosWithoutAITags(limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&domain.Photo{}).
		Where("deleted_at IS NULL").
		Where("hidden = ?", false).
		Where("NOT EXISTS (SELECT 1 FROM photobox.photo_tags WHERE photobox.photo_tags.photo_id = photobox.photos.id AND photobox.photo_tags.source = 'ai')").
		Limit(limit).
		Omit("thumbnail").
		Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) ListPhotosPendingAnalysis(limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Model(&domain.Photo{}).
		Where("deleted_at IS NULL").
		Where("NOT EXISTS (SELECT 1 FROM photobox.photo_analyses WHERE photobox.photo_analyses.photo_id = photobox.photos.id AND photobox.photo_analyses.status = 'completed')").
		Where("NOT EXISTS (SELECT 1 FROM photobox.photo_analyses WHERE photobox.photo_analyses.photo_id = photobox.photos.id AND photobox.photo_analyses.status = 'failed' AND photobox.photo_analyses.last_attempt_at > NOW() - INTERVAL '1 hour' * LEAST(photobox.photo_analyses.attempts, 24))").
		Limit(limit).
		Omit("thumbnail").
		Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) SavePhotoAnalysis(analysis domain.PhotoAnalysis) error {
	return pr.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&analysis).Error
}

func (pr *PhotoRepository) GetPhotoAnalysis(photoId string) (*domain.PhotoAnalysis, error) {
	var analysis domain.PhotoAnalysis
	result := pr.dbEnv.Db.Where("photo_id = ?", photoId).First(&analysis)
	if result.Error != nil {
		return nil, db.HandleError(result)
	}
	return &analysis, nil
}

func (pr *PhotoRepository) GetDuplicatePhotos() ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.dbEnv.Db.Raw(`
		SELECT p.* FROM photobox.photos p
		INNER JOIN (
			SELECT file_hash FROM photobox.photos
			WHERE deleted_at IS NULL AND hidden = false AND file_hash <> ''
			GROUP BY file_hash
			HAVING COUNT(*) > 1
		) dup ON p.file_hash = dup.file_hash
		WHERE p.deleted_at IS NULL AND p.hidden = false
		ORDER BY p.file_hash, p.created_epoch ASC
	`).Scan(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) GetPhotoIndexCache() (map[string]struct{ FileHash string; FileModifiedTime int64 }, error) {
	type cacheEntry struct {
		ID               string
		FileHash         string
		FileModifiedTime int64
	}
	var entries []cacheEntry
	result := pr.dbEnv.Db.Model(&domain.Photo{}).Select("id", "file_hash", "file_modified_time").Where("deleted_at IS NULL").Find(&entries)
	if result.Error != nil {
		return nil, result.Error
	}
	cache := make(map[string]struct{ FileHash string; FileModifiedTime int64 }, len(entries))
	for _, e := range entries {
		cache[e.ID] = struct{ FileHash string; FileModifiedTime int64 }{FileHash: e.FileHash, FileModifiedTime: e.FileModifiedTime}
	}
	return cache, nil
}

func (pr *PhotoRepository) HidePhoto(photoId string) error {
	result := pr.dbEnv.Db.Model(&domain.Photo{}).Where("id = ?", photoId).Update("hidden", true)
	return result.Error
}
