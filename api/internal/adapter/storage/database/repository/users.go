package repository

import (
	db "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	dbEnv *db.Env
}

func NewUserRepository(dbEnv *db.Env) *UserRepository {
	return &UserRepository{
		dbEnv,
	}
}

func (ur *UserRepository) ListUsers(pageNumber, pageSize int) ([]domain.User, error) {
	var user []domain.User
	result := ur.dbEnv.Db.Scopes(db.Paginate(pageNumber, pageSize)).Find(&user)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetUserById(id string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "id = ? ", id)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "username = ? ", username)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetApprovedUserByUsername(username string) (*domain.User, error) {
	var user *domain.User
	result := ur.dbEnv.Db.Find(&user, "username = ? AND approved = true", username)
	err := db.HandleError(result)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	// Admin-created users join the shared default workspace, preserving the
	// home-instance shared-library semantics (issue #74 D6). Signups instead
	// create an isolated personal workspace via
	// CreateUserWithPersonalWorkspace. Both happen in one transaction so a
	// user is never left without a workspace.
	err := ur.dbEnv.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(user).Error; err != nil {
			return err
		}
		// Ensure the deterministic default workspace exists, then join it.
		if err := tx.Exec(`INSERT INTO photobox.workspaces (id, name, slug, owner_user_id, created_at, updated_at)
			VALUES (?, 'Default', 'default', ?, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING`, domain.DefaultWorkspaceID, user.ID).Error; err != nil {
			return err
		}
		return tx.Exec(`INSERT INTO photobox.workspace_members (workspace_id, user_id, role, created_at, updated_at)
			VALUES (?, ?, ?, NOW(), NOW())
			ON CONFLICT (workspace_id, user_id) DO NOTHING`, domain.DefaultWorkspaceID, user.ID, domain.WorkspaceMemberRole).Error
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUserWithPersonalWorkspace inserts a user, their personal workspace,
// and the owner membership in one transaction (issue #74 §8). If the
// personal-workspace insert fails (e.g. slug conflict), the user insert
// rolls back so no half-registered account survives.
func (ur *UserRepository) CreateUserWithPersonalWorkspace(user *domain.User, workspace *domain.Workspace) error {
	return ur.dbEnv.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(user).Error; err != nil {
			return err
		}
		if err := tx.Create(workspace).Error; err != nil {
			return err
		}
		member := &domain.WorkspaceMember{
			WorkspaceID: workspace.ID,
			UserID:      user.ID,
			Role:        domain.WorkspaceOwner,
		}
		return tx.Create(member).Error
	})
}

func (ur *UserRepository) UpdateUser(user *domain.User) error {
	result := ur.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Updates(&user)
	return result.Error
}

func (ur *UserRepository) DeleteUser(id string) error {
	var user *domain.User
	user.ID = id
	result := ur.dbEnv.Db.Delete(&user)
	return result.Error
}
