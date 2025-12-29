package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockUtilityRepository is a mock implementation of port.UtilityRepository
type MockUtilityRepository struct {
	mock.Mock
}

func (m *MockUtilityRepository) GetAllSettings() ([]*domain.Setting, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Setting), args.Error(1)
}

func (m *MockUtilityRepository) GetSetting(key string) (*domain.Setting, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Setting), args.Error(1)
}

func (m *MockUtilityRepository) UpdateSetting(setting *domain.Setting) error {
	args := m.Called(setting)
	return args.Error(0)
}

func (m *MockUtilityRepository) UpdateAllSettings(settings []*domain.Setting) error {
	args := m.Called(settings)
	return args.Error(0)
}

func (m *MockUtilityRepository) CreateBaseSettings(reset bool) error {
	args := m.Called(reset)
	return args.Error(0)
}

// TestHealthcheck tests the healthcheck functionality
func TestHealthcheck(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockUtilityRepository)
		expectedError error
		validate      func(*testing.T, []*domain.Setting, error)
	}{
		{
			name: "successful healthcheck with settings",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetAllSettings").Return([]*domain.Setting{
					{Key: "app_version", Value: "1.0.0", Description: "Application version"},
					{Key: "db_status", Value: "connected", Description: "Database status"},
					{Key: "default_new_albums_dir", Value: "/photos", Description: "Default album directory"},
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, settings []*domain.Setting, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, settings)
				assert.Equal(t, 3, len(settings))
				assert.Equal(t, "app_version", settings[0].Key)
				assert.Equal(t, "1.0.0", settings[0].Value)
			},
		},
		{
			name: "healthcheck with empty settings",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetAllSettings").Return([]*domain.Setting{}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, settings []*domain.Setting, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, settings)
				assert.Equal(t, 0, len(settings))
			},
		},
		{
			name: "repository error during healthcheck",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetAllSettings").Return(nil, domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, settings []*domain.Setting, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Nil(t, settings)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUtilityRepository)
			tt.mockSetup(mockRepo)

			service := NewUtilityService(mockRepo)
			result, err := service.Healthcheck()

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetSetting tests retrieving a single setting by key
func TestGetSetting(t *testing.T) {
	tests := []struct {
		name          string
		settingKey    string
		mockSetup     func(*MockUtilityRepository)
		expectedError error
		validate      func(*testing.T, *domain.Setting, error)
	}{
		{
			name:       "successful setting retrieval",
			settingKey: "default_new_albums_dir",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetSetting", "default_new_albums_dir").Return(&domain.Setting{
					Key:         "default_new_albums_dir",
					Value:       "/photos",
					Description: "Default album directory",
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, setting *domain.Setting, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, setting)
				assert.Equal(t, "default_new_albums_dir", setting.Key)
				assert.Equal(t, "/photos", setting.Value)
			},
		},
		{
			name:       "setting not found",
			settingKey: "nonexistent_key",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetSetting", "nonexistent_key").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrDataNotFound,
			validate: func(t *testing.T, setting *domain.Setting, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrDataNotFound, err)
				assert.Nil(t, setting)
			},
		},
		{
			name:       "repository error retrieving setting",
			settingKey: "app_version",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetSetting", "app_version").Return(nil, domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, setting *domain.Setting, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Nil(t, setting)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUtilityRepository)
			tt.mockSetup(mockRepo)

			service := NewUtilityService(mockRepo)
			result, err := service.GetSetting(tt.settingKey)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListAllSettings tests listing all settings
func TestListAllSettings(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockUtilityRepository)
		expectedError error
		validate      func(*testing.T, []*domain.Setting, error)
	}{
		{
			name: "successful settings listing",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetAllSettings").Return([]*domain.Setting{
					{Key: "setting1", Value: "value1", Description: "First setting"},
					{Key: "setting2", Value: "value2", Description: "Second setting"},
					{Key: "setting3", Value: "value3", Description: "Third setting"},
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, settings []*domain.Setting, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, settings)
				assert.Equal(t, 3, len(settings))
				assert.Equal(t, "setting1", settings[0].Key)
				assert.Equal(t, "value1", settings[0].Value)
			},
		},
		{
			name: "empty settings list",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetAllSettings").Return([]*domain.Setting{}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, settings []*domain.Setting, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, settings)
				assert.Equal(t, 0, len(settings))
			},
		},
		{
			name: "repository error listing settings",
			mockSetup: func(m *MockUtilityRepository) {
				m.On("GetAllSettings").Return(nil, domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, settings []*domain.Setting, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Nil(t, settings)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUtilityRepository)
			tt.mockSetup(mockRepo)

			service := NewUtilityService(mockRepo)
			result, err := service.ListAllSettings()

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateAllSettings tests updating all settings
func TestUpdateAllSettings(t *testing.T) {
	tests := []struct {
		name          string
		settings      []*domain.Setting
		mockSetup     func(*MockUtilityRepository)
		expectedError error
		validate      func(*testing.T, error)
	}{
		{
			name: "successful settings update",
			settings: []*domain.Setting{
				{Key: "setting1", Value: "new_value1", Description: "Updated setting 1"},
				{Key: "setting2", Value: "new_value2", Description: "Updated setting 2"},
			},
			mockSetup: func(m *MockUtilityRepository) {
				m.On("UpdateAllSettings", mock.MatchedBy(func(settings []*domain.Setting) bool {
					return len(settings) == 2 && settings[0].Value == "new_value1"
				})).Return(nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:     "update with empty settings list",
			settings: []*domain.Setting{},
			mockSetup: func(m *MockUtilityRepository) {
				m.On("UpdateAllSettings", mock.MatchedBy(func(settings []*domain.Setting) bool {
					return len(settings) == 0
				})).Return(nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "repository error updating settings",
			settings: []*domain.Setting{
				{Key: "setting1", Value: "value1", Description: "Setting 1"},
			},
			mockSetup: func(m *MockUtilityRepository) {
				m.On("UpdateAllSettings", mock.Anything).Return(domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUtilityRepository)
			tt.mockSetup(mockRepo)

			service := NewUtilityService(mockRepo)
			err := service.UpdateAllSettings(tt.settings)

			tt.validate(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestCreateBaseSettings tests creating base settings
func TestCreateBaseSettings(t *testing.T) {
	tests := []struct {
		name          string
		reset         bool
		mockSetup     func(*MockUtilityRepository)
		expectedError error
		validate      func(*testing.T, error)
	}{
		{
			name:  "create base settings without reset",
			reset: false,
			mockSetup: func(m *MockUtilityRepository) {
				m.On("CreateBaseSettings", false).Return(nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:  "create base settings with reset",
			reset: true,
			mockSetup: func(m *MockUtilityRepository) {
				m.On("CreateBaseSettings", true).Return(nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:  "repository error creating base settings",
			reset: false,
			mockSetup: func(m *MockUtilityRepository) {
				m.On("CreateBaseSettings", false).Return(domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUtilityRepository)
			tt.mockSetup(mockRepo)

			service := NewUtilityService(mockRepo)
			err := service.CreateBaseSettings(tt.reset)

			tt.validate(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
