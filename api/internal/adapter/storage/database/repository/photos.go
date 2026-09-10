package repository

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PhotoRepository struct {
	dbEnv *db.Env
	// workspaceID scopes every tenant-facing query to one workspace (issue
	// #74). Empty means unscoped, which is reserved for background/AI paths
	// that legitimately span workspaces and for tests. Handlers obtain a
	// scoped repository via WithWorkspace.
	workspaceID string
	// scoped records that this repository was explicitly bound to a workspace
	// via WithWorkspace. Once scoped, an empty workspace fails closed (matches
	// nothing) rather than silently widening the query.
	scoped bool
}

func NewPhotoRepository(dbEnv *db.Env) *PhotoRepository {
	return &PhotoRepository{
		dbEnv: dbEnv,
	}
}

// WithWorkspace returns a repository whose tenant-facing queries are scoped to
// workspaceID. Fail-closed: an empty workspace yields a repository that
// matches nothing, so a missing workspace can never widen a query.
func (pr *PhotoRepository) WithWorkspace(workspaceID string) port.PhotoRepository {
	return &PhotoRepository{dbEnv: pr.dbEnv, workspaceID: workspaceID, scoped: true}
}

// scope applies the workspace filter to a photos query. Unscoped repositories
// (background/AI paths) pass through untouched.
func (pr *PhotoRepository) scope(tx *gorm.DB) *gorm.DB {
	if !pr.scoped {
		return tx
	}
	if pr.workspaceID == "" {
		return tx.Where("1 = 0") // fail closed
	}
	return tx.Where("photos.workspace_id = ?", pr.workspaceID)
}

// scopeTags applies the workspace filter to a photo_tags query. photo_tags has
// no workspace column by design (join-through-photos), so membership is
// resolved via a subquery on photos.
func (pr *PhotoRepository) scopeTags(tx *gorm.DB) *gorm.DB {
	if !pr.scoped {
		return tx
	}
	if pr.workspaceID == "" {
		return tx.Where("1 = 0") // fail closed
	}
	return tx.Where("photo_tags.photo_id IN (SELECT id FROM photobox.photos WHERE workspace_id = ?)", pr.workspaceID)
}

// failClosed returns a condition that matches nothing, used when a scoped
// operation must not proceed.
func (pr *PhotoRepository) failClosed(tx *gorm.DB) *gorm.DB {
	if pr.workspaceID == "" {
		return tx
	}
	return tx.Where("1 = 0")
}

// newThumbCap returns a fresh random capability for a photo's thumbnail URL
// (issue #74 D1). Random 128-bit UUID, never derived from the photo ID.
func newThumbCap() string {
	return uuid.NewString()
}

func (pr *PhotoRepository) GetPhotoById(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	var photo domain.Photo
	photo.ID = photoId
	result := pr.scope(pr.dbEnv.Db)
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

// GetPhotoByThumbCap resolves a photo by its capability (issue #74 D1). The
// cap is a random UUID, never derivable from the photo ID. Not-found maps to
// ErrDataNotFound so the unauthenticated capability route returns 404 (no
// existence oracle) for both unknown and wrong caps.
func (pr *PhotoRepository) GetPhotoByThumbCap(thumbCap string) (*domain.Photo, error) {
	var photo domain.Photo
	result := pr.dbEnv.Db.Omit("thumbnail").First(&photo, "thumb_cap = ?", thumbCap)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

// ListMemories returns photos taken on the given month/day in previous years
// ("On This Day"), ordered by year descending. Excludes the current year.
func (pr *PhotoRepository) ListMemories(month, day int, limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).
		Where("deleted_at IS NULL").
		Where("hidden = ?", false).
		Where("EXTRACT(MONTH FROM to_timestamp(created_epoch / 1000.0)) = ?", month).
		Where("EXTRACT(DAY FROM to_timestamp(created_epoch / 1000.0)) = ?", day).
		Where("EXTRACT(YEAR FROM to_timestamp(created_epoch / 1000.0)) < ?", time.Now().Year()).
		Order("created_epoch DESC").
		Omit("thumbnail").
		Limit(limit).
		Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) ListAllPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).Where("deleted_at IS NULL").Where("hidden = ?", false).Order("created_epoch ASC").Omit("thumbnail")
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

