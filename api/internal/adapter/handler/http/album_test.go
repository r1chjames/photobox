package http

import (
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
