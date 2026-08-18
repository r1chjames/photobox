package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/datatypes"
)

// MockSearchPhotoService is a mock implementation of port.PhotoService for search tests
type MockSearchPhotoService struct {
	mock.Mock
}

func (m *MockSearchPhotoService) GetPhoto(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	args := m.Called(photoId, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) ListPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) ListPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	args := m.Called(albumId, fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoCount(albumId string) (int64, error) {
	args := m.Called(albumId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoBinary(photoId string) (string, error) {
	args := m.Called(photoId)
	return args.String(0), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoThumbnail(photoId string) ([]byte, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoThumbnailBytes(photoId string) ([]byte, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoThumbnailPath(photoId string) (string, error) {
	args := m.Called(photoId)
	return args.String(0), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoThumbnails(photoIds []string) (map[string][]byte, error) {
	args := m.Called(photoIds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]byte), args.Error(1)
}

func (m *MockSearchPhotoService) SavePhoto(photo domain.PhotoFile) error {
	args := m.Called(photo)
	return args.Error(0)
}

func (m *MockSearchPhotoService) SavePhotos(photos []domain.PhotoFile) error {
	args := m.Called(photos)
	return args.Error(0)
}

func (m *MockSearchPhotoService) PerformPhotoIndex(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockSearchPhotoService) DeletePhoto(photoId string) error {
	args := m.Called(photoId)
	return args.Error(0)
}

func (m *MockSearchPhotoService) RestorePhoto(photoId string) error {
	args := m.Called(photoId)
	return args.Error(0)
}

func (m *MockSearchPhotoService) ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) EmptyTrash() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockSearchPhotoService) PurgeExpiredTrash(cutoff time.Time) (int, error) {
	args := m.Called(cutoff)
	return args.Int(0), args.Error(1)
}

func (m *MockSearchPhotoService) SetFavorite(photoId string, favorite bool) (*domain.Photo, error) {
	args := m.Called(photoId, favorite)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) UpdatePhotoMetadata(photoId string, description string, latitude, longitude *float64, dateTaken *string) (*domain.Photo, error) {
	args := m.Called(photoId, description, latitude, longitude, dateTaken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) ReverseGeocode(ctx context.Context, photoId string) (string, error) {
	args := m.Called(ctx, photoId)
	return args.String(0), args.Error(1)
}
func (m *MockSearchPhotoService) BatchSetFavorite(photoIds []string, favorite bool) error {
	args := m.Called(photoIds, favorite)
	return args.Error(0)
}

func (m *MockSearchPhotoService) BatchAddToAlbum(photoIds []string, albumId string) error {
	args := m.Called(photoIds, albumId)
	return args.Error(0)
}

func (m *MockSearchPhotoService) BatchDeletePhotos(photoIds []string) error {
	args := m.Called(photoIds)
	return args.Error(0)
}

func (m *MockSearchPhotoService) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) ListLowQualityPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) ListMemories(month, day int, maxPerYear int) ([]domain.MemoryGroup, error) {
	args := m.Called(month, day, maxPerYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MemoryGroup), args.Error(1)
}

func (m *MockSearchPhotoService) SearchPhotosWithFilters(filters domain.PhotoSearchFilters, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(filters, fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) ScorePhotoQuality(photoId string) (bool, error) {
	args := m.Called(photoId)
	return args.Bool(0), args.Error(1)
}

func (m *MockSearchPhotoService) ScoreAllPhotoQuality(batchSize int) (int, error) {
	args := m.Called(batchSize)
	return args.Int(0), args.Error(1)
}

func (m *MockSearchPhotoService) Search(query string, limit int) ([]*domain.Photo, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) GetTimeline() ([]domain.TimelineEntry, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.TimelineEntry), args.Error(1)
}

func (m *MockSearchPhotoService) GetGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error) {
	args := m.Called(north, south, east, west)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PhotoGeoData), args.Error(1)
}

func (m *MockSearchPhotoService) RotatePhoto(photoId string, direction string) (*domain.Photo, error) {
	args := m.Called(photoId, direction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) EditPhoto(photoId string, params domain.EditParams) (*domain.Photo, error) {
	args := m.Called(photoId, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) ClearEdits(photoId string) (*domain.Photo, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) DownloadPhotos(photoIds []string, writer io.Writer) error {
	args := m.Called(photoIds, writer)
	return args.Error(0)
}

func (m *MockSearchPhotoService) GetAllTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockSearchPhotoService) ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(tags, fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) UpdatePhotoTags(photoId string, tags []string) (*domain.Photo, error) {
	args := m.Called(photoId, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) BatchUpdatePhotoTags(photoIds []string, tags []string, operation string) error {
	args := m.Called(photoIds, tags, operation)
	return args.Error(0)
}

func (m *MockSearchPhotoService) GetDuplicatePhotos() ([]*domain.Photo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockSearchPhotoService) GenerateThumbnailForPhoto(photoId string) (string, error) {
	args := m.Called(photoId)
	return args.String(0), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoThumbnailPathForSize(photoId string, size string) (string, error) {
	args := m.Called(photoId, size)
	return args.String(0), args.Error(1)
}

func (m *MockSearchPhotoService) PhotoThumbnailBytesForSize(photoId string, size string) ([]byte, error) {
	return nil, nil
}

func (m *MockSearchPhotoService) AnalyzeExistingPhotos() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSearchPhotoService) TriggerAIAnalysis() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSearchPhotoService) RegenerateThumbnails(ctx context.Context) {}

// MockSearchAlbumService is a mock implementation of port.AlbumService for search tests
type MockSearchAlbumService struct {
	mock.Mock
}

func (m *MockSearchAlbumService) GetAlbumById(id string) (*domain.Album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) GetAlbumByName(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) ListAlbums(fromId string, pageSize int) ([]*domain.Album, error) {
	args := m.Called(fromId, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) AlbumCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSearchAlbumService) CreateAlbum(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) CreateSmartAlbum(name string, rules domain.SmartAlbumRules) (*domain.Album, error) {
	args := m.Called(name, rules)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) UpdateSmartAlbum(id string, rules domain.SmartAlbumRules) (*domain.Album, error) {
	args := m.Called(id, rules)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) CreateAlbumIfNotExists(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) UpdateAlbum(id string, updates map[string]any) (*domain.Album, error) {
	args := m.Called(id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockSearchAlbumService) DeleteAlbum(id string, deletePhotos bool) error {
	args := m.Called(id, deletePhotos)
	return args.Error(0)
}

func (m *MockSearchAlbumService) SearchAlbums(query string, limit int) ([]*domain.Album, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Album), args.Error(1)
}

// TestSearchHandler_Search_Success tests successful search
func TestSearchHandler_Search_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockSearchPhotoService)
	mockAlbumSvc := new(MockSearchAlbumService)
	handler := NewSearchHandler(mockPhotoSvc, mockAlbumSvc)

	mockPhotoSvc.On("Search", "test", 30).Return([]*domain.Photo{
		{ID: "p1", Name: "Photo 1"},
	}, nil)

	mockAlbumSvc.On("ListAlbums", "", 30).Return([]*domain.Album{
		{ID: "a1", Name: "Album 1", Metadata: datatypes.JSON([]byte("{}"))},
	}, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/search?q=test", nil)

	handler.Search(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool           `json:"success"`
		Data    []searchResult `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, "p1", response.Data[0].ID)
	assert.Equal(t, "photo", response.Data[0].Type)
	assert.Equal(t, "a1", response.Data[1].ID)
	assert.Equal(t, "album", response.Data[1].Type)

	mockPhotoSvc.AssertExpectations(t)
	mockAlbumSvc.AssertExpectations(t)
}

// TestSearchHandler_Search_MissingQuery tests search without query parameter
func TestSearchHandler_Search_MissingQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockSearchPhotoService)
	mockAlbumSvc := new(MockSearchAlbumService)
	handler := NewSearchHandler(mockPhotoSvc, mockAlbumSvc)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/search", nil)

	handler.Search(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "q query parameter is required", response["error"])

	mockPhotoSvc.AssertNotCalled(t, "Search", mock.Anything, mock.Anything)
	mockAlbumSvc.AssertNotCalled(t, "ListAlbums", mock.Anything, mock.Anything)
}