// SearchPhotosWithFilters returns photos matching the combined search
// filters. Unused filters are ignored so any combination works.
func (pr *PhotoRepository) SearchPhotosWithFilters(filters domain.PhotoSearchFilters, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).
		Where("deleted_at IS NULL").
		Where("hidden = ?", false).
		Order("created_epoch ASC").
		Omit("thumbnail")

	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("created_epoch > ?", fromPhoto.CreatedEpoch)
	}

	if filters.StartDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05", filters.StartDate); err == nil {
			result = result.Where("created_epoch >= ?", t.UnixMilli())
		}
	}
	if filters.EndDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05", filters.EndDate); err == nil {
			result = result.Where("created_epoch <= ?", t.UnixMilli())
		}
	}
	if filters.MediaType != "" {
		result = result.Where("media_type = ?", filters.MediaType)
	}
	if filters.Camera != "" {
		// Match EXIF Make/Model inside the metadata JSONB payload. JSON text
		// may include whitespace after colons, so match the key and a
		// whitespace-tolerant value pattern.
		like := "%" + filters.Camera + "%"
		result = result.Where(
			"(metadata::text ILIKE ? OR metadata::text ILIKE ?)",
			quoteJSONPattern("Make", like), quoteJSONPattern("Model", like),
		)
	}
	if filters.HasGPS {
		result = result.Where("latitude IS NOT NULL AND longitude IS NOT NULL AND latitude <> 0 AND longitude <> 0")
	}
	if filters.Orientation != "" {
		switch filters.Orientation {
		case "landscape":
			result = result.Where("width > height")
		case "portrait":
			result = result.Where("height > width")
		case "square":
			result = result.Where("width = height")
		}
	}
	if filters.Favorite {
		result = result.Where("favorite = ?", true)
	}
	if filters.LowQuality {
		result = result.Where("is_low_quality = ?", true)
	}
	if len(filters.Tags) > 0 {
		result = result.Where("id IN (SELECT photo_id FROM photo_tags WHERE tag IN ? GROUP BY photo_id HAVING COUNT(DISTINCT tag) = ?)", filters.Tags, len(filters.Tags))
	}

	result = result.Limit(limit)
	result = result.Find(&photos)
	return photos, result.Error
}

// quoteJSONPattern builds a LIKE fragment matching a JSON key's string value,
// tolerant of whitespace after the colon. For key "Make" and pattern
// "%Canon%": %"Make"%:%"%Canon%"%
func quoteJSONPattern(key, valuePattern string) string {
	return "%\"" + key + "\"%:%\"" + valuePattern + "\"%"
}

func (pr *PhotoRepository) ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).Where("deleted_at IS NULL AND album_id = ?", albumId).Where("hidden = ?", false).Order("created_epoch ASC").Omit("thumbnail")
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
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).Where("album_id = ? AND deleted_at IS NULL", albumId).Count(&count)
	err := db.HandleError(result)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pr *PhotoRepository) CreatePhotoInfo(photo domain.Photo) error {
	if photo.ThumbCap == "" {
		photo.ThumbCap = newThumbCap()
	}
	result := pr.dbEnv.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		// thumb_cap is immutable once assigned (issue #74 D1): re-indexing an
		// existing photo must not rotate its capability URL. Exclude it from
		// the conflict update set.
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "filesystem_path", "source_path", "album_id", "metadata",
			"created_at", "created_epoch", "year", "month", "updated_at",
			"thumbnail", "thumbnail_path", "favorite", "deleted_at", "blurhash",
			"dominant_color", "file_hash", "file_modified_time", "media_type",
			"duration", "width", "height", "latitude", "longitude", "hidden",
			"live_photo_path", "description", "quality_score", "blur_score",
			"is_low_quality", "trash_path", "edit_params", "workspace_id",
		}),
	}).Create(&photo)
	return result.Error
}

