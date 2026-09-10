package service

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

type AlbumService struct {
	repo   port.AlbumRepository
	config appconfig.AppConfig
}

// NewAlbumService creates a new Album service instance
func NewAlbumService(repo port.AlbumRepository, config appconfig.AppConfig) *AlbumService {
	return &AlbumService{
		repo,
		config,
	}
}

// WithWorkspace returns a copy of the service whose queries are scoped to the
// given workspace (issue #74).
func (as *AlbumService) WithWorkspace(workspaceID string) port.AlbumService {
	c := *as
	c.repo = as.repo.WithWorkspace(workspaceID)
	return &c
}

func (as *AlbumService) GetAlbumById(id string) (*domain.Album, error) {
	return as.repo.GetAlbumById(id)
}

func (as *AlbumService) GetAlbumByName(name string) (*domain.Album, error) {
	return as.repo.GetAlbumByName(name)
}

func (as *AlbumService) ListAlbums(fromId string, limit int) ([]*domain.Album, error) {
	return as.repo.ListAllAlbums(fromId, limit)
}

func (as *AlbumService) AlbumCount() (int64, error) {
	return as.repo.AlbumCount()
}

func (as *AlbumService) CreateAlbum(name string) (*domain.Album, error) {
	return as.repo.CreateAlbum(name)
}

// CreateSmartAlbum creates an album whose contents are defined by filter
// rules rather than explicit membership. The rules are stored in the
// album's metadata along with a smart marker.
func (as *AlbumService) CreateSmartAlbum(name string, rules domain.SmartAlbumRules) (*domain.Album, error) {
	meta, err := json.Marshal(map[string]any{
		"smart": true,
		"rules": rules,
	})
	if err != nil {
		return nil, err
	}

	album := &domain.Album{
		Name:     name,
		Metadata: meta,
	}
	if err := as.repo.CreateAlbumWithMetadata(album); err != nil {
		return nil, err
	}
	return album, nil
}

// UpdateSmartAlbum replaces the rules of a smart album.
func (as *AlbumService) UpdateSmartAlbum(id string, rules domain.SmartAlbumRules) (*domain.Album, error) {
	album, err := as.repo.GetAlbumById(id)
	if err != nil {
		return nil, err
	}
	meta, err := json.Marshal(map[string]any{
		"smart": true,
		"rules": rules,
	})
	if err != nil {
		return nil, err
	}
	album.Metadata = meta
	if err := as.repo.UpdateAlbum(album); err != nil {
		return nil, err
	}
	return album, nil
}

func (as *AlbumService) UpdateAlbum(id string, updates map[string]any) (*domain.Album, error) {
	album, err := as.repo.GetAlbumById(id)
	if err != nil {
		return nil, err
	}

	oldName := album.Name

	if name, ok := updates["name"].(string); ok && name != "" && name != album.Name {
		album.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		album.Description = description
	}
	if coverPhotoId, ok := updates["coverPhotoId"].(string); ok {
		album.CoverPhotoId = coverPhotoId
	}
	if tags, ok := updates["tags"].(string); ok {
		album.Tags = tags
	}

	err = as.repo.UpdateAlbum(album)
	if err != nil {
		return nil, err
	}

	// Rename filesystem directory if name changed
	if oldName != album.Name {
		oldPath := filepath.Join(as.config.PhotoDir, oldName)
		newPath := filepath.Join(as.config.PhotoDir, album.Name)
		// We'll handle this at the handler/service level by passing filesystem service
		_ = oldPath
		_ = newPath
	}

	return album, nil
}

func (as *AlbumService) SearchAlbums(query string, limit int) ([]*domain.Album, error) {
	return as.repo.SearchAlbums(query, limit)
}

func (as *AlbumService) DeleteAlbum(id string, deletePhotos bool) error {
	if deletePhotos {
		// Move photos to trash instead of deleting
		// This requires photo service; we'll handle at handler level
		return as.repo.DeleteAlbum(id)
	}

	// Create or find "Uncategorized" album
	uncategorized, err := as.repo.GetAlbumByName("Uncategorized")
	if err != nil || uncategorized == nil {
		uncategorized, err = as.repo.CreateAlbum("Uncategorized")
		if err != nil {
			return fmt.Errorf("failed to create uncategorized album: %w", err)
		}
	}

	err = as.repo.ReassignPhotosToAlbum(id, uncategorized.ID)
	if err != nil {
		return fmt.Errorf("failed to reassign photos: %w", err)
	}

	return as.repo.DeleteAlbum(id)
}
