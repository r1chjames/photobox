package http

import (
	"bytes"
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
)

// MockPhotoService is a mock implementation of port.PhotoService
type MockPhotoService struct {
	mock.Mock
}

func (m *MockPhotoService) GetPhoto(photoId string, includeThumbnail bool) (*domain.Photo, error) {
	args := m.Called(photoId, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) ListPhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) ListPhotosInAlbum(albumId string, fromId string, limit int, includeThumbnail bool, startDate string, endDate string, mediaType string) ([]*domain.Photo, error) {
	args := m.Called(albumId, fromId, limit, includeThumbnail, startDate, endDate, mediaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) PhotoCount(albumId string) (int64, error) {
	args := m.Called(albumId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPhotoService) PhotoBinary(photoId string) (string, error) {
	args := m.Called(photoId)
	return args.String(0), args.Error(1)
}

func (m *MockPhotoService) PhotoThumbnail(photoId string) ([]byte, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPhotoService) PhotoThumbnailBytes(photoId string) ([]byte, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockPhotoService) PhotoThumbnailPath(photoId string) (string, error) {
	args := m.Called(photoId)
	return args.String(0), args.Error(1)
}

func (m *MockPhotoService) PhotoThumbnails(photoIds []string) (map[string][]byte, error) {
	args := m.Called(photoIds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]byte), args.Error(1)
}

func (m *MockPhotoService) SavePhoto(photo domain.PhotoFile) error {
	args := m.Called(photo)
	return args.Error(0)
}

func (m *MockPhotoService) SavePhotos(photos []domain.PhotoFile) error {
	args := m.Called(photos)
	return args.Error(0)
}

func (m *MockPhotoService) PerformPhotoIndex(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockPhotoService) DeletePhoto(photoId string) error {
	args := m.Called(photoId)
	return args.Error(0)
}

func (m *MockPhotoService) RestorePhoto(photoId string) error {
	args := m.Called(photoId)
	return args.Error(0)
}

func (m *MockPhotoService) ListTrashPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) EmptyTrash() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPhotoService) PurgeExpiredTrash(cutoff time.Time) (int, error) {
	args := m.Called(cutoff)
	return args.Int(0), args.Error(1)
}

func (m *MockPhotoService) SetFavorite(photoId string, favorite bool) (*domain.Photo, error) {
	args := m.Called(photoId, favorite)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) UpdatePhotoMetadata(photoId string, description string, latitude, longitude *float64, dateTaken *string) (*domain.Photo, error) {
	args := m.Called(photoId, description, latitude, longitude, dateTaken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) ReverseGeocode(ctx context.Context, photoId string) (string, error) {
	args := m.Called(ctx, photoId)
	return args.String(0), args.Error(1)
}

func (m *MockPhotoService) BatchSetFavorite(photoIds []string, favorite bool) error {
	args := m.Called(photoIds, favorite)
	return args.Error(0)
}

func (m *MockPhotoService) BatchAddToAlbum(photoIds []string, albumId string) error {
	args := m.Called(photoIds, albumId)
	return args.Error(0)
}

func (m *MockPhotoService) BatchDeletePhotos(photoIds []string) error {
	args := m.Called(photoIds)
	return args.Error(0)
}

func (m *MockPhotoService) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) ListLowQualityPhotos(fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) ScorePhotoQuality(photoId string) (bool, error) {
	args := m.Called(photoId)
	return args.Bool(0), args.Error(1)
}

func (m *MockPhotoService) ScoreAllPhotoQuality(batchSize int) (int, error) {
	args := m.Called(batchSize)
	return args.Int(0), args.Error(1)
}

func (m *MockPhotoService) Search(query string, limit int) ([]*domain.Photo, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) GetTimeline() ([]domain.TimelineEntry, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.TimelineEntry), args.Error(1)
}

func (m *MockPhotoService) GetGeodata(north, south, east, west float64) ([]domain.PhotoGeoData, error) {
	args := m.Called(north, south, east, west)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PhotoGeoData), args.Error(1)
}

