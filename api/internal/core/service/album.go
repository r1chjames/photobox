package service

import (
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
