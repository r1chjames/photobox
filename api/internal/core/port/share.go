package port

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// ShareRepository is an interface for interacting with share-related data
type ShareRepository interface {
	CreateShare(share *domain.SharedLink) error
	GetShareByToken(token string) (*domain.SharedLink, error)
	ListShares() ([]*domain.SharedLink, error)
	DeleteShare(token string) error
	IncrementViewCount(token string) error
}

// ShareService is an interface for interacting with share-related business logic
type ShareService interface {
	CreateShare(resourceType, resourceId, createdBy string, expiry *string, password *string) (*domain.SharedLink, error)
	GetSharedResource(token string, password *string) (*domain.SharedLink, error)
	ListShares() ([]*domain.SharedLink, error)
	RevokeShare(token string) error
}
