package repository

import (
	"time"

	"github.com/google/uuid"
	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FaceRepository persists face detections, people, and clusters.
type FaceRepository struct {
	dbEnv *db.Env
}

func NewFaceRepository(dbEnv *db.Env) *FaceRepository {
	return &FaceRepository{
		dbEnv,
	}
}

// InsertDetectionsForPhoto replaces all detections for a photo with the given
// set (delete-then-insert per photo). Idempotent re-analysis is safe.
func (fr *FaceRepository) InsertDetectionsForPhoto(photoID string, dets []domain.FaceDetection) error {
	return fr.dbEnv.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("photo_id = ?", photoID).Delete(&domain.FaceDetection{}).Error; err != nil {
			return err
		}
		if len(dets) == 0 {
			return nil
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&dets).Error
	})
}

// DeleteDetectionsForPhoto removes all detections belonging to a photo
// (cascade on photo delete).
func (fr *FaceRepository) DeleteDetectionsForPhoto(photoID string) error {
	return fr.dbEnv.Db.Where("photo_id = ?", photoID).Delete(&domain.FaceDetection{}).Error
}

// DeleteAllFaceData wipes every face table (privacy "delete all data" action).
func (fr *FaceRepository) DeleteAllFaceData() (int64, error) {
	var total int64
	err := fr.dbEnv.Db.Transaction(func(tx *gorm.DB) error {
		for _, model := range []any{&domain.FaceCluster{}, &domain.FacePerson{}, &domain.FaceDetection{}} {
			res := tx.Session(&gorm.Session{}).Delete(model, "1 = 1")
			if res.Error != nil {
				return res.Error
			}
			total += res.RowsAffected
		}
		return nil
	})
	return total, err
}

// FacesForPhoto returns detections for a photo, newest first.
func (fr *FaceRepository) FacesForPhoto(photoID string) ([]domain.FaceDetection, error) {
	var dets []domain.FaceDetection
	result := fr.dbEnv.Db.Where("photo_id = ?", photoID).Order("created_at ASC").Find(&dets)
	return dets, result.Error
}

// FacesForPerson returns detections assigned to a person.
func (fr *FaceRepository) FacesForPerson(personID string, limit int) ([]domain.FaceDetection, error) {
	var dets []domain.FaceDetection
	q := fr.dbEnv.Db.Where("person_id = ?", personID).Order("created_at ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	result := q.Find(&dets)
	return dets, result.Error
}

// PersonPhotoIDs returns the distinct photo IDs containing a person's faces.
func (fr *FaceRepository) PersonPhotoIDs(personID string) ([]string, error) {
	var ids []string
	result := fr.dbEnv.Db.Model(&domain.FaceDetection{}).
		Where("person_id = ?", personID).
		Distinct().
		Pluck("photo_id", &ids)
	return ids, result.Error
}

// Persons returns every person with a distinct photo count, ordered by name.
func (fr *FaceRepository) Persons() ([]domain.FacePerson, error) {
	var persons []domain.FacePerson
	result := fr.dbEnv.Db.Model(&domain.FacePerson{}).
		Select("face_persons.*, COUNT(DISTINCT fd.photo_id) AS photo_count").
		Joins("LEFT JOIN face_detections fd ON fd.person_id = face_persons.id").
		Group("face_persons.id").
		Order("face_persons.name ASC").
		Scan(&persons)
	return persons, result.Error
}

// PersonByID returns a single person.
func (fr *FaceRepository) PersonByID(id string) (*domain.FacePerson, error) {
	var person domain.FacePerson
	person.ID = id
	result := fr.dbEnv.Db.First(&person)
	if err := db.HandleError(result); err != nil {
		return nil, err
	}
	return &person, nil
}

// CreatePerson inserts a person and returns it with its ID populated.
func (fr *FaceRepository) CreatePerson(name string) (*domain.FacePerson, error) {
	now := time.Now()
	person := domain.FacePerson{ID: uuid.NewString(), Name: name, CreatedAt: now, UpdatedAt: now}
	result := fr.dbEnv.Db.Create(&person)
	if err := db.HandleError(result); err != nil {
		return nil, err
	}
	return &person, nil
}

// UpdatePerson renames a person.
func (fr *FaceRepository) UpdatePerson(id, name string) error {
	result := fr.dbEnv.Db.Model(&domain.FacePerson{}).Where("id = ?", id).
		Updates(map[string]any{"name": name, "updated_at": time.Now()})
	return db.HandleError(result)
}

// DeletePerson removes the person row and detaches all of its faces (faces
// return to the unassigned pool rather than being deleted).
func (fr *FaceRepository) DeletePerson(id string) error {
	return fr.dbEnv.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.FaceDetection{}).Where("person_id = ?", id).
			Updates(map[string]any{"person_id": nil, "status": "detected"}).Error; err != nil {
			return err
		}
		res := tx.Where("id = ?", id).Delete(&domain.FacePerson{})
		return db.HandleError(res)
	})
}

