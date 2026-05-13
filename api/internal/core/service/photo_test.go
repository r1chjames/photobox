package service

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockPhotoRepository is a mock implementation of port.PhotoRepository
type MockPhotoRepository struct {
	mock.Mock
}

func (m *MockPhotoRepository) GetPhotoById(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	args := m.Called(photoId, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) ListAllPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	args := m.Called(albumId, fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) GetPhotosInAlbumCount(albumId string) (int64, error) {
	args := m.Called(albumId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPhotoRepository) CreatePhotoInfo(photo domain.Photo) error {
	args := m.Called(photo)
	return args.Error(0)
}

func (m *MockPhotoRepository) CreatePhotosInfo(photos []domain.Photo) error {
	args := m.Called(photos)
	return args.Error(0)
}

func (m *MockPhotoRepository) SoftDeletePhoto(photoId string) (*domain.Photo, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) RestorePhoto(photoId string) (*domain.Photo, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) EmptyTrash() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPhotoRepository) UpdatePhoto(photo domain.Photo) error {
	args := m.Called(photo)
	return args.Error(0)
}

func (m *MockPhotoRepository) SetFavorite(photoId string, favorite bool) error {
	args := m.Called(photoId, favorite)
	return args.Error(0)
}

func (m *MockPhotoRepository) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) SearchPhotos(query string, limit int) ([]*domain.Photo, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) GetTimeline() ([]domain.TimelineEntry, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.TimelineEntry), args.Error(1)
}

func (m *MockPhotoRepository) GetPhotosWithGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error) {
	args := m.Called(north, south, east, west)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PhotoGeoData), args.Error(1)
}

func (m *MockPhotoRepository) GetPhotoThumbnails(photoIds []string) (map[string][]byte, error) {
	args := m.Called(photoIds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]byte), args.Error(1)
}

func (m *MockPhotoRepository) GetAllTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPhotoRepository) ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(tags, fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) UpdatePhotoTags(photoId string, tags string) error {
	args := m.Called(photoId, tags)
	return args.Error(0)
}

func (m *MockPhotoRepository) GetDuplicatePhotos() ([]*domain.Photo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) GetThumbnailBytes(photoId string) ([]byte, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPhotoRepository) GetThumbnailPath(photoId string) (string, error) {
	args := m.Called(photoId)
	return args.String(0), args.Error(1)
}

func (m *MockPhotoRepository) GetPhotoIndexCache() (map[string]struct{ FileHash string; FileModifiedTime int64 }, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]struct{ FileHash string; FileModifiedTime int64 }), args.Error(1)
}

func (m *MockPhotoRepository) AddAITags(photoId string, tags []string) error {
	args := m.Called(photoId, tags)
	return args.Error(0)
}

func (m *MockPhotoRepository) ListPhotosWithoutAITags(limit int) ([]*domain.Photo, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

// MockAlbumService is a mock implementation of port.AlbumService
type MockAlbumService struct {
	mock.Mock
}

func (m *MockAlbumService) GetAlbumByName(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) CreateAlbum(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) GetAlbum(albumId string) (*domain.Album, error) {
	args := m.Called(albumId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) ListAlbums(fromId string, limit int) ([]*domain.Album, error) {
	args := m.Called(fromId, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Album), args.Error(1)
}

func (m *MockAlbumService) AlbumCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAlbumService) GetAlbumById(id string) (*domain.Album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) UpdateAlbum(id string, updates map[string]any) (*domain.Album, error) {
	args := m.Called(id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) DeleteAlbum(id string, deletePhotos bool) error {
	args := m.Called(id, deletePhotos)
	return args.Error(0)
}

func (m *MockAlbumService) SearchAlbums(query string, limit int) ([]*domain.Album, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Album), args.Error(1)
}

// MockFilesystemService is a mock implementation of port.FilesystemService
type MockFilesystemService struct {
	mock.Mock
}

func (m *MockFilesystemService) PerformPhotoIndex(callback func([]domain.PhotoFile) error, indexCache map[string]struct{ FileHash string; FileModifiedTime int64 }) {
	m.Called(callback, indexCache)
}

func (m *MockFilesystemService) GenerateThumbnail(path string, exifData exif.Exif, width, height int) []byte {
	args := m.Called(path, exifData, width, height)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]byte)
}

