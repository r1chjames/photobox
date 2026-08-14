package http

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"net/http"
	"time"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

type paginationMetadata struct {
	FromId   string `json:"fromId" example:"1"`
	ToId     string `json:"toId" example:"1"`
	Count    int    `json:"count" example:"1"`
	NextPage string `json:"nextPage" example:"1"`
}

type pageableResponse struct {
	Success  bool               `json:"success" example:"true"`
	Message  string             `json:"message" example:"Success"`
	Data     any                `json:"data,omitempty"`
	Metadata paginationMetadata `json:"metadata,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// newResponse is a helper function to create a response body
func newPageableResponse(success bool, message string, data any, fromId string, toId string, count int, nextPage string) pageableResponse {
	return pageableResponse{
		Success: success,
		Message: message,
		Data:    data,
		Metadata: paginationMetadata{
			FromId:   b64.StdEncoding.EncodeToString([]byte(fromId)),
			ToId:     b64.StdEncoding.EncodeToString([]byte(toId)),
			Count:    count,
			NextPage: fmt.Sprintf(nextPage, b64.StdEncoding.EncodeToString([]byte(fromId))),
		},
	}
}

// meta represents metadata for a paginated response
type meta struct {
	Total int `json:"total" example:"100"`
	Limit int `json:"limit" example:"10"`
	Skip  int `json:"skip" example:"0"`
}

// newMeta is a helper function to create metadata for a paginated response
func newMeta(total, limit, skip int) meta {
	return meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// authResponse represents an authentication response body
type authResponse struct {
	AccessToken string    `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
	ExpiresAt   time.Time `json:"expires_at" example:"2026-08-13T12:00:00Z"`
}

// newAuthResponse is a helper function to create a response body for handling authentication data
func newAuthResponse(token string, expiresAt time.Time) authResponse {
	return authResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
	}
}

// userResponse represents a user response body
type userResponse struct {
	ID        string    `json:"id" example:"b8d44111-bb92-4013-899f-a4fb9f9436bb"`
	Username  string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	Role      string    `json:"role" example:"viewer"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	domain.ErrInternal:                   http.StatusInternalServerError,
	domain.ErrDataNotFound:               http.StatusNotFound,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrInvalidToken:               http.StatusUnauthorized,
	domain.ErrExpiredToken:               http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,
	domain.ErrNoUpdatedData:              http.StatusBadRequest,
	domain.ErrInvalidRequest:             http.StatusBadRequest,
	domain.ErrJobAlreadyRunning:          http.StatusConflict,
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

    if errors.Is(err, domain.ErrDataNotFound) {
        ctx.JSON(http.StatusNotFound, newErrorResponse([]string{err.Error()}))
        return
    }
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// handleAbort sends an error response and aborts the request with the specified status code and error message
func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}

// handleSuccess sends a success response with the specified status code and optional data
func handlePaginatedSuccess(ctx *gin.Context, data any, fromId string, toId string, count int, nextPage string) {
	rsp := newPageableResponse(true, "Success", data, fromId, toId, count, nextPage)
	ctx.JSON(http.StatusOK, rsp)
}