// CreatePhotosInfo creates multiple photo records in a single batch operation
func (pr *PhotoRepository) CreatePhotosInfo(photos []domain.Photo) error {
	if len(photos) == 0 {
		return nil
	}
	for i := range photos {
		if photos[i].ThumbCap == "" {
			photos[i].ThumbCap = newThumbCap()
		}
	}
	result := pr.dbEnv.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "filesystem_path", "source_path", "album_id", "metadata",
			"created_at", "created_epoch", "year", "month", "updated_at",
			"thumbnail", "thumbnail_path", "favorite", "deleted_at", "blurhash",
			"dominant_color", "file_hash", "file_modified_time", "media_type",
			"duration", "width", "height", "latitude", "longitude", "hidden",
			"live_photo_path", "description", "quality_score", "blur_score",
			"is_low_quality", "trash_path", "edit_params", "workspace_id",
		}),
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
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Update("deleted_at", now)
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
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Update("deleted_at", nil)
	if result.Error != nil {
		return nil, result.Error
	}
	return photo, nil
}

// SetTrashPath records where the original file lives while the photo is in
// the trash (used by RestorePhoto to move the file back).
func (pr *PhotoRepository) SetTrashPath(photoId, trashPath string) error {
	return pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Update("trash_path", trashPath).Error
}

// SetEditParams stores the non-destructive edit parameters for a photo
// (nil clears them).
func (pr *PhotoRepository) SetEditParams(photoId string, params []byte) error {
	if params == nil {
		return pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Update("edit_params", nil).Error
	}
	return pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Update("edit_params", string(params)).Error
}

func (pr *PhotoRepository) ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).Where("deleted_at IS NOT NULL").Order("created_epoch ASC").Omit("thumbnail")
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
	result := pr.scope(pr.dbEnv.Db).Where("deleted_at IS NOT NULL").Delete(&domain.Photo{})
	return result.Error
}

// ListExpiredTrashPhotos returns photos whose deleted_at predates the given
// cutoff. Used by the trash retention job to find photos eligible for
// permanent deletion.
func (pr *PhotoRepository) ListExpiredTrashPhotos(cutoff time.Time) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).
		Where("deleted_at IS NOT NULL").
		Where("deleted_at < ?", cutoff).
		Order("deleted_at ASC").
		Find(&photos)
	return photos, result.Error
}

// PermanentlyDeletePhotos hard-deletes the given photo rows from the database.
func (pr *PhotoRepository) PermanentlyDeletePhotos(photoIds []string) error {
	if len(photoIds) == 0 {
		return nil
	}
	result := pr.scope(pr.dbEnv.Db).Where("id IN ?", photoIds).Delete(&domain.Photo{})
	return result.Error
}

func (pr *PhotoRepository) UpdatePhoto(photo domain.Photo) error {
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photo.ID).Updates(map[string]interface{}{
		"thumbnail_path": photo.ThumbnailPath,
		"blurhash":       photo.Blurhash,
		"dominant_color": photo.DominantColor,
		"updated_at":     time.Now(),
	})
	return result.Error
}

// UpdatePhotoMetadata persists user-editable metadata overrides for a photo.
// Only non-zero fields are updated so partial edits preserve the rest.
func (pr *PhotoRepository) UpdatePhotoMetadata(photoID string, updates domain.Photo) error {
	fields := map[string]interface{}{"updated_at": time.Now()}
	if updates.Description != "" {
		fields["description"] = updates.Description
	}
	if updates.Latitude != 0 || updates.Longitude != 0 {
		fields["latitude"] = updates.Latitude
		fields["longitude"] = updates.Longitude
	}
	if updates.CreatedEpoch != 0 {
		fields["created_epoch"] = updates.CreatedEpoch
		fields["year"] = updates.Year
		fields["month"] = updates.Month
	}
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoID).Updates(fields)
	return result.Error
}

func (pr *PhotoRepository) SetFavorite(photoId string, favorite bool) error {
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Update("favorite", favorite)
	return result.Error
}

// SetFavoriteForPhotos updates the favorite flag for multiple photos in a
// single query.
func (pr *PhotoRepository) SetFavoriteForPhotos(photoIds []string, favorite bool) error {
	if len(photoIds) == 0 {
		return nil
	}
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id IN ?", photoIds).Update("favorite", favorite)
	return result.Error
}