func (m *MockFilesystemService) WriteFileToFilesystem(photo domain.PhotoUpload) domain.PhotoFile {
	args := m.Called(photo)
	return args.Get(0).(domain.PhotoFile)
}

func (m *MockFilesystemService) MoveToTrash(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func (m *MockFilesystemService) RestoreFromTrash(trashPath, originalPath string) error {
	args := m.Called(trashPath, originalPath)
	return args.Error(0)
}

func (m *MockFilesystemService) RenameDirectory(oldPath, newPath string) error {
	args := m.Called(oldPath, newPath)
	return args.Error(0)
}

// MockCacheService is a mock implementation of port.CacheService
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) Set(key, value string, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheService) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCacheService) DeletePattern(pattern string) error {
	args := m.Called(pattern)
	return args.Error(0)
}

func (m *MockCacheService) GetBytes(key string) ([]byte, error) {
	return nil, nil
}

func (m *MockCacheService) SetBytes(key string, value []byte, ttl time.Duration) error {
	return nil
}

// TestGetPhoto tests retrieving a photo by ID
func TestGetPhoto(t *testing.T) {
	tests := []struct {
		name             string
		photoId          string
		includeThumbnail bool
		mockSetup        func(*MockPhotoRepository)
		validate         func(*testing.T, *domain.Photo, error)
	}{
		{
			name:             "successfully get photo with thumbnail",
			photoId:          "photo-123",
			includeThumbnail: true,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "photo-123", true).Return(&domain.Photo{
					ID:             "photo-123",
					Name:           "test.jpg",
					FilesystemPath: "/path/to/photo.jpg",
					Thumbnail:      []byte("thumbnail-data"),
				}, nil)
			},
			validate: func(t *testing.T, photo *domain.Photo, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, photo)
				assert.Equal(t, "photo-123", photo.ID)
				assert.Equal(t, "test.jpg", photo.Name)
				assert.Equal(t, "photo/photo-123/bin", photo.SourcePath)
				assert.NotEmpty(t, photo.Thumbnail)
			},
		},
		{
			name:             "successfully get photo without thumbnail",
			photoId:          "photo-456",
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "photo-456", false).Return(&domain.Photo{
					ID:             "photo-456",
					Name:           "photo2.jpg",
					FilesystemPath: "/path/to/photo2.jpg",
				}, nil)
			},
			validate: func(t *testing.T, photo *domain.Photo, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, photo)
				assert.Equal(t, "photo-456", photo.ID)
				assert.Equal(t, "photo/photo-456/bin", photo.SourcePath)
			},
		},
		{
			name:             "photo not found",
			photoId:          "nonexistent",
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "nonexistent", false).Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, photo *domain.Photo, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrDataNotFound))
				assert.Nil(t, photo)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.GetPhoto(tt.photoId, tt.includeThumbnail)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListPhotos tests listing all photos with pagination
