package repository

import (
	"time"

	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

type ApiKeyRepository struct {
	dbEnv *db.Env
}

func NewApiKeyRepository(dbEnv *db.Env) *ApiKeyRepository {
	return &ApiKeyRepository{dbEnv}
}

func (r *ApiKeyRepository) Create(key *domain.ApiKey) error {
	return r.dbEnv.Db.Create(key).Error
}

func (r *ApiKeyRepository) GetByHash(hash string) (*domain.ApiKey, error) {
	var key domain.ApiKey
	result := r.dbEnv.Db.Where("key_hash = ? AND revoked = ?", hash, false).First(&key)
	if err := db.HandleError(result); err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *ApiKeyRepository) List() ([]*domain.ApiKey, error) {
	var keys []*domain.ApiKey
	result := r.dbEnv.Db.Order("created_at DESC").Find(&keys)
	if err := db.HandleError(result); err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *ApiKeyRepository) Revoke(id string) error {
	return r.dbEnv.Db.Model(&domain.ApiKey{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *ApiKeyRepository) TouchLastUsed(id string) error {
	now := time.Now()
	return r.dbEnv.Db.Model(&domain.ApiKey{}).Where("id = ?", id).Update("last_used", now).Error
}
