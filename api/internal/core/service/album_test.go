package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockAlbumRepository is a mock implementation of port.AlbumRepository
type MockAlbumRepository struct {
	mock.Mock
}

func (m *MockAlbumRepository) GetAlbumById(id string) (*domain.Album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumRepository) GetAlbumByName(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumRepository) ListAllAlbums(fromId string, limit int) ([]*domain.Album, error) {
	args := m.Called(fromId, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Album), args.Error(1)
}

func (m *MockAlbumRepository) AlbumCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAlbumRepository) CreateAlbum(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumRepository) CreateAlbumIfNotExists(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

// TestGetAlbumById tests retrieving an album by ID
func TestGetAlbumById(t *testing.T) {
	tests := []struct {
		name          string
		albumID       string
		mockSetup     func(*MockAlbumRepository)
		expectedError error
		validate      func(*testing.T, *domain.Album, error)
	}{
		{
			name:    "successful album retrieval",
			albumID: "album123",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumById", "album123").Return(&domain.Album{
					ID:          "album123",
					Name:        "Test Album",
					Description: "Test Description",
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, album)
				assert.Equal(t, "album123", album.ID)
				assert.Equal(t, "Test Album", album.Name)
				assert.Equal(t, "Test Description", album.Description)
			},
		},
		{
			name:    "album not found",
			albumID: "nonexistent",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumById", "nonexistent").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrDataNotFound,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrDataNotFound, err)
				assert.Nil(t, album)
			},
		},
		{
			name:    "repository error",
			albumID: "album456",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumById", "album456").Return(nil, domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Nil(t, album)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlbumRepository)
			tt.mockSetup(mockRepo)

			service := NewAlbumService(mockRepo, appconfig.AppConfig{})
			result, err := service.GetAlbumById(tt.albumID)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetAlbumByName tests retrieving an album by name
func TestGetAlbumByName(t *testing.T) {
	tests := []struct {
		name          string
		albumName     string
		mockSetup     func(*MockAlbumRepository)
		expectedError error
		validate      func(*testing.T, *domain.Album, error)
	}{
		{
			name:      "successful album retrieval by name",
			albumName: "Vacation 2024",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumByName", "Vacation 2024").Return(&domain.Album{
					ID:          "album789",
					Name:        "Vacation 2024",
					Description: "Summer vacation photos",
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, album)
				assert.Equal(t, "album789", album.ID)
				assert.Equal(t, "Vacation 2024", album.Name)
			},
		},
		{
			name:      "album name not found",
			albumName: "Nonexistent Album",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumByName", "Nonexistent Album").Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrDataNotFound,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrDataNotFound, err)
				assert.Nil(t, album)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlbumRepository)
			tt.mockSetup(mockRepo)

			service := NewAlbumService(mockRepo, appconfig.AppConfig{})
			result, err := service.GetAlbumByName(tt.albumName)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListAlbums tests listing albums with pagination
func TestListAlbums(t *testing.T) {
	tests := []struct {
		name          string
		fromID        string
		limit         int
		mockSetup     func(*MockAlbumRepository)
		expectedError error
		validate      func(*testing.T, []*domain.Album, error)
	}{
		{
			name:   "successful album listing without pagination",
			fromID: "",
			limit:  10,
			mockSetup: func(m *MockAlbumRepository) {
				m.On("ListAllAlbums", "", 10).Return([]*domain.Album{
					{ID: "album1", Name: "Album One", Description: "First album"},
					{ID: "album2", Name: "Album Two", Description: "Second album"},
					{ID: "album3", Name: "Album Three", Description: "Third album"},
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, albums []*domain.Album, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, albums)
				assert.Equal(t, 3, len(albums))
				assert.Equal(t, "album1", albums[0].ID)
				assert.Equal(t, "album2", albums[1].ID)
				assert.Equal(t, "album3", albums[2].ID)
			},
		},
		{
			name:   "successful album listing with pagination",
			fromID: "album1",
			limit:  5,
			mockSetup: func(m *MockAlbumRepository) {
				m.On("ListAllAlbums", "album1", 5).Return([]*domain.Album{
					{ID: "album2", Name: "Album Two", Description: "Second album"},
					{ID: "album3", Name: "Album Three", Description: "Third album"},
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, albums []*domain.Album, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, albums)
				assert.Equal(t, 2, len(albums))
				assert.Equal(t, "album2", albums[0].ID)
			},
		},
		{
			name:   "empty album list",
			fromID: "",
			limit:  10,
			mockSetup: func(m *MockAlbumRepository) {
				m.On("ListAllAlbums", "", 10).Return(nil, domain.ErrDataNotFound)
			},
			expectedError: domain.ErrDataNotFound,
			validate: func(t *testing.T, albums []*domain.Album, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrDataNotFound, err)
				assert.Nil(t, albums)
			},
		},
		{
			name:   "repository error during listing",
			fromID: "",
			limit:  10,
			mockSetup: func(m *MockAlbumRepository) {
				m.On("ListAllAlbums", "", 10).Return(nil, domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, albums []*domain.Album, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Nil(t, albums)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlbumRepository)
			tt.mockSetup(mockRepo)

			service := NewAlbumService(mockRepo, appconfig.AppConfig{})
			result, err := service.ListAlbums(tt.fromID, tt.limit)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestAlbumCount tests counting total albums
func TestAlbumCount(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockAlbumRepository)
		expectedError error
		validate      func(*testing.T, int64, error)
	}{
		{
			name: "successful album count",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("AlbumCount").Return(int64(42), nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(42), count)
			},
		},
		{
			name: "zero albums",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("AlbumCount").Return(int64(0), nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "repository error during count",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("AlbumCount").Return(int64(0), domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Equal(t, int64(0), count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlbumRepository)
			tt.mockSetup(mockRepo)

			service := NewAlbumService(mockRepo, appconfig.AppConfig{})
			result, err := service.AlbumCount()

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestCreateAlbum tests creating a new album
func TestCreateAlbum(t *testing.T) {
	tests := []struct {
		name          string
		albumName     string
		mockSetup     func(*MockAlbumRepository)
		expectedError error
		validate      func(*testing.T, *domain.Album, error)
	}{
		{
			name:      "successful album creation",
			albumName: "New Album",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("CreateAlbum", "New Album").Return(&domain.Album{
					ID:          "album-new",
					Name:        "New Album",
					Description: "",
				}, nil)
			},
			expectedError: nil,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, album)
				assert.Equal(t, "album-new", album.ID)
				assert.Equal(t, "New Album", album.Name)
			},
		},
		{
			name:      "album already exists",
			albumName: "Existing Album",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("CreateAlbum", "Existing Album").Return(nil, domain.ErrConflictingData)
			},
			expectedError: domain.ErrConflictingData,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrConflictingData, err)
				assert.Nil(t, album)
			},
		},
		{
			name:      "repository error during creation",
			albumName: "Failed Album",
			mockSetup: func(m *MockAlbumRepository) {
				m.On("CreateAlbum", "Failed Album").Return(nil, domain.ErrInternal)
			},
			expectedError: domain.ErrInternal,
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrInternal, err)
				assert.Nil(t, album)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlbumRepository)
			tt.mockSetup(mockRepo)

			service := NewAlbumService(mockRepo, appconfig.AppConfig{})
			result, err := service.CreateAlbum(tt.albumName)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
