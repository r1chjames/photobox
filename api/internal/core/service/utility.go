package service

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

/**
 * UtilityService implements port.UtilityService interface
 * and provides an access to the utility repository
 */
type UtilityService struct {
	repo port.UtilityRepository
}

// NewUtilityService creates a new auth service instance
func NewUtilityService(repo port.UtilityRepository) *UtilityService {
	return &UtilityService{
		repo,
	}
}

func (us *UtilityService) Healthcheck() ([]*domain.Setting, error) {
	return us.repo.GetAllSettings()
}

func (us *UtilityService) GetSetting(key string) (*domain.Setting, error) {
	return us.repo.GetSetting(key)
}

func (us *UtilityService) ListAllSettings() ([]*domain.Setting, error) {
	return us.repo.GetAllSettings()
}

func (us *UtilityService) UpdateAllSettings(settings []*domain.Setting) error {
	return us.repo.UpdateAllSettings(settings)
}

func (us *UtilityService) CreateBaseSettings(reset bool) error {
	return us.repo.CreateBaseSettings(reset)
}
