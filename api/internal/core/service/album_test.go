package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/datatypes"
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

func (m *MockAlbumRepository) CreateAlbumWithMetadata(album *domain.Album) error {
	args := m.Called(album)
	return args.Error(0)
}

func (m *MockAlbumRepository) UpdateAlbum(album *domain.Album) error {
	args := m.Called(album)
	return args.Error(0)
}

func (m *MockAlbumRepository) DeleteAlbum(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAlbumRepository) ReassignPhotosToAlbum(fromAlbumId, toAlbumId string) error {
	args := m.Called(fromAlbumId, toAlbumId)
	return args.Error(0)
}

func (m *MockAlbumRepository) SearchAlbums(query string, limit int) ([]*domain.Album, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Album), args.Error(1)
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

// TestUpdateAlbum tests updating an album
func TestUpdateAlbum(t *testing.T) {
	tests := []struct {
		name      string
		albumID   string
		updates   map[string]any
		mockSetup func(*MockAlbumRepository)
		validate  func(*testing.T, *domain.Album, error)
	}{
		{
			name:    "successfully update album",
			albumID: "album-123",
			updates: map[string]any{
				"name":        "Updated Name",
				"description": "Updated Description",
			},
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumById", "album-123").Return(&domain.Album{
					ID:          "album-123",
					Name:        "Original Name",
					Description: "Original Description",
				}, nil)
				m.On("UpdateAlbum", mock.AnythingOfType("*domain.Album")).Return(nil)
			},
			validate: func(t *testing.T, album *domain.Album, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, album)
				assert.Equal(t, "Updated Name", album.Name)
				assert.Equal(t, "Updated Description", album.Description)
			},
		},
		{
			name:    "album not found",
			albumID: "nonexistent",
			updates: map[string]any{
				"name": "Updated Name",
			},
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumById", "nonexistent").Return(nil, domain.ErrDataNotFound)
			},
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

			service := NewAlbumService(mockRepo, appconfig.AppConfig{PhotoDir: "/tmp/photos"})
			result, err := service.UpdateAlbum(tt.albumID, tt.updates)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestDeleteAlbum tests deleting an album
func TestDeleteAlbum(t *testing.T) {
	tests := []struct {
		name         string
		albumID      string
		deletePhotos bool
		mockSetup    func(*MockAlbumRepository)
		validate     func(*testing.T, error)
	}{
		{
			name:         "successfully delete album without deleting photos",
			albumID:      "album-123",
			deletePhotos: false,
			mockSetup: func(m *MockAlbumRepository) {
				m.On("GetAlbumByName", "Uncategorized").Return(&domain.Album{
					ID:   "album-uncategorized",
					Name: "Uncategorized",
				}, nil)
				m.On("ReassignPhotosToAlbum", "album-123", "album-uncategorized").Return(nil)
				m.On("DeleteAlbum", "album-123").Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "successfully delete album and photos",
			albumID:      "album-123",
			deletePhotos: true,
			mockSetup: func(m *MockAlbumRepository) {
				m.On("DeleteAlbum", "album-123").Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlbumRepository)
			tt.mockSetup(mockRepo)

			service := NewAlbumService(mockRepo, appconfig.AppConfig{})
			err := service.DeleteAlbum(tt.albumID, tt.deletePhotos)

			tt.validate(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestSearchAlbums tests searching albums
func TestSearchAlbums(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		limit     int
		mockSetup func(*MockAlbumRepository)
		validate  func(*testing.T, []*domain.Album, error)
	}{
		{
			name:  "successfully search albums",
			query: "vacation",
			limit: 10,
			mockSetup: func(m *MockAlbumRepository) {
				albums := []*domain.Album{
					{ID: "album-1", Name: "Vacation 2024"},
					{ID: "album-2", Name: "Vacation 2023"},
				}
				m.On("SearchAlbums", "vacation", 10).Return(albums, nil)
			},
			validate: func(t *testing.T, albums []*domain.Album, err error) {
				assert.NoError(t, err)
				assert.Len(t, albums, 2)
			},
		},
		{
			name:  "search returns error",
			query: "nonexistent",
			limit: 10,
			mockSetup: func(m *MockAlbumRepository) {
				m.On("SearchAlbums", "nonexistent", 10).Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, albums []*domain.Album, err error) {
				assert.Error(t, err)
				assert.Nil(t, albums)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlbumRepository)
			tt.mockSetup(mockRepo)

			service := NewAlbumService(mockRepo, appconfig.AppConfig{})
			result, err := service.SearchAlbums(tt.query, tt.limit)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestSmartAlbums tests smart album creation and rule parsing
func TestSmartAlbums(t *testing.T) {
	t.Run("create smart album stores rules in metadata", func(t *testing.T) {
		mockRepo := new(MockAlbumRepository)
		config := appconfig.AppConfig{}

		rules := domain.SmartAlbumRules{Camera: "Canon", Favorite: true}
		mockRepo.On("CreateAlbumWithMetadata", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			album := args.Get(0).(*domain.Album)
			assert.Equal(t, "Best 2024", album.Name)
			assert.True(t, album.IsSmart())
			assert.Equal(t, "Canon", album.SmartRules().Camera)
			assert.True(t, album.SmartRules().Favorite)
		})

		svc := NewAlbumService(mockRepo, config)
		album, err := svc.CreateSmartAlbum("Best 2024", rules)

		assert.NoError(t, err)
		assert.NotNil(t, album)
		assert.True(t, album.IsSmart())
		mockRepo.AssertExpectations(t)
	})

	t.Run("rules convert to photo search filters", func(t *testing.T) {
		rules := domain.SmartAlbumRules{Camera: "Canon", HasGPS: true, Orientation: "landscape"}
		filters := rules.ToFilters()
		assert.Equal(t, "Canon", filters.Camera)
		assert.True(t, filters.HasGPS)
		assert.Equal(t, "landscape", filters.Orientation)
	})

	t.Run("regular album is not smart", func(t *testing.T) {
		album := &domain.Album{Name: "Manual", Metadata: datatypes.JSON([]byte(`{"foo":"bar"}`))}
		assert.False(t, album.IsSmart())
		assert.Equal(t, domain.SmartAlbumRules{}, album.SmartRules())
	})
}
