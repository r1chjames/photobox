package apiServer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"testing"
)


func TestGetAlbumById(t *testing.T) {
	expectedBody := gin.H{
		"id": "album1",
		"name": "name1",
		"description": "desc1",
		"tags": "",
		"metadata": "",
		"createdAt": "0001-01-01T00:00:00Z",
		"updatedAt": "0001-01-01T00:00:00Z",
	}

	router := setupRouter(TestAppConfig())
	mockEnv, mock, _ := database.MockDB(t)
	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `albums` WHERE `albums`.`id` = (.+) ORDER BY `albums`.`id` LIMIT 1", sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("album1", "name1", "desc1"))
	setDbEnv(mockEnv)

	w := PerformRequest(router, "GET", "/test/albums", map[string]string{"albumId": "album1"}, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody["id"], response["id"])
	assert.Equal(t, expectedBody["name"], response["name"])
	assert.Equal(t, expectedBody["description"], response["description"])
	assert.Equal(t, expectedBody["tags"], response["tags"])
	assert.Equal(t, expectedBody["metadata"], response["metadata"])
	assert.Equal(t, expectedBody["createdAt"], response["createdAt"])
	assert.Equal(t, expectedBody["updatedAt"], response["updatedAt"])
}

func TestGetAlbumByIdWhenNotExists(t *testing.T) {
	expectedBody := apiError{
		Code:    404,
		Message: "Requested album not found",
	}

	router := setupRouter(TestAppConfig())
	mockEnv, mock, _ := database.MockDB(t)
	database.ShouldReturnNotFoundErrorForQuery(mock, "SELECT (.+) FROM `albums` WHERE `albums`.`id` = (.+) ORDER BY `albums`.`id` LIMIT 1")
	setDbEnv(mockEnv)

	w := PerformRequest(router, "GET", "/test/albums", map[string]string{"albumId": "album1"}, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var responseBody apiError
	json.Unmarshal([]byte(w.Body.String()), &responseBody)
	assert.Equal(t, expectedBody, responseBody)
}

func TestGetAlbumCount(t *testing.T) {
	expectedBody := gin.H{
		"albumCount": 1,
	}

	router := setupRouter(TestAppConfig())
	mockEnv, mock, _ := database.MockDB(t)
	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `albums`", sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("album1", "name1", "desc1"))
	setDbEnv(mockEnv)

	w := PerformRequest(router, "GET", "/test/albums/count", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]int
	json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Equal(t, expectedBody["albumCount"], response["albumCount"])
}

func TestGetAlbumCountWhenNoneExist(t *testing.T) {
	expectedBody := gin.H{
		"albumCount": 0,
	}

	router := setupRouter(TestAppConfig())
	mockEnv, mock, _ := database.MockDB(t)
	database.ShouldReturnNotFoundErrorForQuery(mock, "SELECT (.+) FROM `albums`")
	setDbEnv(mockEnv)

	w := PerformRequest(router, "GET", "/test/albums/count", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]int
	json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Equal(t, expectedBody["albumCount"], response["albumCount"])
}
