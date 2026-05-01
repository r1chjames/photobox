package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

// MockShareService is a mock implementation of port.ShareService
type MockShareService struct {
	mock.Mock
}

func (m *MockShareService) CreateShare(resourceType, resourceId, createdBy string, expiry *string, password *string) (*domain.SharedLink, error) {
	args := m.Called(resourceType, resourceId, createdBy, expiry, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SharedLink), args.Error(1)
}

func (m *MockShareService) GetSharedResource(token string, password *string) (*domain.SharedLink, error) {
	args := m.Called(token, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SharedLink), args.Error(1)
}

func (m *MockShareService) ListShares() ([]*domain.SharedLink, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.SharedLink), args.Error(1)
}

func (m *MockShareService) GetSharedResourceData(token string, password *string) (*port.SharedResourceData, error) {
	args := m.Called(token, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.SharedResourceData), args.Error(1)
}

func (m *MockShareService) RevokeShare(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

// TestShareHandler_CreateShare_Success tests successful share creation
func TestShareHandler_CreateShare_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockShareService)
	handler := NewShareHandler(mockService)

	share := &domain.SharedLink{
		Token:        "abc123",
		ResourceType: "photo",
		ResourceId:   "photo1",
		CreatedBy:    "user1",
	}

	mockService.On("CreateShare", "photo", "photo1", mock.Anything, (*string)(nil), (*string)(nil)).
		Return(share, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	payload := &domain.TokenPayload{
		ID:       uuid.New(),
		Username: "admin",
		Role:     domain.ADMINISTRATOR,
	}
	ctx.Set("authorization_payload", payload)

	body, _ := json.Marshal(createShareRequest{ResourceType: "photo", ResourceId: "photo1"})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/shares", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.CreateShare(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "abc123", response.Data["token"])
	assert.Contains(t, response.Data["url"], "abc123")

	mockService.AssertExpectations(t)
}

// TestShareHandler_CreateShare_InvalidBody tests creating a share with invalid body
func TestShareHandler_CreateShare_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockShareService)
	handler := NewShareHandler(mockService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	payload := &domain.TokenPayload{
		ID:       uuid.New(),
		Username: "admin",
		Role:     domain.ADMINISTRATOR,
	}
	ctx.Set("authorization_payload", payload)

	body, _ := json.Marshal(map[string]interface{}{})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/shares", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.CreateShare(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)

	mockService.AssertNotCalled(t, "CreateShare", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// TestShareHandler_GetShared_Success tests retrieving a shared resource
func TestShareHandler_GetShared_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockShareService)
	handler := NewShareHandler(mockService)

	share := &domain.SharedLink{
		Token:        "token123",
		ResourceType: "photo",
		ResourceId:   "photo1",
	}

	mockService.On("GetSharedResource", "token123", (*string)(nil)).Return(share, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "token", Value: "token123"}}
	ctx.Request = httptest.NewRequest(http.MethodGet, "/shared/token123", nil)

	handler.GetShared(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool             `json:"success"`
		Data    domain.SharedLink `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "token123", response.Data.Token)

	mockService.AssertExpectations(t)
}

// TestShareHandler_GetShared_NotFound tests retrieving a non-existent shared resource
func TestShareHandler_GetShared_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockShareService)
	handler := NewShareHandler(mockService)

	mockService.On("GetSharedResource", "badtoken", (*string)(nil)).Return(nil, domain.ErrDataNotFound)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "token", Value: "badtoken"}}
	ctx.Request = httptest.NewRequest(http.MethodGet, "/shared/badtoken", nil)

	handler.GetShared(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

// TestShareHandler_GetShared_Expired tests retrieving an expired shared resource
func TestShareHandler_GetShared_Expired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockShareService)
	handler := NewShareHandler(mockService)

	mockService.On("GetSharedResource", "expired", (*string)(nil)).Return(nil, domain.ErrSharedLinkExpired)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "token", Value: "expired"}}
	ctx.Request = httptest.NewRequest(http.MethodGet, "/shared/expired", nil)

	handler.GetShared(ctx)

	assert.Equal(t, http.StatusGone, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "shared link expired", response["error"])

	mockService.AssertExpectations(t)
}

// TestShareHandler_ListShares_Success tests listing shares
func TestShareHandler_ListShares_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockShareService)
	handler := NewShareHandler(mockService)

	shares := []*domain.SharedLink{
		{Token: "token1", ResourceType: "photo"},
		{Token: "token2", ResourceType: "album"},
	}

	mockService.On("ListShares").Return(shares, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/shares", nil)

	handler.ListShares(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                `json:"success"`
		Data    []domain.SharedLink `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)

	mockService.AssertExpectations(t)
}

// TestShareHandler_RevokeShare_Success tests revoking a share
func TestShareHandler_RevokeShare_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockShareService)
	handler := NewShareHandler(mockService)

	mockService.On("RevokeShare", "token123").Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Params = gin.Params{{Key: "token", Value: "token123"}}
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/shares/token123", nil)

	handler.RevokeShare(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Share revoked", response.Data["message"])

	mockService.AssertExpectations(t)
}