// AssignPhotosToAlbum moves the given photos into an album (single query).
func (pr *PhotoRepository) AssignPhotosToAlbum(photoIds []string, albumId string) error {
	if len(photoIds) == 0 {
		return nil
	}
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id IN ?", photoIds).Update("album_id", albumId)
	return result.Error
}

func (pr *PhotoRepository) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).Where("favorite = ? AND deleted_at IS NULL", true).Where("hidden = ?", false).Order("created_epoch ASC").Omit("thumbnail")
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

// ListLowQualityPhotos returns photos flagged as low quality (blurry,
// badly exposed, or mostly-solid). Ordered by worst quality first.
func (pr *PhotoRepository) ListLowQualityPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).
		Where("is_low_quality = ? AND deleted_at IS NULL", true).
		Where("hidden = ?", false).
		Order("quality_score ASC").
		Omit("thumbnail")
	if fromId != "" {
		fromPhoto, err := pr.GetPhotoById(fromId, false)
		if err != nil {
			return nil, err
		}
		result = result.Where("quality_score < ?", fromPhoto.QualityScore)
	}
	result = result.Limit(limit)
	result = result.Find(&photos)
	return photos, result.Error
}

// ListUnscoredPhotos returns non-video photos without a quality score yet,
// used by the background quality pass to backfill scores for pre-existing
// photos (the index only scores newly-added files).
func (pr *PhotoRepository) ListUnscoredPhotos(limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).
		Where("quality_score = 0 AND deleted_at IS NULL").
		Where("media_type <> ? OR media_type IS NULL", "video").
		Omit("thumbnail").
		Limit(limit).
		Find(&photos)
	return photos, result.Error
}

// UpdatePhotoQuality persists the quality metrics for a photo.
func (pr *PhotoRepository) UpdatePhotoQuality(photoId string, qualityScore int, blurScore float64, isLowQuality bool) error {
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Updates(map[string]interface{}{
		"quality_score":  qualityScore,
		"blur_score":     blurScore,
		"is_low_quality": isLowQuality,
		"updated_at":     time.Now(),
	})
	return result.Error
}

func (pr *PhotoRepository) SearchPhotos(query string, limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	// Full-text match on name (the photos.tags column was dropped in favour of
	// the photo_tags junction table, and idx_photo_search indexes name only).
	result := pr.scope(pr.dbEnv.Db.Model(&[]domain.Photo{})).
		Limit(limit).
		Where("deleted_at IS NULL").
		Where("hidden = ?", false).
		Where("to_tsvector('english', coalesce(name, '')) @@ plainto_tsquery('english', ?)", query).
		Order("created_epoch DESC").
		Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) GetTimeline() ([]domain.TimelineEntry, error) {
	var entries []domain.TimelineEntry
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).
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
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).
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
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).
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
	err := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Select("thumbnail").Where("id = ?", photoId).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result.Thumbnail, nil
}

func (pr *PhotoRepository) GetThumbnailPath(photoId string) (string, error) {
	var result struct{ ThumbnailPath string }
	err := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Select("thumbnail_path").Where("id = ?", photoId).Scan(&result).Error
	if err != nil {
		return "", err
	}
	return result.ThumbnailPath, nil
}

func (pr *PhotoRepository) GetAllTags() ([]string, error) {
	var tags []string
	result := pr.scopeTags(pr.dbEnv.Db.Model(&domain.PhotoTag{})).
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
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).
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
	result := pr.scopeTags(pr.dbEnv.Db.Model(&domain.PhotoTag{})).
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
	pr.scopeTags(pr.dbEnv.Db).Where("photo_id = ? AND (source = '' OR source IS NULL)", photoId).Delete(&domain.PhotoTag{})

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
	pr.scopeTags(pr.dbEnv.Db).Where("photo_id = ? AND source = 'ai'", photoId).Delete(&domain.PhotoTag{})

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
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).
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
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).
		Where("deleted_at IS NULL").
		Where("NOT EXISTS (SELECT 1 FROM photobox.photo_analysis WHERE photobox.photo_analysis.photo_id = photobox.photos.id AND photobox.photo_analysis.status = 'completed')").
		Where("NOT EXISTS (SELECT 1 FROM photobox.photo_analysis WHERE photobox.photo_analysis.photo_id = photobox.photos.id AND photobox.photo_analysis.status = 'failed' AND photobox.photo_analysis.retryable = false)").
		Where("NOT EXISTS (SELECT 1 FROM photobox.photo_analysis WHERE photobox.photo_analysis.photo_id = photobox.photos.id AND photobox.photo_analysis.status = 'failed' AND photobox.photo_analysis.retryable = true AND photobox.photo_analysis.last_attempt_at > NOW() - INTERVAL '1 hour' * LEAST(photobox.photo_analysis.attempts, 24))").
		Limit(limit).
		Omit("thumbnail").
		Find(&photos)
	return photos, result.Error
}

