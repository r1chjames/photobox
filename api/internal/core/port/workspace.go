package port

import (
	"context"

	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

//go:generate mockgen -source=workspace.go -destination=mock/workspace.go -package=mock

// WorkspaceRepository persists workspaces and memberships.
type WorkspaceRepository interface {
	// CreateWorkspace inserts a new workspace.
	CreateWorkspace(workspace *domain.Workspace) error
	// GetWorkspaceById returns a workspace or ErrDataNotFound.
	GetWorkspaceById(id string) (*domain.Workspace, error)
	// GetWorkspaceBySlug returns a workspace or ErrDataNotFound.
	GetWorkspaceBySlug(slug string) (*domain.Workspace, error)
	// UpdateWorkspace persists name/slug/settings changes.
	UpdateWorkspace(workspace *domain.Workspace) error
	// DeleteWorkspace removes a workspace and its memberships.
	DeleteWorkspace(id string) error

	// AddMember inserts a membership row. Duplicate membership is an error.
	AddMember(workspaceID, userID string, role domain.WorkspaceRole) error
	// RemoveMember deletes a membership row.
	RemoveMember(workspaceID, userID string) error
	// UpdateMemberRole changes a member's role.
	UpdateMemberRole(workspaceID, userID string, role domain.WorkspaceRole) error
	// GetMembership returns the caller's membership in a workspace (nil when
	// the user is not a member). The single membership lookup that backs
	// per-request workspace authorization.
	GetMembership(workspaceID, userID string) (*domain.WorkspaceMember, error)
	// ListMembers returns all members of a workspace.
	ListMembers(workspaceID string) ([]domain.WorkspaceMember, error)
	// ListWorkspacesForUser returns the workspaces the user belongs to,
	// newest first.
	ListWorkspacesForUser(userID string) ([]domain.Workspace, error)
	// ListMemberUsers returns the users who belong to a workspace.
	ListMemberUsers(workspaceID string) ([]domain.User, error)
	// AdjustStorageUsed atomically adds delta to a workspace's
	// storage_used_bytes (negative on delete). Returns the new total.
	AdjustStorageUsed(ctx context.Context, workspaceID string, delta int64) (int64, error)
}

// WorkspaceService is the application-facing workspace API.
type WorkspaceService interface {
	// CreateWorkspace creates a workspace owned by userID and adds the owner
	// as a member.
	CreateWorkspace(name, slug, userID string) (*domain.Workspace, error)
	// GetWorkspace returns a workspace if the user is a member.
	GetWorkspace(workspaceID, userID string) (*domain.Workspace, error)
	// ListWorkspaces returns the user's workspaces.
	ListWorkspaces(userID string) ([]domain.Workspace, error)
	// UpdateWorkspace renames a workspace (owner/admin only).
	UpdateWorkspace(workspaceID, userID, name string) (*domain.Workspace, error)
	// DeleteWorkspace removes a workspace (owner only).
	DeleteWorkspace(workspaceID, userID string) error
	// AddMember adds a user to a workspace (owner/admin only).
	AddMember(workspaceID, actorUserID, targetUserID string, role domain.WorkspaceRole) error
	// RemoveMember removes a user from a workspace (owner/admin only; owner
	// cannot remove themselves while others remain).
	RemoveMember(workspaceID, actorUserID, targetUserID string) error
	// UpdateMemberRole changes a member's role (owner/admin only).
	UpdateMemberRole(workspaceID, actorUserID, targetUserID string, role domain.WorkspaceRole) error
	// ListMembers returns the workspace members (any workspace member).
	ListMembers(workspaceID, userID string) ([]domain.WorkspaceMember, error)
	// GetMembership returns the user's role in a workspace (nil if none).
	GetMembership(workspaceID, userID string) (*domain.WorkspaceMember, error)
	// WorkspaceExists reports whether a workspace ID exists. Used to resolve
	// deployment-level credentials (API keys), which carry no membership.
	WorkspaceExists(workspaceID string) (bool, error)
	// CanManage reports whether the user may manage members/settings
	// (owner or admin).
	CanManage(membership *domain.WorkspaceMember) bool
}
