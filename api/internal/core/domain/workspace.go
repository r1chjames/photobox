package domain

import "time"

// WorkspaceRole is a member's role within a workspace. The workspace role is
// orthogonal to the deployment-level UserRole on the token (ADMINISTRATOR
// etc.); a user may hold different workspace roles in different workspaces.
type WorkspaceRole string

const (
	// WorkspaceOwner is the workspace creator; can transfer/delete the
	// workspace and manage all members.
	WorkspaceOwner WorkspaceRole = "owner"
	// WorkspaceAdmin manages members and workspace settings.
	WorkspaceAdmin WorkspaceRole = "admin"
	// WorkspaceMemberRole can read and add content (photos, albums, tags).
	WorkspaceMemberRole WorkspaceRole = "member"
	// WorkspaceViewer is read-only.
	WorkspaceViewer WorkspaceRole = "viewer"
)

// Workspace is the isolation container for a photo library (issue #74).
// Every content row (photos, albums, shared_links) carries a workspace_id;
// access is mediated by a workspace_members row, never by the token alone.
type Workspace struct {
	ID          string `gorm:"primaryKey;size:36" json:"id"`
	Name        string `json:"name"`
	Slug        string `gorm:"uniqueIndex;size:64" json:"slug"`
	OwnerUserID string `gorm:"index;size:36" json:"ownerUserId"`
	// StorageUsedBytes tracks the sum of original-bytes stored in this
	// workspace; maintained transactionally on upload/delete (issue #74).
	StorageUsedBytes int64 `json:"storageUsedBytes" gorm:"default:0"`
	// StorageLimitBytes caps the workspace; 0 = unlimited. Quota columns ship
	// in Phase 1 but enforcement lands with the upload path (Phase 2).
	StorageLimitBytes int64 `json:"storageLimitBytes" gorm:"default:0"`
	Settings          *string   `json:"settings" gorm:"type:jsonb"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// WorkspaceMember is a (workspace, user) membership with a workspace role.
type WorkspaceMember struct {
	WorkspaceID string        `gorm:"primaryKey;size:36" json:"workspaceId"`
	UserID      string        `gorm:"primaryKey;size:36" json:"userId"`
	Role        WorkspaceRole `gorm:"size:16" json:"role"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// DefaultWorkspaceID is the deterministic UUID used for the single shared
// workspace that adopts pre-multi-tenancy data (issue #74 D6). Existing
// deployments keep today's shared-library semantics: every existing user is
// a member of this one workspace and all pre-existing content lives in it.
const DefaultWorkspaceID = "00000000-0000-0000-0000-000000000001"
