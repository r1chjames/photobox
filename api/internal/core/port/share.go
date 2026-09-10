package port

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// ShareRepository is an interface for interacting with share-related data
type ShareRepository interface {
	// WithWorkspace returns a repository whose tenant queries are scoped to
	// the given workspace (issue #74). An empty workspace is fail-closed.
	WithWorkspace(workspaceID string) ShareRepository
	CreateShare(share *domain.SharedLink) error
	GetShareByToken(token string) (*domain.SharedLink, error)
	ListShares() ([]*domain.SharedLink, error)
	ListSharesByOwner(createdBy string) ([]*domain.SharedLink, error)
	DeleteShare(token string) error
	IncrementViewCount(token string) error
}

// SharedResourceData wraps a shared link with its actual resource
type SharedResourceData struct {
	Share    *domain.SharedLink `json:"share"`
	Resource interface{}        `json:"resource"`
}

// ShareService is an interface for interacting with share-related business logic
type ShareService interface {
	// WithWorkspace returns a service scoped to the given workspace
	// (issue #74). Handlers must use this; empty is fail-closed.
	WithWorkspace(workspaceID string) ShareService
	CreateShare(resourceType, resourceId, createdBy string, expiry *string, password *string) (*domain.SharedLink, error)
	GetSharedResource(token string, password *string) (*domain.SharedLink, error)
	GetSharedResourceData(token string, password *string) (*SharedResourceData, error)
	ListShares(createdBy string) ([]*domain.SharedLink, error)
	RevokeShare(token, createdBy string) error
}