// UnassignedDetections returns detections with no person (review pool).
func (fr *FaceRepository) UnassignedDetections(limit int) ([]domain.FaceDetection, error) {
	var dets []domain.FaceDetection
	q := fr.dbEnv.Db.Where("person_id IS NULL").Order("created_at ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	result := q.Find(&dets)
	return dets, result.Error
}

// AllDetections returns detections up to limit (for full recluster passes).
func (fr *FaceRepository) AllDetections(limit int) ([]domain.FaceDetection, error) {
	var dets []domain.FaceDetection
	q := fr.dbEnv.Db.Order("created_at ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	result := q.Find(&dets)
	return dets, result.Error
}

// SetDetectionPerson assigns (or, with a nil personID, unassigns) a detection.
func (fr *FaceRepository) SetDetectionPerson(detID string, personID *string) error {
	status := "detected"
	if personID != nil {
		status = "assigned"
	}
	result := fr.dbEnv.Db.Model(&domain.FaceDetection{}).Where("id = ?", detID).
		Updates(map[string]any{"person_id": personID, "status": status, "updated_at": time.Now()})
	return db.HandleError(result)
}

// AssignDetectionsToPerson assigns a set of detections to a person in bulk.
func (fr *FaceRepository) AssignDetectionsToPerson(detIDs []string, personID string) error {
	if len(detIDs) == 0 {
		return nil
	}
	result := fr.dbEnv.Db.Model(&domain.FaceDetection{}).Where("id IN ?", detIDs).
		Updates(map[string]any{"person_id": personID, "status": "assigned", "updated_at": time.Now()})
	return result.Error
}

// FaceCount returns the total number of stored detections.
func (fr *FaceRepository) FaceCount() (int64, error) {
	var count int64
	result := fr.dbEnv.Db.Model(&domain.FaceDetection{}).Count(&count)
	return count, result.Error
}

// SetCoverFace records the best (thumbnail) face for a person.
func (fr *FaceRepository) SetCoverFace(personID, faceID string) error {
	result := fr.dbEnv.Db.Model(&domain.FacePerson{}).Where("id = ?", personID).
		Updates(map[string]any{"cover_face_id": faceID, "updated_at": time.Now()})
	return db.HandleError(result)
}

// UpsertCluster inserts or updates a cluster centroid.
func (fr *FaceRepository) UpsertCluster(c domain.FaceCluster) error {
	c.UpdatedAt = time.Now()
	result := fr.dbEnv.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&c)
	return result.Error
}

// DeleteCluster removes a cluster row.
func (fr *FaceRepository) DeleteCluster(id string) error {
	return fr.dbEnv.Db.Where("id = ?", id).Delete(&domain.FaceCluster{}).Error
}

// Clusters returns all cluster rows (for centroid lookups during matching).
func (fr *FaceRepository) Clusters() ([]domain.FaceCluster, error) {
	var clusters []domain.FaceCluster
	result := fr.dbEnv.Db.Find(&clusters)
	return clusters, result.Error
}

// FaceSetting reads the workspace opt-in flag (default false when unset).
func (fr *FaceRepository) FaceSetting() (bool, error) {
	var setting domain.Setting
	result := fr.dbEnv.Db.Where("key = ?", faceRecognitionSettingKey).First(&setting)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, db.HandleError(result)
	}
	return setting.Value == "true", nil
}

// SetFaceSetting writes the workspace opt-in flag.
func (fr *FaceRepository) SetFaceSetting(enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	setting := domain.Setting{
		Key:          faceRecognitionSettingKey,
		Value:        value,
		FriendlyName: "Face recognition",
		Category:     "Privacy",
		Type:         "boolean",
		Description:  "Detect and cluster faces locally to organize photos by person. Photos never leave this server.",
	}
	return fr.dbEnv.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&setting).Error
}

const faceRecognitionSettingKey = "face_recognition_enabled"