func (m *MockPhotoService) RotatePhoto(photoId string, direction string) (*domain.Photo, error) {
	args := m.Called(photoId, direction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) EditPhoto(photoId string, params domain.EditParams) (*domain.Photo, error) {
	args := m.Called(photoId, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) ClearEdits(photoId string) (*domain.Photo, error) {
	args := m.Called(photoId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) DownloadPhotos(photoIds []string, writer io.Writer) error {
	args := m.Called(photoIds, writer)
	return args.Error(0)
}

func (m *MockPhotoService) GetAllTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPhotoService) ListPhotosByTags(tags []string, fromId string, limit int, includeThumbnail bool) ([]*domain.Photo, error) {
	args := m.Called(tags, fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) UpdatePhotoTags(photoId string, tags []string) (*domain.Photo, error) {
	args := m.Called(photoId, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) BatchUpdatePhotoTags(photoIds []string, tags []string, operation string) error {
	args := m.Called(photoIds, tags, operation)
	return args.Error(0)
}

func (m *MockPhotoService) GetDuplicatePhotos() ([]*domain.Photo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) GenerateThumbnailForPhoto(photoId string) (string, error) {
	args := m.Called(photoId)
	return args.String(0), args.Error(1)
}

func (m *MockPhotoService) PhotoThumbnailPathForSize(photoId string, size string) (string, error) {
	args := m.Called(photoId, size)
	return args.String(0), args.Error(1)
}

func (m *MockPhotoService) PhotoThumbnailBytesForSize(photoId string, size string) ([]byte, error) {
	return nil, nil
}

func (m *MockPhotoService) AnalyzeExistingPhotos() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPhotoService) TriggerAIAnalysis() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPhotoService) RegenerateThumbnails(ctx context.Context) {}

// MockJobService is a mock implementation of port.JobService
type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) IsJobRunning(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func (m *MockJobService) GetAllJobs() ([]domain.Job, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Job), args.Error(1)
}

func (m *MockJobService) UpdateAllJobsStatus(status string) error {
	args := m.Called(status)
	return args.Error(0)
}

func (m *MockJobService) JobStart(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockJobService) JobComplete(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockJobService) CreateBaseJobs() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockJobService) StartJobIfNotRunning(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

// TestPhotoHandler_DeletePhoto_Success tests successful photo deletion
func TestPhotoHandler_DeletePhoto_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockPhotoSvc.On("DeletePhoto", "photo123").Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/photos/photo123", nil)

	handler.DeletePhoto(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Photo moved to trash", response.Data["message"])

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_DeletePhoto_MissingID tests deleting a photo without an ID
func TestPhotoHandler_DeletePhoto_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: ""}}
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/photos/", nil)

	handler.DeletePhoto(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "photo ID is required", response["error"])

	mockPhotoSvc.AssertNotCalled(t, "DeletePhoto", mock.Anything)
}

