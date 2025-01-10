package port

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock

// UtilityRepository is an interface for interacting with utility-related business logic
type UtilityRepository interface {
	// GetSetting returns a setting by key
	GetSetting(key string) (*domain.Setting, error)
	// GetAllSettings returns all settings
	GetAllSettings() ([]*domain.Setting, error)
	// UpdateSetting updates a single setting
	UpdateSetting(setting *domain.Setting) error
	// UpdateAllSettings updates all settings
	UpdateAllSettings(settings []*domain.Setting) error
	//CreateBaseSettings creates initial settings required for new instance
	CreateBaseSettings(reset bool) error
}

// UtilityService is an interface for interacting with user utility-related business logic
type UtilityService interface {
	// Healthcheck returns all settings as a healthcheck
	Healthcheck() ([]*domain.Setting, error)
	// GetSetting returns a setting by key
	GetSetting(key string) (*domain.Setting, error)
	// ListAllSettings returns all settings
	ListAllSettings() ([]*domain.Setting, error)
	//UpdateAllSettings updates all settings
	UpdateAllSettings(settings []*domain.Setting) error
	//CreateBaseSettings creates initial settings required for new instance
	CreateBaseSettings(reset bool) error
}
