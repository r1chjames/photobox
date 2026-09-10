package repository

import (
	"time"

	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gorm.io/gorm"
)

type ShareRepository struct {
	dbEnv *db.Env
	// workspaceID scopes tenant-facing share queries to one workspace (issue
	// #74). Empty means unscoped, reserved for the anonymous share-resolution
	// path (where the share token itself is the credential) and for tests.
	workspaceID string
	// scoped records explicit binding via WithWorkspace; once scoped an empty
	// workspace fails closed rather than widening the query.
	scoped bool
}

func NewShareRepository(dbEnv *db.Env) *ShareRepository {
	return &ShareRepository{
		dbEnv: dbEnv,
	}
}

// WithWorkspace returns a repository scoped to workspaceID. Fail-closed: an
// empty workspace matches nothing, so a missing workspace cannot widen a query.
func (sr *ShareRepository) WithWorkspace(workspaceID string) port.ShareRepository {
	return &ShareRepository{dbEnv: sr.dbEnv, workspaceID: workspaceID, scoped: true}
}

// scope applies the workspace filter to a shared_links query. Unscoped
// repositories pass through untouched.
func (sr *ShareRepository) scope(tx *gorm.DB) *gorm.DB {
	if !sr.scoped {
		return tx
	}
	if sr.workspaceID == "" {
		return tx.Where("1 = 0") // fail closed
	}
	return tx.Where("shared_links.workspace_id = ?", sr.workspaceID)
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
	result := sr.scope(sr.dbEnv.Db).Find(&shares)
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
	result := sr.scope(sr.dbEnv.Db).Where("created_by = ?", createdBy).Find(&shares)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (sr *ShareRepository) DeleteShare(token string) error {
	result := sr.scope(sr.dbEnv.Db).Delete(&domain.SharedLink{}, "token = ?", token)
	return db.HandleError(result)
}

func (sr *ShareRepository) IncrementViewCount(token string) error {
	result := sr.dbEnv.Db.Model(&domain.SharedLink{}).
		Where("token = ?", token).
		Update("view_count", sr.dbEnv.Db.Raw("view_count + 1")).
		Update("updated_at", time.Now())
	return result.Error
}
