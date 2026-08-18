package repository

import (
	"time"

	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

type ShareRepository struct {
	dbEnv *db.Env
}

func NewShareRepository(dbEnv *db.Env) *ShareRepository {
	return &ShareRepository{
		dbEnv,
	}
}

func (sr *ShareRepository) CreateShare(share *domain.SharedLink) error {
	result := sr.dbEnv.Db.Create(share)
	return result.Error
}

func (sr *ShareRepository) GetShareByToken(token string) (*domain.SharedLink, error) {
	var share domain.SharedLink
	result := sr.dbEnv.Db.First(&share, "token = ?", token)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (sr *ShareRepository) ListShares() ([]*domain.SharedLink, error) {
	var shares []*domain.SharedLink
	result := sr.dbEnv.Db.Find(&shares)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

// ListSharesByOwner returns only the shares created by the given user
// (owner scoping — issue #151).
func (sr *ShareRepository) ListSharesByOwner(createdBy string) ([]*domain.SharedLink, error) {
	var shares []*domain.SharedLink
	result := sr.dbEnv.Db.Where("created_by = ?", createdBy).Find(&shares)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (sr *ShareRepository) DeleteShare(token string) error {
	result := sr.dbEnv.Db.Delete(&domain.SharedLink{}, "token = ?", token)
	return db.HandleError(result)
}

func (sr *ShareRepository) IncrementViewCount(token string) error {
	result := sr.dbEnv.Db.Model(&domain.SharedLink{}).
		Where("token = ?", token).
		Update("view_count", sr.dbEnv.Db.Raw("view_count + 1")).
		Update("updated_at", time.Now())
	return result.Error
}