func TestListPhotos(t *testing.T) {
	tests := []struct {
		name             string
		fromId           string
		limit            int
		includeThumbnail bool
		mockSetup        func(*MockPhotoRepository)
		validate         func(*testing.T, []*domain.Photo, error)
	}{
		{
			name:             "successfully list photos",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				photos := []*domain.Photo{
					{ID: "photo-1", Name: "photo1.jpg"},
					{ID: "photo-2", Name: "photo2.jpg"},
					{ID: "photo-3", Name: "photo3.jpg"},
				}
				m.On("ListAllPhotos", "", 10, false, "", "", "").Return(photos, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Len(t, photos, 3)
				assert.Equal(t, "photo-1", photos[0].ID)
				assert.Equal(t, "photo/photo-1/bin", photos[0].SourcePath)
				assert.Equal(t, "photo-2", photos[1].ID)
				assert.Equal(t, "photo/photo-2/bin", photos[1].SourcePath)
			},
		},
		{
			name:             "successfully list photos with thumbnail",
			fromId:           "photo-10",
			limit:            5,
			includeThumbnail: true,
			mockSetup: func(m *MockPhotoRepository) {
				photos := []*domain.Photo{
					{ID: "photo-11", Name: "photo11.jpg", Thumbnail: []byte("thumb1")},
					{ID: "photo-12", Name: "photo12.jpg", Thumbnail: []byte("thumb2")},
				}
				m.On("ListAllPhotos", "photo-10", 5, true, "", "", "").Return(photos, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Len(t, photos, 2)
				assert.NotEmpty(t, photos[0].Thumbnail)
			},
		},
		{
			name:             "empty photo list",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("ListAllPhotos", "", 10, false, "", "", "").Return([]*domain.Photo{}, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Empty(t, photos)
			},
		},
		{
			name:             "repository error",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("ListAllPhotos", "", 10, false, "", "", "").Return(nil, errors.New("database error"))
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrDataNotFound))
				assert.Nil(t, photos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.ListPhotos(tt.fromId, tt.limit, tt.includeThumbnail, "", "", "")

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListPhotosInAlbum tests listing photos in a specific album
func TestListPhotosInAlbum(t *testing.T) {
	tests := []struct {
		name             string
		albumId          string
		fromId           string
		limit            int
		includeThumbnail bool
		mockSetup        func(*MockPhotoRepository)
		validate         func(*testing.T, []*domain.Photo, error)
	}{
		{
			name:             "successfully list photos in album",
			albumId:          "album-123",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				photos := []*domain.Photo{
					{ID: "photo-1", Name: "photo1.jpg", AlbumId: "album-123"},
					{ID: "photo-2", Name: "photo2.jpg", AlbumId: "album-123"},
				}
				m.On("ListAllPhotosInAlbum", "album-123", "", 10, false, "", "", "").Return(photos, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Len(t, photos, 2)
				assert.Equal(t, "album-123", photos[0].AlbumId)
				assert.Equal(t, "photo/photo-1/bin", photos[0].SourcePath)
			},
		},
		{
			name:             "empty album",
			albumId:          "album-456",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("ListAllPhotosInAlbum", "album-456", "", 10, false, "", "", "").Return([]*domain.Photo{}, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Empty(t, photos)
			},
		},
		{
			name:             "repository error",
			albumId:          "album-789",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("ListAllPhotosInAlbum", "album-789", "", 10, false, "", "", "").Return(nil, errors.New("database error"))
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrDataNotFound))
				assert.Nil(t, photos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.ListPhotosInAlbum(tt.albumId, tt.fromId, tt.limit, tt.includeThumbnail, "", "", "")

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPhotoCount tests getting count of photos in an album
func TestPhotoCount(t *testing.T) {
	tests := []struct {
		name      string
		albumId   string
		mockSetup func(*MockPhotoRepository)
		validate  func(*testing.T, int64, error)
	}{
		{
			name:    "successfully get photo count",
			albumId: "album-123",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotosInAlbumCount", "album-123").Return(int64(42), nil)
			},
			validate: func(t *testing.T, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(42), count)
			},
		},
		{
			name:    "zero photos in album",
			albumId: "album-456",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotosInAlbumCount", "album-456").Return(int64(0), nil)
			},
			validate: func(t *testing.T, count int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name:    "repository error",
			albumId: "album-789",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotosInAlbumCount", "album-789").Return(int64(0), errors.New("database error"))
			},
			validate: func(t *testing.T, count int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.PhotoCount(tt.albumId)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPhotoBinary tests getting photo binary path
func TestPhotoBinary(t *testing.T) {
	tests := []struct {
		name      string
		photoId   string
		mockSetup func(*MockPhotoRepository)
		validate  func(*testing.T, string, error)
	}{
		{
			name:    "successfully get photo binary path",
			photoId: "photo-123",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "photo-123", false).Return(&domain.Photo{
					ID:             "photo-123",
					FilesystemPath: "/storage/photos/test.jpg",
				}, nil)
			},
			validate: func(t *testing.T, path string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "/storage/photos/test.jpg", path)
			},
		},
		{
			name:    "photo not found",
			photoId: "nonexistent",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "nonexistent", false).Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, path string, err error) {
				assert.Error(t, err)
				assert.Empty(t, path)
			},
		},
		{
			name:    "path traversal detected",
			photoId: "photo-evil",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "photo-evil", false).Return(&domain.Photo{
					ID:             "photo-evil",
					FilesystemPath: "/etc/passwd",
				}, nil)
			},
			validate: func(t *testing.T, path string, err error) {
				assert.ErrorIs(t, err, domain.ErrForbidden)
				assert.Empty(t, path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{PhotoDir: "/storage/photos"}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.PhotoBinary(tt.photoId)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPhotoThumbnail tests getting photo thumbnail
func TestPhotoThumbnail(t *testing.T) {
	tests := []struct {
		name      string
		photoId   string
		mockSetup func(*MockPhotoRepository)
		validate  func(*testing.T, []byte, error)
	}{
		{
			name:    "successfully get photo thumbnail",
			photoId: "photo-123",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "photo-123", true).Return(&domain.Photo{
					ID:        "photo-123",
					Thumbnail: []byte("thumbnail-data"),
				}, nil)
			},
			validate: func(t *testing.T, thumbnail []byte, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []byte("thumbnail-data"), thumbnail)
			},
		},
		{
			name:    "photo not found",
			photoId: "nonexistent",
			mockSetup: func(m *MockPhotoRepository) {
				m.On("GetPhotoById", "nonexistent", true).Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, thumbnail []byte, err error) {
				assert.Error(t, err)
				assert.Nil(t, thumbnail)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.PhotoThumbnail(tt.photoId)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestSavePhoto tests saving a single photo
func TestSavePhoto(t *testing.T) {
	tests := []struct {
		name      string
		photoFile domain.PhotoFile
		mockSetup func(*MockPhotoRepository, *MockAlbumService)
		validate  func(*testing.T, error)
	}{
		{
			name: "successfully save photo with existing album",
			photoFile: domain.PhotoFile{
				Name:      "test.jpg",
				Path:      "/photos/vacation/test.jpg",
				Directory: "vacation",
				Thumbnail: []byte("thumb"),
			},
			mockSetup: func(mPhoto *MockPhotoRepository, mAlbum *MockAlbumService) {
				mAlbum.On("GetAlbumByName", "vacation").Return(&domain.Album{
					ID:   "album-123",
					Name: "vacation",
				}, nil)
				mPhoto.On("CreatePhotoInfo", mock.AnythingOfType("domain.Photo")).Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "successfully save photo with new album creation",
			photoFile: domain.PhotoFile{
				Name:      "new.jpg",
				Path:      "/photos/newalbum/new.jpg",
				Directory: "newalbum",
			},
			mockSetup: func(mPhoto *MockPhotoRepository, mAlbum *MockAlbumService) {
				mAlbum.On("GetAlbumByName", "newalbum").Return(nil, errors.New("not found"))
				mAlbum.On("CreateAlbum", "newalbum").Return(&domain.Album{
					ID:   "album-456",
					Name: "newalbum",
				}, nil)
				mPhoto.On("CreatePhotoInfo", mock.AnythingOfType("domain.Photo")).Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "save photo when album creation fails",
			photoFile: domain.PhotoFile{
				Name:      "test.jpg",
				Path:      "/photos/failalbum/test.jpg",
				Directory: "failalbum",
			},
			mockSetup: func(mPhoto *MockPhotoRepository, mAlbum *MockAlbumService) {
				mAlbum.On("GetAlbumByName", "failalbum").Return(nil, errors.New("not found"))
				mAlbum.On("CreateAlbum", "failalbum").Return(nil, errors.New("album creation failed"))
				// Service should still try to create photo with empty albumId
				mPhoto.On("CreatePhotoInfo", mock.AnythingOfType("domain.Photo")).Return(nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "fail to save photo info",
			photoFile: domain.PhotoFile{
				Name:      "test.jpg",
				Path:      "/photos/vacation/test.jpg",
				Directory: "vacation",
			},
			mockSetup: func(mPhoto *MockPhotoRepository, mAlbum *MockAlbumService) {
				mAlbum.On("GetAlbumByName", "vacation").Return(&domain.Album{
					ID:   "album-123",
					Name: "vacation",
				}, nil)
				mPhoto.On("CreatePhotoInfo", mock.AnythingOfType("domain.Photo")).Return(errors.New("database error"))
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo, mockAlbumSvc)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			err := service.SavePhoto(tt.photoFile)

			tt.validate(t, err)
			mockRepo.AssertExpectations(t)
			mockAlbumSvc.AssertExpectations(t)
		})
	}
}

// TestSavePhotos tests saving multiple photos
func TestSavePhotos(t *testing.T) {
	tests := []struct {
		name       string
		photoFiles []domain.PhotoFile
		mockSetup  func(*MockPhotoRepository, *MockAlbumService)
		validate   func(*testing.T, error)
	}{
		{
			name: "successfully save multiple photos",
			photoFiles: []domain.PhotoFile{
				{Name: "photo1.jpg", Path: "/photos/album1/photo1.jpg", Directory: "album1"},
				{Name: "photo2.jpg", Path: "/photos/album1/photo2.jpg", Directory: "album1"},
				{Name: "photo3.jpg", Path: "/photos/album2/photo3.jpg", Directory: "album2"},
			},
			mockSetup: func(mPhoto *MockPhotoRepository, mAlbum *MockAlbumService) {
				// Album lookups are cached, so album1 is only looked up once
				mAlbum.On("GetAlbumByName", "album1").Return(&domain.Album{
					ID:   "album-1",
					Name: "album1",
				}, nil).Once()
				// album2 is looked up once
				mAlbum.On("GetAlbumByName", "album2").Return(&domain.Album{
					ID:   "album-2",
					Name: "album2",
				}, nil).Once()
				// All photos get saved in a single batch
				mPhoto.On("CreatePhotosInfo", mock.AnythingOfType("[]domain.Photo")).Return(nil).Once()
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "fail to save batch",
			photoFiles: []domain.PhotoFile{
				{Name: "photo1.jpg", Path: "/photos/album1/photo1.jpg", Directory: "album1"},
				{Name: "photo2.jpg", Path: "/photos/album1/photo2.jpg", Directory: "album1"},
			},
			mockSetup: func(mPhoto *MockPhotoRepository, mAlbum *MockAlbumService) {
				// Album lookup cached, called once
				mAlbum.On("GetAlbumByName", "album1").Return(&domain.Album{
					ID:   "album-1",
					Name: "album1",
				}, nil).Once()
				// Batch save fails
				mPhoto.On("CreatePhotosInfo", mock.AnythingOfType("[]domain.Photo")).Return(errors.New("save error")).Once()
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "save error")
			},
		},
		{
			name:       "empty photo list",
			photoFiles: []domain.PhotoFile{},
			mockSetup:  func(mPhoto *MockPhotoRepository, mAlbum *MockAlbumService) {},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo, mockAlbumSvc)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			err := service.SavePhotos(tt.photoFiles)

			tt.validate(t, err)
			mockRepo.AssertExpectations(t)
			mockAlbumSvc.AssertExpectations(t)
		})
	}
}

// TestPerformPhotoIndex tests photo indexing
func TestPerformPhotoIndex(t *testing.T) {
	t.Run("calls filesystem service to perform index", func(t *testing.T) {
		mockRepo := new(MockPhotoRepository)
		mockAlbumSvc := new(MockAlbumService)
		mockFsSvc := new(MockFilesystemService)
		config := appconfig.AppConfig{}

		// Mock GetPhotoIndexCache to return empty cache
		mockRepo.On("GetPhotoIndexCache").Return(map[string]struct{ FileHash string; FileModifiedTime int64 }{}, nil)
		// Mock expects the callback function and indexCache to be passed
		mockFsSvc.On("PerformPhotoIndex", mock.AnythingOfType("func([]domain.PhotoFile) error"), mock.Anything).Return()

		service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)
		service.PerformPhotoIndex()

		mockRepo.AssertExpectations(t)
		mockFsSvc.AssertExpectations(t)
	})
}

// TestDeletePhoto tests deleting a photo
func TestDeletePhoto(t *testing.T) {
	tests := []struct {
		name      string
		photoId   string
		mockSetup func(*MockPhotoRepository, *MockFilesystemService)
		validate  func(*testing.T, error)
	}{
		{
			name:    "successfully delete photo",
			photoId: "photo-123",
			mockSetup: func(mPhoto *MockPhotoRepository, mFs *MockFilesystemService) {
				mPhoto.On("SoftDeletePhoto", "photo-123").Return(&domain.Photo{
					ID:             "photo-123",
					FilesystemPath: "/storage/photos/test.jpg",
				}, nil)
				mFs.On("MoveToTrash", mock.AnythingOfType("string")).Return("/trash/test.jpg", nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "photo not found",
			photoId: "nonexistent",
			mockSetup: func(mPhoto *MockPhotoRepository, mFs *MockFilesystemService) {
				mPhoto.On("SoftDeletePhoto", "nonexistent").Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo, mockFsSvc)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			err := service.DeletePhoto(tt.photoId)

			tt.validate(t, err)
			mockRepo.AssertExpectations(t)
			mockFsSvc.AssertExpectations(t)
		})
	}
}

// TestRestorePhoto tests restoring a photo from trash
func TestRestorePhoto(t *testing.T) {
	tests := []struct {
		name      string
		photoId   string
		mockSetup func(*MockPhotoRepository)
		validate  func(*testing.T, error)
	}{
		{
			name:    "successfully restore photo",
			photoId: "photo-123",
			mockSetup: func(mPhoto *MockPhotoRepository) {
				mPhoto.On("RestorePhoto", "photo-123").Return(&domain.Photo{
					ID:             "photo-123",
					FilesystemPath: "/storage/photos/test.jpg",
				}, nil)
			},
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "photo not found",
			photoId: "nonexistent",
			mockSetup: func(mPhoto *MockPhotoRepository) {
				mPhoto.On("RestorePhoto", "nonexistent").Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			err := service.RestorePhoto(tt.photoId)

			tt.validate(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListTrashPhotos tests listing trashed photos
func TestListTrashPhotos(t *testing.T) {
	tests := []struct {
		name             string
		fromId           string
		limit            int
		includeThumbnail bool
		mockSetup        func(*MockPhotoRepository)
		validate         func(*testing.T, []*domain.Photo, error)
	}{
		{
			name:             "successfully list trash photos",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				photos := []*domain.Photo{
					{ID: "photo-1", Name: "photo1.jpg"},
					{ID: "photo-2", Name: "photo2.jpg"},
				}
				m.On("ListTrashPhotos", "", 10, false).Return(photos, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Len(t, photos, 2)
				assert.Equal(t, "photo/photo-1/bin", photos[0].SourcePath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.ListTrashPhotos(tt.fromId, tt.limit, tt.includeThumbnail)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestEmptyTrash tests emptying the trash
func TestEmptyTrash(t *testing.T) {
	t.Run("successfully empty trash", func(t *testing.T) {
		mockRepo := new(MockPhotoRepository)
		mockAlbumSvc := new(MockAlbumService)
		mockFsSvc := new(MockFilesystemService)
		config := appconfig.AppConfig{}

		mockRepo.On("EmptyTrash").Return(nil)

		service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)
		err := service.EmptyTrash()

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestSetFavorite tests setting a photo as favorite
func TestSetFavorite(t *testing.T) {
	tests := []struct {
		name      string
		photoId   string
		favorite  bool
		mockSetup func(*MockPhotoRepository)
		validate  func(*testing.T, *domain.Photo, error)
	}{
		{
			name:     "successfully set favorite",
			photoId:  "photo-123",
			favorite: true,
			mockSetup: func(mPhoto *MockPhotoRepository) {
				mPhoto.On("GetPhotoById", "photo-123", false).Return(&domain.Photo{
					ID:       "photo-123",
					Name:     "test.jpg",
					Favorite: false,
				}, nil)
				mPhoto.On("SetFavorite", "photo-123", true).Return(nil)
			},
			validate: func(t *testing.T, photo *domain.Photo, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, photo)
				assert.True(t, photo.Favorite)
				assert.Equal(t, "photo/photo-123/bin", photo.SourcePath)
			},
		},
		{
			name:     "photo not found",
			photoId:  "nonexistent",
			favorite: true,
			mockSetup: func(mPhoto *MockPhotoRepository) {
				mPhoto.On("GetPhotoById", "nonexistent", false).Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, photo *domain.Photo, err error) {
				assert.Error(t, err)
				assert.Nil(t, photo)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.SetFavorite(tt.photoId, tt.favorite)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestListFavoritePhotos tests listing favorite photos
func TestListFavoritePhotos(t *testing.T) {
	tests := []struct {
		name             string
		fromId           string
		limit            int
		includeThumbnail bool
		mockSetup        func(*MockPhotoRepository)
		validate         func(*testing.T, []*domain.Photo, error)
	}{
		{
			name:             "successfully list favorite photos",
			fromId:           "",
			limit:            10,
			includeThumbnail: false,
			mockSetup: func(m *MockPhotoRepository) {
				photos := []*domain.Photo{
					{ID: "photo-1", Name: "photo1.jpg", Favorite: true},
					{ID: "photo-2", Name: "photo2.jpg", Favorite: true},
				}
				m.On("ListFavoritePhotos", "", 10, false).Return(photos, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Len(t, photos, 2)
				assert.Equal(t, "photo/photo-1/bin", photos[0].SourcePath)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.ListFavoritePhotos(tt.fromId, tt.limit, tt.includeThumbnail, "", "")

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestSearch tests searching photos
func TestSearch(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		limit     int
		mockSetup func(*MockPhotoRepository)
		validate  func(*testing.T, []*domain.Photo, error)
	}{
		{
			name:  "successfully search photos",
			query: "vacation",
			limit: 10,
			mockSetup: func(m *MockPhotoRepository) {
				photos := []*domain.Photo{
					{ID: "photo-1", Name: "vacation1.jpg"},
					{ID: "photo-2", Name: "vacation2.jpg"},
				}
				m.On("SearchPhotos", "vacation", 10).Return(photos, nil)
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.NoError(t, err)
				assert.Len(t, photos, 2)
			},
		},
		{
			name:  "search returns error",
			query: "nonexistent",
			limit: 10,
			mockSetup: func(m *MockPhotoRepository) {
				m.On("SearchPhotos", "nonexistent", 10).Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, photos []*domain.Photo, err error) {
				assert.Error(t, err)
				assert.Nil(t, photos)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.Search(tt.query, tt.limit)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetTimeline tests retrieving the photo timeline
func TestGetTimeline(t *testing.T) {
	t.Run("successfully get timeline", func(t *testing.T) {
		mockRepo := new(MockPhotoRepository)
		mockAlbumSvc := new(MockAlbumService)
		mockFsSvc := new(MockFilesystemService)
		config := appconfig.AppConfig{}

		entries := []domain.TimelineEntry{
			{Year: 2024, Month: 1, Count: 10},
			{Year: 2024, Month: 2, Count: 5},
		}
		mockRepo.On("GetTimeline").Return(entries, nil)

		service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)
		result, err := service.GetTimeline()

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2024, result[0].Year)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetGeodata tests retrieving photos with geodata
func TestGetGeodata(t *testing.T) {
	t.Run("successfully get geodata", func(t *testing.T) {
		mockRepo := new(MockPhotoRepository)
		mockAlbumSvc := new(MockAlbumService)
		mockFsSvc := new(MockFilesystemService)
		config := appconfig.AppConfig{}

		geoData := []domain.PhotoGeoData{
			{ID: "photo-1", Lat: 51.5, Lng: -0.1},
			{ID: "photo-2", Lat: 40.7, Lng: -74.0},
		}
		mockRepo.On("GetPhotosWithGeodata", 52.0, 50.0, 1.0, -1.0).Return(geoData, nil)

		service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)
		result, err := service.GetGeodata(52.0, 50.0, 1.0, -1.0)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "photo-1", result[0].ID)
		mockRepo.AssertExpectations(t)
	})
}

// TestRotatePhoto tests rotating a photo
func TestRotatePhoto(t *testing.T) {
	tests := []struct {
		name      string
		photoId   string
		direction string
		mockSetup func(*MockPhotoRepository)
		validate  func(*testing.T, *domain.Photo, error)
	}{
		{
			name:      "successfully rotate photo",
			photoId:   "photo-123",
			direction: "right",
			mockSetup: func(mPhoto *MockPhotoRepository) {
				mPhoto.On("GetPhotoById", "photo-123", false).Return(&domain.Photo{
					ID:   "photo-123",
					Name: "test.jpg",
				}, nil)
			},
			validate: func(t *testing.T, photo *domain.Photo, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, photo)
				assert.Equal(t, "photo-123", photo.ID)
			},
		},
		{
			name:      "photo not found",
			photoId:   "nonexistent",
			direction: "left",
			mockSetup: func(mPhoto *MockPhotoRepository) {
				mPhoto.On("GetPhotoById", "nonexistent", false).Return(nil, errors.New("not found"))
			},
			validate: func(t *testing.T, photo *domain.Photo, err error) {
				assert.Error(t, err)
				assert.Nil(t, photo)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

			result, err := service.RotatePhoto(tt.photoId, tt.direction)

			tt.validate(t, result, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestDownloadPhotos tests downloading photos as a zip archive
func TestDownloadPhotos(t *testing.T) {
	t.Run("successfully download photos", func(t *testing.T) {
		mockRepo := new(MockPhotoRepository)
		mockAlbumSvc := new(MockAlbumService)
		mockFsSvc := new(MockFilesystemService)
		config := appconfig.AppConfig{}

		// Create a temporary file to act as the photo on disk
		tmpFile, err := os.CreateTemp("", "photo-*.jpg")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("photo-binary-data")
		assert.NoError(t, err)
		err = tmpFile.Close()
		assert.NoError(t, err)

		mockRepo.On("GetPhotoById", "photo-123", false).Return(&domain.Photo{
			ID:             "photo-123",
			Name:           "test.jpg",
			FilesystemPath: tmpFile.Name(),
		}, nil)

		service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, new(MockCacheService), nil, config)

		var buf bytes.Buffer
		err = service.DownloadPhotos([]string{"photo-123"}, &buf)

		assert.NoError(t, err)
		assert.Greater(t, buf.Len(), 0)
		mockRepo.AssertExpectations(t)
	})
}