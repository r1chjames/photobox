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
)

// MockApiKeyService is a mock implementation of port.ApiKeyService.
type MockApiKeyService struct {
	mock.Mock
}

func (m *MockApiKeyService) CreateKey(name string, scope domain.ApiKeyScope, createdBy string) (string, *domain.ApiKey, error) {
	args := m.Called(name, scope, createdBy)
	return args.String(0), args.Get(1).(*domain.ApiKey), args.Error(2)
}

func (m *MockApiKeyService) ValidateKey(key string) (*domain.ApiKey, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ApiKey), args.Error(1)
}

func (m *MockApiKeyService) ListKeys() ([]*domain.ApiKey, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ApiKey), args.Error(1)
}

func (m *MockApiKeyService) RevokeKey(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestApiKeyHandler_CreateApiKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockApiKeyService)
	handler := NewApiKeyHandler(mockSvc)

	apiKey := &domain.ApiKey{ID: "key-1", Name: "test", Scope: domain.ApiKeyReadOnly}
	mockSvc.On("CreateKey", "test", domain.ApiKeyReadOnly, "00000000-0000-0000-0000-000000000001").Return("pb_secret", apiKey, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(authorizationPayloadKey, &domain.TokenPayload{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")})
	body, _ := json.Marshal(map[string]string{"name": "test", "scope": "read_only"})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api-keys", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.CreateApiKey(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}
