package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type ShareService struct {
	repo port.ShareRepository
}

func NewShareService(repo port.ShareRepository) *ShareService {
	return &ShareService{repo}
}

func (ss *ShareService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (ss *ShareService) CreateShare(resourceType, resourceId, createdBy string, expiry *string, password *string) (*domain.SharedLink, error) {
	token, err := ss.generateToken()
	if err != nil {
		return nil, err
	}

	var expiryTime *time.Time
	if expiry != nil && *expiry != "" {
		t, err := time.Parse(time.RFC3339, *expiry)
		if err != nil {
			return nil, err
		}
		expiryTime = &t
	}

	var passwordHash string
	if password != nil && *password != "" {
		// In a real app, hash the password with bcrypt/argon2
		// For simplicity, we store it plaintext here (NOT for production)
		passwordHash = *password
	}

	share := &domain.SharedLink{
		Token:        token,
		ResourceType: resourceType,
		ResourceId:   resourceId,
		CreatedBy:    createdBy,
		Expiry:       expiryTime,
		PasswordHash: passwordHash,
	}

	err = ss.repo.CreateShare(share)
	if err != nil {
		return nil, err
	}
	return share, nil
}

func (ss *ShareService) GetSharedResource(token string, password *string) (*domain.SharedLink, error) {
	share, err := ss.repo.GetShareByToken(token)
	if err != nil {
		return nil, domain.ErrDataNotFound
	}

	if share.Expiry != nil && share.Expiry.Before(time.Now()) {
		return nil, domain.ErrSharedLinkExpired
	}

	if share.PasswordHash != "" {
		if password == nil || *password != share.PasswordHash {
			return nil, domain.ErrSharedLinkPasswordRequired
		}
	}

	_ = ss.repo.IncrementViewCount(token)
	return share, nil
}

func (ss *ShareService) ListShares() ([]*domain.SharedLink, error) {
	return ss.repo.ListShares()
}

func (ss *ShareService) RevokeShare(token string) error {
	return ss.repo.DeleteShare(token)
}