// ListPhotosPendingFaceDetection returns non-deleted photos that have no
// stored face detections yet (face recognition batch processing).
func (pr *PhotoRepository) ListPhotosPendingFaceDetection(limit int) ([]*domain.Photo, error) {
	var photos []*domain.Photo
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).
		Where("deleted_at IS NULL").
		Where("NOT EXISTS (SELECT 1 FROM photobox.face_detections WHERE photobox.face_detections.photo_id = photobox.photos.id)").
		Limit(limit).
		Omit("thumbnail").
		Find(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) SavePhotoAnalysis(analysis domain.PhotoAnalysis) error {
	// Use explicit schema-qualified table to bypass GORM's TablePrefix quoting issue
	return pr.dbEnv.Db.Table("photobox.photo_analysis").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&analysis).Error
}

func (pr *PhotoRepository) GetPhotoAnalysis(photoId string) (*domain.PhotoAnalysis, error) {
	var analysis domain.PhotoAnalysis
	result := pr.dbEnv.Db.Table("photobox.photo_analysis").Where("photo_id = ?", photoId).First(&analysis)
	if result.Error != nil {
		return nil, db.HandleError(result)
	}
	return &analysis, nil
}

func (pr *PhotoRepository) GetDuplicatePhotos() ([]*domain.Photo, error) {
	var photos []*domain.Photo
	sql := `
		SELECT p.* FROM photobox.photos p
		INNER JOIN (
			SELECT file_hash FROM photobox.photos
			WHERE deleted_at IS NULL AND hidden = false AND file_hash <> ''
			{wsFilter}
			GROUP BY file_hash
			HAVING COUNT(*) > 1
		) dup ON p.file_hash = dup.file_hash
		WHERE p.deleted_at IS NULL AND p.hidden = false
	`
	args := []interface{}{}
	if pr.scoped {
		// Scope BOTH the outer scan and the inner hash grouping to the
		// caller's workspace (issue #74). Scoping only the outer query would
		// still count another tenant's identical hash as a duplicate group and
		// falsely report the tenant as having duplicates.
		sql = strings.Replace(sql, "{wsFilter}", "AND workspace_id = ?", 1)
		args = append(args, pr.workspaceID)
		sql += " AND p.workspace_id = ?"
		args = append(args, pr.workspaceID)
	} else {
		sql = strings.Replace(sql, "{wsFilter}", "", 1)
	}
	sql += " ORDER BY p.file_hash, p.created_epoch ASC"
	result := pr.dbEnv.Db.Raw(sql, args...).Scan(&photos)
	return photos, result.Error
}

func (pr *PhotoRepository) GetPhotoIndexCache() (map[string]struct {
	FileHash         string
	FileModifiedTime int64
}, error) {
	type cacheEntry struct {
		ID               string
		FileHash         string
		FileModifiedTime int64
	}
	var entries []cacheEntry
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Select("id", "file_hash", "file_modified_time").Where("deleted_at IS NULL").Find(&entries)
	if result.Error != nil {
		return nil, result.Error
	}
	cache := make(map[string]struct {
		FileHash         string
		FileModifiedTime int64
	}, len(entries))
	for _, e := range entries {
		cache[e.ID] = struct {
			FileHash         string
			FileModifiedTime int64
		}{FileHash: e.FileHash, FileModifiedTime: e.FileModifiedTime}
	}
	return cache, nil
}

func (pr *PhotoRepository) HidePhoto(photoId string) error {
	result := pr.scope(pr.dbEnv.Db.Model(&domain.Photo{})).Where("id = ?", photoId).Update("hidden", true)
	return result.Error
}
