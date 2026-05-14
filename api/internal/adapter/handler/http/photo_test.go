package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func (m *MockPhotoService) PerformPhotoIndex() {
	m.Called()
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

func (m *MockPhotoService) SetFavorite(photoId string, favorite bool) (*domain.Photo, error) {
	args := m.Called(photoId, favorite)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Photo), args.Error(1)
}

func (m *MockPhotoService) ListFavoritePhotos(fromId string, limit int, includeThumbnail bool, startDate string, endDate string) ([]*domain.Photo, error) {
	args := m.Called(fromId, limit, includeThumbnail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Photo), args.Error(1)
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

func (m *MockPhotoService) RegenerateThumbnails() {}

// MockJobService is a mock implementation of port.JobService
type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) IsJobRunning(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
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
		writer.Write([]byte("PK"))
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