// TestPhotoHandler_RestorePhoto_Success tests successful photo restoration
func TestPhotoHandler_RestorePhoto_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockPhotoSvc.On("RestorePhoto", "photo123").Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/photo123/restore", nil)

	handler.RestorePhoto(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Photo restored", response.Data["message"])

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_ListTrashPhotos_Success tests listing trash photos
func TestPhotoHandler_ListTrashPhotos_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expectedPhotos := []*domain.Photo{
		{ID: "p1", Name: "Trashed Photo 1", CreatedEpoch: 1000},
		{ID: "p2", Name: "Trashed Photo 2", CreatedEpoch: 2000},
	}

	mockPhotoSvc.On("ListTrashPhotos", "", 10, false).Return(expectedPhotos, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/photos/trash", nil)

	handler.ListTrashPhotos(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response pageableResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_EmptyTrash_Success tests emptying the trash
func TestPhotoHandler_EmptyTrash_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockPhotoSvc.On("EmptyTrash").Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodDelete, "/photos/trash", nil)

	handler.EmptyTrash(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Trash emptied", response.Data["message"])

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_SetFavorite_Success tests setting a photo as favorite
func TestPhotoHandler_SetFavorite_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expectedPhoto := &domain.Photo{
		ID:       "photo123",
		Name:     "Test Photo",
		Favorite: true,
	}

	mockPhotoSvc.On("SetFavorite", "photo123", true).Return(expectedPhoto, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
	body, _ := json.Marshal(map[string]interface{}{"favorite": true})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/photo123/favorite", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.SetFavorite(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool         `json:"success"`
		Data    domain.Photo `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "photo123", response.Data.ID)
	assert.True(t, response.Data.Favorite)

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_SetFavorite_DefaultsToFalse tests setting favorite with empty body defaults to false
func TestPhotoHandler_SetFavorite_DefaultsToFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockPhotoSvc.On("SetFavorite", "photo123", false).Return(&domain.Photo{ID: "photo123", Favorite: false}, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
	body, _ := json.Marshal(map[string]interface{}{})
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/photos/photo123/favorite", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.SetFavorite(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Message string `json:"message"`
		Data    any  `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_DownloadPhotos_Success tests downloading photos
func TestPhotoHandler_DownloadPhotos_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockPhotoSvc.On("DownloadPhotos", []string{"photo1", "photo2"}, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		writer := args.Get(1).(io.Writer)
		_, _ = writer.Write([]byte("PK"))
	})

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(downloadPhotosRequest{PhotoIds: []string{"photo1", "photo2"}})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/download", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.DownloadPhotos(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/zip", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "photobox-download.zip")

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_RotatePhoto_Success tests rotating a photo
func TestPhotoHandler_RotatePhoto_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expectedPhoto := &domain.Photo{
		ID:   "photo123",
		Name: "Test Photo",
	}

	mockPhotoSvc.On("RotatePhoto", "photo123", "cw").Return(expectedPhoto, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/photo123/rotate?direction=cw", nil)

	handler.RotatePhoto(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool         `json:"success"`
		Data    domain.Photo `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "photo123", response.Data.ID)

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_RotatePhoto_InvalidDirection tests rotating with invalid direction
func TestPhotoHandler_RotatePhoto_InvalidDirection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
	ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/photo123/rotate?direction=invalid", nil)

	handler.RotatePhoto(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "direction must be cw or ccw", response["error"])

	mockPhotoSvc.AssertNotCalled(t, "RotatePhoto", mock.Anything, mock.Anything)
}

// TestPhotoHandler_GetTimeline_Success tests getting the photo timeline
func TestPhotoHandler_GetTimeline_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expectedEntries := []domain.TimelineEntry{
		{Year: 2024, Month: 1, Count: 10},
		{Year: 2024, Month: 2, Count: 5},
	}

	mockPhotoSvc.On("GetTimeline").Return(expectedEntries, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/photos/timeline", nil)

	handler.GetTimeline(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Data    []domain.TimelineEntry `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, 2024, response.Data[0].Year)

	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_GetGeodata_Success tests getting photo geodata
func TestPhotoHandler_GetGeodata_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expectedGeoData := []domain.PhotoGeoData{
		{ID: "photo1", Lat: 51.5, Lng: -0.1},
		{ID: "photo2", Lat: 40.7, Lng: -74.0},
	}

	mockPhotoSvc.On("GetGeodata", 90.0, -90.0, 180.0, -180.0).Return(expectedGeoData, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/photos/geodata", nil)

	handler.GetGeodata(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                `json:"success"`
		Data    []domain.PhotoGeoData `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, "photo1", response.Data[0].ID)

	mockPhotoSvc.AssertExpectations(t)
}

// TestGetJobStatuses_Success tests successful retrieval of all job statuses
func TestGetJobStatuses_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockJobSvc.On("GetAllJobs").Return([]domain.Job{{Name: "Photo_index", Status: "RUNNING"}}, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/jobs/statuses", nil)

	handler.GetJobStatuses(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Jobs []domain.Job `json:"jobs"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Len(t, response.Data.Jobs, 1)
	assert.Equal(t, "Photo_index", response.Data.Jobs[0].Name)
	assert.Equal(t, "RUNNING", response.Data.Jobs[0].Status)

	mockJobSvc.AssertExpectations(t)
}

// TestGetJobStatuses_Error tests error when GetAllJobs fails
func TestGetJobStatuses_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockJobSvc.On("GetAllJobs").Return(nil, domain.ErrInternal)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/jobs/statuses", nil)

	handler.GetJobStatuses(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockJobSvc.AssertExpectations(t)
}

// TestStopAllJobs_Success tests successfully stopping all jobs
func TestStopAllJobs_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockJobSvc.On("UpdateAllJobsStatus", "NOT_RUNNING").Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/jobs/stop", nil)

	handler.StopAllJobs(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "All jobs stopped", response["message"])

	mockJobSvc.AssertExpectations(t)
}

// TestStopAllJobs_Error tests error when stopping all jobs fails
func TestStopAllJobs_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	mockJobSvc.On("UpdateAllJobsStatus", "NOT_RUNNING").Return(domain.ErrInternal)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/jobs/stop", nil)

	handler.StopAllJobs(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockJobSvc.AssertExpectations(t)
}

// TestPhotoHandler_BatchPhotos tests the unified bulk endpoint
func TestPhotoHandler_BatchPhotos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("add_to_album success", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		mockPhotoSvc.On("BatchAddToAlbum", []string{"p1", "p2"}, "album1").Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(map[string]interface{}{
			"photoIds": []string{"p1", "p2"},
			"action":   "add_to_album",
			"albumId":  "album1",
		})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/batch", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.BatchPhotos(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockPhotoSvc.AssertExpectations(t)
	})

	t.Run("delete success", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		mockPhotoSvc.On("BatchDeletePhotos", []string{"p1", "p2"}).Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(map[string]interface{}{
			"photoIds": []string{"p1", "p2"},
			"action":   "delete",
		})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/batch", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.BatchPhotos(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockPhotoSvc.AssertExpectations(t)
	})

	t.Run("favorite success", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		fav := true
		mockPhotoSvc.On("BatchSetFavorite", []string{"p1"}, true).Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(map[string]interface{}{
			"photoIds": []string{"p1"},
			"action":   "favorite",
			"favorite": true,
		})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/batch", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.BatchPhotos(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		_ = fav
		mockPhotoSvc.AssertExpectations(t)
	})

	t.Run("missing albumId validation", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(map[string]interface{}{
			"photoIds": []string{"p1"},
			"action":   "add_to_album",
		})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/batch", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.BatchPhotos(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockPhotoSvc.AssertNotCalled(t, "BatchAddToAlbum", mock.Anything, mock.Anything)
	})

	t.Run("unsupported action", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(map[string]interface{}{
			"photoIds": []string{"p1"},
			"action":   "explode",
		})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/batch", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.BatchPhotos(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestPhotoHandler_AddPhotosToAlbum tests POST /albums/:id/photos
func TestPhotoHandler_AddPhotosToAlbum(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		mockPhotoSvc.On("BatchAddToAlbum", []string{"p1", "p2"}, "album1").Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "album1"}}
		body, _ := json.Marshal(map[string]interface{}{"photoIds": []string{"p1", "p2"}})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/albums/album1/photos", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.AddPhotosToAlbum(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockPhotoSvc.AssertExpectations(t)
	})

	t.Run("missing photos validation", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "album1"}}
		body, _ := json.Marshal(map[string]interface{}{"photoIds": []string{}})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/albums/album1/photos", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.AddPhotosToAlbum(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockPhotoSvc.AssertNotCalled(t, "BatchAddToAlbum", mock.Anything, mock.Anything)
	})
}

