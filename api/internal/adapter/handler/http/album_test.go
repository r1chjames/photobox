package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/datatypes"
)

// MockAlbumService is a mock implementation of port.AlbumService
type MockAlbumService struct {
	mock.Mock
}

func (m *MockAlbumService) GetAlbumById(id string) (*domain.Album, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) GetAlbumByName(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) ListAlbums(fromId string, pageSize int) ([]*domain.Album, error) {
	args := m.Called(fromId, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Album), args.Error(1)
}

func (m *MockAlbumService) AlbumCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAlbumService) CreateAlbum(name string) (*domain.Album, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Album), args.Error(1)
}

func (m *MockAlbumService) CreateAlbumIfNotExists(name string) (*domain.Album, error) {
	args := m.Called(name)
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

// TestAlbumHandler_GetAlbum_Success tests successful album retrieval
func TestAlbumHandler_GetAlbum_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAlbumService)
	handler := NewAlbumHandler(mockService)

	expectedAlbum := &domain.Album{
		ID:          "1",
		Name:        "Test Album",
		Description: "Test Description",
		Tags:        "test,album",
		Metadata:    datatypes.JSON([]byte("{}")),
	}

	mockService.On("GetAlbumById", "1").Return(expectedAlbum, nil)

	// Use a router to properly set up the context
	router := gin.Default()
	router.GET("/albums/:id", handler.GetAlbum)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/albums/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool          `json:"success"`
		Message string        `json:"message"`
		Data    domain.Album  `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "1", response.Data.ID)
	assert.Equal(t, "Test Album", response.Data.Name)
	assert.Equal(t, "Test Description", response.Data.Description)

	mockService.AssertExpectations(t)
}

//func TestGetAlbumById(t *testing.T) {
//	expectedBody := gin.H{
//		"id":          "album1",
//		"name":        "name1",
//		"description": "desc1",
//		"tags":        "",
//		"metadata":    "",
//		"createdAt":   "0001-01-01T00:00:00Z",
//		"updatedAt":   "0001-01-01T00:00:00Z",
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `albums` WHERE `albums`.`id` = (.+) ORDER BY `albums`.`id` LIMIT 1", sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("album1", "name1", "desc1"))
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/albums", map[string]string{"albumId": "album1"}, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response map[string]string
//	err := json.Unmarshal([]byte(w.Body.String()), &response)
//	assert.Nil(t, err)
//	assert.Equal(t, expectedBody["id"], response["id"])
//	assert.Equal(t, expectedBody["name"], response["name"])
//	assert.Equal(t, expectedBody["description"], response["description"])
//	assert.Equal(t, expectedBody["tags"], response["tags"])
//	assert.Equal(t, expectedBody["metadata"], response["metadata"])
//	assert.Equal(t, expectedBody["createdAt"], response["createdAt"])
//	assert.Equal(t, expectedBody["updatedAt"], response["updatedAt"])
//}
//
//func TestGetAlbumByIdWhenNotExists(t *testing.T) {
//	expectedBody := apiError{
//		Code:    404,
//		Message: "Requested album not found",
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnNotFoundErrorForQuery(mock, "SELECT (.+) FROM `albums` WHERE `albums`.`id` = (.+) ORDER BY `albums`.`id` LIMIT 1")
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/albums", map[string]string{"albumId": "album1"}, nil)
//	assert.Equal(t, http.StatusNotFound, w.Code)
//
//	var responseBody apiError
//	json.Unmarshal([]byte(w.Body.String()), &responseBody)
//	assert.Equal(t, expectedBody, responseBody)
//}
//
//func TestGetAlbumCount(t *testing.T) {
//	expectedBody := gin.H{
//		"albumCount": 1,
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `albums`", sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("album1", "name1", "desc1"))
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/albums/count", nil, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response map[string]int
//	json.Unmarshal([]byte(w.Body.String()), &response)
//
//	assert.Equal(t, expectedBody["albumCount"], response["albumCount"])
//}
//
//func TestGetAlbumCountWhenNoneExist(t *testing.T) {
//	expectedBody := gin.H{
//		"albumCount": 0,
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnNotFoundErrorForQuery(mock, "SELECT (.+) FROM `albums`")
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/albums/count", nil, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response map[string]int
//	json.Unmarshal([]byte(w.Body.String()), &response)
//
//	assert.Equal(t, expectedBody["albumCount"], response["albumCount"])
//}

// TestAlbumHandler_CreateAlbum_Success tests successful album creation
func TestAlbumHandler_CreateAlbum_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAlbumService)
	handler := NewAlbumHandler(mockService)

	expectedAlbum := &domain.Album{
		ID:          "1",
		Name:        "New Album",
		Description: "A new album",
	}

	mockService.On("CreateAlbum", "New Album").Return(&domain.Album{ID: "1", Name: "New Album"}, nil)
	mockService.On("UpdateAlbum", "1", mock.Anything).Return(expectedAlbum, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(createAlbumRequest{Name: "New Album", Description: "A new album"})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/albums", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.CreateAlbum(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool         `json:"success"`
		Data    domain.Album `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "1", response.Data.ID)
	assert.Equal(t, "New Album", response.Data.Name)
	assert.Equal(t, "A new album", response.Data.Description)

	mockService.AssertExpectations(t)
}

// TestAlbumHandler_CreateAlbum_InvalidBody tests creating an album with invalid body
func TestAlbumHandler_CreateAlbum_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAlbumService)
	handler := NewAlbumHandler(mockService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/albums", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.CreateAlbum(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)

	mockService.AssertNotCalled(t, "CreateAlbum", mock.Anything)
}

// TestAlbumHandler_UpdateAlbum_Success tests successful album update
func TestAlbumHandler_UpdateAlbum_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAlbumService)
	handler := NewAlbumHandler(mockService)

	expectedAlbum := &domain.Album{
		ID:          "1",
		Name:        "Updated Album",
		Description: "Updated description",
	}

	mockService.On("UpdateAlbum", "1", mock.Anything).Return(expectedAlbum, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	body, _ := json.Marshal(updateAlbumRequest{Name: "Updated Album", Description: "Updated description"})
	ctx.Request = httptest.NewRequest(http.MethodPut, "/albums/1", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateAlbum(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool         `json:"success"`
		Data    domain.Album `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "1", response.Data.ID)
	assert.Equal(t, "Updated Album", response.Data.Name)

	mockService.AssertExpectations(t)
}

// TestAlbumHandler_UpdateAlbum_NotFound tests updating a non-existent album
func TestAlbumHandler_UpdateAlbum_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAlbumService)
	handler := NewAlbumHandler(mockService)

	mockService.On("UpdateAlbum", "nonexistent", mock.Anything).Return(nil, domain.ErrDataNotFound)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	body, _ := json.Marshal(updateAlbumRequest{Name: "Updated"})
	ctx.Request = httptest.NewRequest(http.MethodPut, "/albums/nonexistent", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateAlbum(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// TestAlbumHandler_DeleteAlbum_Success tests successful album deletion
func TestAlbumHandler_DeleteAlbum_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAlbumService)
	handler := NewAlbumHandler(mockService)

	mockService.On("DeleteAlbum", "1", false).Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/albums/1", nil)

	handler.DeleteAlbum(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Album deleted", response.Data["message"])

	mockService.AssertExpectations(t)
}
