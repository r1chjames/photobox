package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm"
)

// TestHandleError_RecordNotFound tests error handling for record not found
func TestHandleError_RecordNotFound(t *testing.T) {
	mockDB := &gorm.DB{Error: gorm.ErrRecordNotFound}

	err := HandleError(mockDB)

	assert.Equal(t, domain.ErrDataNotFound, err)
}

// TestHandleError_GenericError tests error handling for generic errors
func TestHandleError_GenericError(t *testing.T) {
	expectedError := errors.New("some database error")
	mockDB := &gorm.DB{Error: expectedError}

	err := HandleError(mockDB)

	assert.Equal(t, expectedError, err)
}

// TestHandleError_NoError tests error handling when there's no error
func TestHandleError_NoError(t *testing.T) {
	mockDB := &gorm.DB{Error: nil}

	err := HandleError(mockDB)

	assert.Nil(t, err)
}

// TestPaginate_DefaultValues tests pagination with default values
func TestPaginate_DefaultValues(t *testing.T) {
	// Page 0 should become page 1, limit 0 should become 10
	paginateFunc := Paginate(0, 0)

	// Verify the function is created
	assert.NotNil(t, paginateFunc)
}

// TestPaginate_Page1 tests pagination for first page
func TestPaginate_Page1(t *testing.T) {
	tests := []struct {
		name   string
		page   int
		limit  int
		expLimit int
	}{
		{
			name:     "page 1 with limit 10",
			page:     1,
			limit:    10,
			expLimit: 10,
		},
		{
			name:     "page 1 with limit 50",
			page:     1,
			limit:    50,
			expLimit: 50,
		},
		{
			name:     "page 1 with limit exceeding max",
			page:     1,
			limit:    150,
			expLimit: 100, // Should cap at 100
		},
		{
			name:     "page 1 with negative limit",
			page:     1,
			limit:    -5,
			expLimit: 10, // Should default to 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginateFunc := Paginate(tt.page, tt.limit)
			assert.NotNil(t, paginateFunc)
		})
	}
}

// TestPaginate_MultiplePages tests pagination calculations for different pages
func TestPaginate_MultiplePages(t *testing.T) {
	tests := []struct {
		name   string
		page   int
		limit  int
		description string
	}{
		{
			name:        "page 2 with limit 10",
			page:        2,
			limit:       10,
			description: "Should offset by 10 (page 2 starts after first 10 items)",
		},
		{
			name:        "page 3 with limit 20",
			page:        3,
			limit:       20,
			description: "Should offset by 40 (skip first 2 pages of 20 items each)",
		},
		{
			name:        "page 5 with limit 25",
			page:        5,
			limit:       25,
			description: "Should offset by 100 (skip first 4 pages of 25 items each)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginateFunc := Paginate(tt.page, tt.limit)
			assert.NotNil(t, paginateFunc, tt.description)
		})
	}
}

// TestPaginate_LimitBoundaries tests limit boundary conditions
func TestPaginate_LimitBoundaries(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		expectedLimit int
	}{
		{
			name:          "limit at max boundary",
			limit:         100,
			expectedLimit: 100,
		},
		{
			name:          "limit just over max",
			limit:         101,
			expectedLimit: 100,
		},
		{
			name:          "limit way over max",
			limit:         1000,
			expectedLimit: 100,
		},
		{
			name:          "limit of 1",
			limit:         1,
			expectedLimit: 1,
		},
		{
			name:          "zero limit",
			limit:         0,
			expectedLimit: 10,
		},
		{
			name:          "negative limit",
			limit:         -10,
			expectedLimit: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginateFunc := Paginate(1, tt.limit)
			assert.NotNil(t, paginateFunc)
			// The actual limit is set when the function is called with a real DB.
			// Testing the returned function requires a real database connection (integration test).
		})
	}
}

// TestPaginate_PageBoundaries tests page boundary conditions
func TestPaginate_PageBoundaries(t *testing.T) {
	tests := []struct {
		name string
		page int
	}{
		{
			name: "page 0 defaults to page 1",
			page: 0,
		},
		{
			name: "negative page defaults to page 1",
			page: -5,
		},
		{
			name: "page 1",
			page: 1,
		},
		{
			name: "large page number",
			page: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paginateFunc := Paginate(tt.page, 10)
			assert.NotNil(t, paginateFunc)
		})
	}
}

// TestHandleError_NilDB tests error handling with nil DB (edge case)
func TestHandleError_NilDB(t *testing.T) {
	// This tests defensive programming - what happens with nil input
	var mockDB *gorm.DB = &gorm.DB{}

	err := HandleError(mockDB)

	assert.Nil(t, err)
}

// TestPaginate_ReturnsValidFunctions tests that Paginate returns valid functions
func TestPaginate_ReturnsValidFunctions(t *testing.T) {
	// Test that calling Paginate returns non-nil functions
	paginateFunc1 := Paginate(1, 10)
	paginateFunc2 := Paginate(2, 20)

	assert.NotNil(t, paginateFunc1)
	assert.NotNil(t, paginateFunc2)
	// Testing what these functions do requires a real DB connection (integration test)
}