// TestPhotoHandler_UpdatePhotoMetadata tests PATCH /photos/:id/metadata
func TestPhotoHandler_UpdatePhotoMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		expected := &domain.Photo{ID: "photo123", Description: "Beach day"}
		lat := 52.04
		lng := 0.094
		dateTaken := "2023-01-01T21:43:14Z"
		mockPhotoSvc.On("UpdatePhotoMetadata", "photo123", "Beach day", &lat, &lng, &dateTaken).Return(expected, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
		body, _ := json.Marshal(map[string]interface{}{
			"description": "Beach day",
			"latitude":    52.04,
			"longitude":   0.094,
			"dateTaken":   "2023-01-01T21:43:14Z",
		})
		ctx.Request = httptest.NewRequest(http.MethodPatch, "/photos/photo123/metadata", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.UpdatePhotoMetadata(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockPhotoSvc.AssertExpectations(t)
	})

	t.Run("missing id", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(map[string]interface{}{"description": "x"})
		ctx.Request = httptest.NewRequest(http.MethodPatch, "/photos//metadata", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.UpdatePhotoMetadata(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestPhotoHandler_GetPhotoLocation tests GET /photos/:id/location
func TestPhotoHandler_GetPhotoLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		mockPhotoSvc.On("ReverseGeocode", mock.Anything, "photo123").Return("Cambridge, United Kingdom", nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
		ctx.Request = httptest.NewRequest(http.MethodGet, "/photos/photo123/location", nil)

		handler.GetPhotoLocation(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response struct {
			Success bool `json:"success"`
			Data    PhotoLocationResponse `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Cambridge, United Kingdom", response.Data.Location)
		mockPhotoSvc.AssertExpectations(t)
	})
}

// TestPhotoHandler_ListPhotos_LowQuality tests the lowQuality filter
func TestPhotoHandler_ListPhotos_LowQuality(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expected := []*domain.Photo{
		{ID: "p1", Name: "blurry.jpg", QualityScore: 12, IsLowQuality: true},
		{ID: "p2", Name: "solid.jpg", QualityScore: 5, IsLowQuality: true},
	}
	mockPhotoSvc.On("ListLowQualityPhotos", "", 10, false).Return(expected, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/photos?lowQuality=true", nil)

	handler.ListPhotos(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_ListPhotos_WithoutLowQuality tests default routing
func TestPhotoHandler_ListPhotos_WithoutLowQuality(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expected := []*domain.Photo{{ID: "p1", Name: "normal.jpg"}}
	mockPhotoSvc.On("ListPhotos", "", 10, false, "", "", "").Return(expected, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/photos", nil)

	handler.ListPhotos(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPhotoSvc.AssertExpectations(t)
}

// TestPhotoHandler_EditPhoto tests POST /photos/:id/edit
func TestPhotoHandler_EditPhoto(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		rot := 90
		expected := &domain.Photo{ID: "photo123"}
		params := domain.EditParams{Rotate: &rot}
		mockPhotoSvc.On("EditPhoto", "photo123", params).Return(expected, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
		body, _ := json.Marshal(map[string]interface{}{"rotate": 90})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/photos/photo123/edit", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.EditPhoto(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockPhotoSvc.AssertExpectations(t)
	})

	t.Run("missing id", func(t *testing.T) {
		mockPhotoSvc := new(MockPhotoService)
		mockJobSvc := new(MockJobService)
		handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(map[string]interface{}{"rotate": 90})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/photos//edit", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler.EditPhoto(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestPhotoHandler_ClearEdits tests DELETE /photos/:id/edit
func TestPhotoHandler_ClearEdits(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhotoSvc := new(MockPhotoService)
	mockJobSvc := new(MockJobService)
	handler := NewPhotoHandler(mockPhotoSvc, mockJobSvc)

	expected := &domain.Photo{ID: "photo123"}
	mockPhotoSvc.On("ClearEdits", "photo123").Return(expected, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "id", Value: "photo123"}}
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/photos/photo123/edit", nil)

	handler.ClearEdits(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	mockPhotoSvc.AssertExpectations(t)
}
