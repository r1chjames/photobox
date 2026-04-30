package service

import (
	"errors"
	"testing"

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

func (m *MockPhotoRepository) ListAllPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoRepository) ListAllPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(albumId, fromId, limit, includeThumbnail)
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

func (m *MockPhotoRepository) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
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

func (m *MockFilesystemService) PerformPhotoIndex(callback func(domain.PhotoFile) error) {
	m.Called(callback)
}

func (m *MockFilesystemService) GenerateThumbnail(path string, exifData exif.Exif) []byte {
	args := m.Called(path, exifData)
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
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

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
				m.On("ListAllPhotos", "", 10, false).Return(photos, nil)
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
				m.On("ListAllPhotos", "photo-10", 5, true).Return(photos, nil)
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
				m.On("ListAllPhotos", "", 10, false).Return([]*domain.Photo{}, nil)
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
				m.On("ListAllPhotos", "", 10, false).Return(nil, errors.New("database error"))
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
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

			result, err := service.ListPhotos(tt.fromId, tt.limit, tt.includeThumbnail)

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
				m.On("ListAllPhotosInAlbum", "album-123", "", 10, false).Return(photos, nil)
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
				m.On("ListAllPhotosInAlbum", "album-456", "", 10, false).Return([]*domain.Photo{}, nil)
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
				m.On("ListAllPhotosInAlbum", "album-789", "", 10, false).Return(nil, errors.New("database error"))
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
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

			result, err := service.ListPhotosInAlbum(tt.albumId, tt.fromId, tt.limit, tt.includeThumbnail)

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
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPhotoRepository)
			mockAlbumSvc := new(MockAlbumService)
			mockFsSvc := new(MockFilesystemService)
			config := appconfig.AppConfig{}

			tt.mockSetup(mockRepo)
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

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
				m.On("GetPhotoById", "photo-123", false).Return(&domain.Photo{
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
				m.On("GetPhotoById", "nonexistent", false).Return(nil, errors.New("not found"))
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
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

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
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

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
			service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)

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

		// Mock expects the callback function to be passed
		mockFsSvc.On("PerformPhotoIndex", mock.AnythingOfType("func(domain.PhotoFile) error")).Return()

		service := NewPhotoService(mockRepo, mockAlbumSvc, mockFsSvc, config)
		service.PerformPhotoIndex()

		mockFsSvc.AssertExpectations(t)
	})
}