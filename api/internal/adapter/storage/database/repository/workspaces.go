package repository

import (
	"context"

	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm"
)

// WorkspaceRepository persists workspaces and memberships.
type WorkspaceRepository struct {
	dbEnv *db.Env
}

func NewWorkspaceRepository(dbEnv *db.Env) *WorkspaceRepository {
	return &WorkspaceRepository{dbEnv: dbEnv}
}

func (r *WorkspaceRepository) CreateWorkspace(workspace *domain.Workspace) error {
	return r.dbEnv.Db.Create(workspace).Error
}

func (r *WorkspaceRepository) GetWorkspaceById(id string) (*domain.Workspace, error) {
	var w domain.Workspace
	if err := r.dbEnv.Db.First(&w, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *WorkspaceRepository) GetWorkspaceBySlug(slug string) (*domain.Workspace, error) {
	var w domain.Workspace
	if err := r.dbEnv.Db.First(&w, "slug = ?", slug).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}
	return &w, nil
}

func (r *WorkspaceRepository) UpdateWorkspace(workspace *domain.Workspace) error {
	return r.dbEnv.Db.Model(&domain.Workspace{}).Where("id = ?", workspace.ID).
		Updates(map[string]interface{}{
			"name":       workspace.Name,
			"slug":       workspace.Slug,
			"settings":   workspace.Settings,
			"updated_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *WorkspaceRepository) DeleteWorkspace(id string) error {
	err := r.dbEnv.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("workspace_id = ?", id).Delete(&domain.WorkspaceMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.Workspace{}, "id = ?", id).Error
	})
	return err
}

func (r *WorkspaceRepository) AddMember(workspaceID, userID string, role domain.WorkspaceRole) error {
	m := &domain.WorkspaceMember{WorkspaceID: workspaceID, UserID: userID, Role: role}
	return r.dbEnv.Db.Create(m).Error
}

func (r *WorkspaceRepository) RemoveMember(workspaceID, userID string) error {
	return r.dbEnv.Db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Delete(&domain.WorkspaceMember{}).Error
}

func (r *WorkspaceRepository) UpdateMemberRole(workspaceID, userID string, role domain.WorkspaceRole) error {
	return r.dbEnv.Db.Model(&domain.WorkspaceMember{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Update("role", role).Error
}

func (r *WorkspaceRepository) GetMembership(workspaceID, userID string) (*domain.WorkspaceMember, error) {
	var m domain.WorkspaceMember
	if err := r.dbEnv.Db.First(&m, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *WorkspaceRepository) ListMembers(workspaceID string) ([]domain.WorkspaceMember, error) {
	var members []domain.WorkspaceMember
	err := r.dbEnv.Db.Where("workspace_id = ?", workspaceID).
		Order("created_at ASC").Find(&members).Error
	return members, err
}

func (r *WorkspaceRepository) ListWorkspacesForUser(userID string) ([]domain.Workspace, error) {
	var workspaces []domain.Workspace
	err := r.dbEnv.Db.
		Joins("JOIN photobox.workspace_members wm ON wm.workspace_id = photobox.workspaces.id AND wm.user_id = ?", userID).
		Order("photobox.workspaces.created_at DESC").
		Find(&workspaces).Error
	return workspaces, err
}

func (r *WorkspaceRepository) ListMemberUsers(workspaceID string) ([]domain.User, error) {
	var users []domain.User
	err := r.dbEnv.Db.
		Joins("JOIN photobox.workspace_members wm ON wm.user_id = photobox.users.id AND wm.workspace_id = ?", workspaceID).
		Find(&users).Error
	return users, err
}

func (r *WorkspaceRepository) AdjustStorageUsed(ctx context.Context, workspaceID string, delta int64) (int64, error) {
	var newTotal int64
	err := r.dbEnv.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Atomic update with row lock; never below zero.
		res := tx.Model(&domain.Workspace{}).
			Where("id = ?", workspaceID).
			UpdateColumn("storage_used_bytes", gorm.Expr("GREATEST(storage_used_bytes + ?, 0)", delta))
		if res.Error != nil {
			return res.Error
		}
		return tx.Model(&domain.Workspace{}).Select("storage_used_bytes").
			Where("id = ?", workspaceID).Scan(&newTotal).Error
	})
	return newTotal, err
}
