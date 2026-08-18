package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/adapter/handler/auth"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type ShareService struct {
	repo      port.ShareRepository
	photoSvc  port.PhotoService
	albumSvc  port.AlbumService
}

func NewShareService(repo port.ShareRepository, photoSvc port.PhotoService, albumSvc port.AlbumService) *ShareService {
	return &ShareService{repo, photoSvc, albumSvc}
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
		// Hash share passwords at rest (argon2id, matching user passwords).
		hashed, err := auth.CreateHash(*password, auth.DefaultArgon2idHash())
		if err != nil {
			return nil, err
		}
		passwordHash = hashed
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
		if password == nil {
			return nil, domain.ErrSharedLinkPasswordRequired
		}
		// Verify against the argon2id hash. Legacy plaintext rows (created
		// before hashing) are compared directly and re-hashed on success.
		match, err := auth.ComparePasswordAndHash(*password, share.PasswordHash)
		if err != nil {
			// Not an argon2id hash — treat as legacy plaintext.
			if *password != share.PasswordHash {
				return nil, domain.ErrSharedLinkPasswordRequired
			}
		} else if !match {
			return nil, domain.ErrSharedLinkPasswordRequired
		}
	}

	_ = ss.repo.IncrementViewCount(token)
	return share, nil
}

func (ss *ShareService) GetSharedResourceData(token string, password *string) (*port.SharedResourceData, error) {
	share, err := ss.GetSharedResource(token, password)
	if err != nil {
		return nil, err
	}

	var resource interface{}
	switch share.ResourceType {
	case "photo":
		photo, err := ss.photoSvc.GetPhoto(share.ResourceId, false)
		if err != nil {
			return nil, err
		}
		resource = photo
	case "album":
		album, err := ss.albumSvc.GetAlbumById(share.ResourceId)
		if err != nil {
			return nil, err
		}
		resource = album
	}

	return &port.SharedResourceData{
		Share:    share,
		Resource: resource,
	}, nil
}

// ListShares returns shares. When createdBy is non-empty, only that user's
// shares are returned (owner scoping — issue #151). Empty createdBy lists all
// (admin/ops path).
func (ss *ShareService) ListShares(createdBy string) ([]*domain.SharedLink, error) {
	if createdBy != "" {
		return ss.repo.ListSharesByOwner(createdBy)
	}
	return ss.repo.ListShares()
}

// RevokeShare deletes a share. When createdBy is non-empty, the share is only
// deleted if it belongs to that user (owner scoping — issue #151); otherwise
// it is treated as not found to avoid leaking existence.
func (ss *ShareService) RevokeShare(token, createdBy string) error {
	if createdBy != "" {
		share, err := ss.repo.GetShareByToken(token)
		if err != nil {
			return domain.ErrDataNotFound
		}
		if share.CreatedBy != createdBy {
			return domain.ErrDataNotFound
		}
	}
	return ss.repo.DeleteShare(token)
}
