package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockAuthService is a mock implementation of port.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

// TestAuthHandler_Login_Success tests successful login
func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService, time.Hour)

	loginReq := loginRequest{
		Username: "testuser",
		Password: "password123",
	}

	expectedToken := "mock.jwt.token"

	mockService.On("Login", mock.Anything, "testuser", "password123").Return(expectedToken, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(loginReq)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Login(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool         `json:"success"`
		Message string       `json:"message"`
		Data    authResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, expectedToken, response.Data.AccessToken)
	assert.False(t, response.Data.ExpiresAt.IsZero())
	assert.WithinDuration(t, time.Now().Add(time.Hour), response.Data.ExpiresAt, time.Minute)

	mockService.AssertExpectations(t)
}

// TestAuthHandler_Login_InvalidCredentials tests login with invalid credentials
func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService, time.Hour)

	loginReq := loginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	mockService.On("Login", mock.Anything, "testuser", "wrongpassword").
		Return("", domain.ErrInvalidCredentials)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(loginReq)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Login(ctx)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Messages[0], "invalid")

	mockService.AssertExpectations(t)
}

// TestAuthHandler_Login_ValidationError tests login with missing required fields
func TestAuthHandler_Login_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		requestBody map[string]interface{}
		description string
	}{
		{
			name:        "missing username",
			requestBody: map[string]interface{}{"password": "password123"},
			description: "Username is required",
		},
		{
			name:        "missing password",
			requestBody: map[string]interface{}{"username": "testuser"},
			description: "Password is required",
		},
		{
			name:        "empty username",
			requestBody: map[string]interface{}{"username": "", "password": "password123"},
			description: "Username cannot be empty",
		},
		{
			name:        "empty password",
			requestBody: map[string]interface{}{"username": "testuser", "password": ""},
			description: "Password cannot be empty",
		},
		{
			name:        "password too short",
			requestBody: map[string]interface{}{"username": "testuser", "password": "short"},
			description: "Password must be at least 8 characters",
		},
		{
			name:        "empty body",
			requestBody: map[string]interface{}{},
			description: "Both fields required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			handler := NewAuthHandler(mockService, time.Hour)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			handler.Login(ctx)

			assert.Equal(t, http.StatusBadRequest, w.Code, tt.description)

			var response errorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.False(t, response.Success)
			assert.NotEmpty(t, response.Messages)

			// Service should not be called for validation errors
			mockService.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

// TestAuthHandler_Login_InternalError tests login with internal server error
func TestAuthHandler_Login_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService, time.Hour)

	loginReq := loginRequest{
		Username: "testuser",
		Password: "password123",
	}

	mockService.On("Login", mock.Anything, "testuser", "password123").
		Return("", domain.ErrInternal)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(loginReq)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Login(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Messages[0], "internal")

	mockService.AssertExpectations(t)
}

// TestAuthHandler_Login_InvalidJSON tests login with malformed JSON
func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService, time.Hour)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Send invalid JSON
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid json"))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.Login(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)

	mockService.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}
