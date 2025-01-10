package http

//
//import (
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
//	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
//	"gopkg.in/DATA-DOG/go-sqlmock.v1"
//	"net/http"
//	"testing"
//)
//
//func TestGetPhotoByPhotoId(t *testing.T) {
//	expectedBody := gin.H{
//		"id":             "photo1",
//		"name":           "name1",
//		"filesystemPath": "/dir1",
//		"sourcePath":     "https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg",
//		"albumId":        "album1",
//		"tags":           "",
//		"metadata":       "",
//		"createdAt":      "0001-01-01T00:00:00Z",
//		"updatedAt":      "0001-01-01T00:00:00Z",
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `photos` WHERE `photos`.`id` = (.+) ORDER BY `photos`.`id` LIMIT 1", sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id"}).AddRow("photo1", "name1", "/dir1", "album1"))
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/photos", map[string]string{"photoId": "photo1"}, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response map[string]string
//	json.Unmarshal([]byte(w.Body.String()), &response)
//
//	assert.Equal(t, expectedBody["id"], response["id"])
//	assert.Equal(t, expectedBody["name"], response["name"])
//	assert.Equal(t, expectedBody["filesystemPath"], response["filesystemPath"])
//	assert.Equal(t, expectedBody["sourcePath"], response["sourcePath"])
//	assert.Equal(t, expectedBody["albumId"], response["albumId"])
//	assert.Equal(t, expectedBody["tags"], response["tags"])
//	assert.Equal(t, expectedBody["metadata"], response["metadata"])
//	assert.Equal(t, expectedBody["createdAt"], response["createdAt"])
//	assert.Equal(t, expectedBody["updatedAt"], response["updatedAt"])
//}
//
//func TestGetPhotoByPhotoIdWhenIdDoesNotExist(t *testing.T) {
//	expectedBody := apiError{
//		Code:    404,
//		Message: "Requested photo not found",
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnNotFoundErrorForQuery(mock, "SELECT (.+) FROM `photos` WHERE `photos`.`id` = (.+) ORDER BY `photos`.`id` LIMIT 1")
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/photos", map[string]string{"photoId": "photo1"}, nil)
//	assert.Equal(t, http.StatusNotFound, w.Code)
//
//	var responseBody apiError
//	json.Unmarshal([]byte(w.Body.String()), &responseBody)
//	assert.Equal(t, expectedBody, responseBody)
//}
//
//func TestGetPhotoByAlbumId(t *testing.T) {
//	photo1 := domain.Photo{
//		ID:             "photo1",
//		Name:           "name1",
//		FilesystemPath: "/dir1",
//		SourcePath:     "https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg",
//		AlbumId:        "album1",
//	}
//	photo2 := domain.Photo{
//		ID:             "photo2",
//		Name:           "name2",
//		FilesystemPath: "/dir1",
//		SourcePath:     "https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg",
//		AlbumId:        "album1",
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `photos` WHERE album_id = (.+) LIMIT 10", sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id"}).AddRow("photo1", "name1", "/dir1", "album1").AddRow("photo2", "name2", "/dir1", "album1"))
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/photos", map[string]string{"albumId": "album1"}, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response []map[string]string
//	json.Unmarshal([]byte(w.Body.String()), &response)
//
//	assert.Equal(t, photo1.ID, response[0]["id"])
//	assert.Equal(t, photo1.Name, response[0]["name"])
//	assert.Equal(t, photo1.FilesystemPath, response[0]["filesystemPath"])
//	assert.Equal(t, photo1.SourcePath, response[0]["sourcePath"])
//	assert.Equal(t, photo1.AlbumId, response[0]["albumId"])
//	assert.Equal(t, photo2.ID, response[1]["id"])
//	assert.Equal(t, photo2.Name, response[1]["name"])
//	assert.Equal(t, photo2.FilesystemPath, response[1]["filesystemPath"])
//	assert.Equal(t, photo2.SourcePath, response[1]["sourcePath"])
//	assert.Equal(t, photo2.AlbumId, response[1]["albumId"])
//}
//
//func TestGetAllPhotos(t *testing.T) {
//	photo1 := domain.Photo{
//		ID:             "photo1",
//		Name:           "name1",
//		SourcePath:     "https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg",
//		FilesystemPath: "/dir1",
//		AlbumId:        "album1",
//	}
//	photo2 := domain.Photo{
//		ID:             "photo2",
//		Name:           "name2",
//		SourcePath:     "https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg",
//		FilesystemPath: "/dir1",
//		AlbumId:        "album1",
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `photos` LIMIT 10", sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id"}).AddRow("photo1", "name1", "/dir1", "album1").AddRow("photo2", "name2", "/dir1", "album1"))
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/photos", nil, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response []map[string]string
//	json.Unmarshal([]byte(w.Body.String()), &response)
//
//	assert.Equal(t, photo1.ID, response[0]["id"])
//	assert.Equal(t, photo1.Name, response[0]["name"])
//	assert.Equal(t, photo1.SourcePath, response[0]["sourcePath"])
//	assert.Equal(t, photo1.FilesystemPath, response[0]["filesystemPath"])
//	assert.Equal(t, photo1.AlbumId, response[0]["albumId"])
//	assert.Equal(t, photo2.ID, response[1]["id"])
//	assert.Equal(t, photo2.Name, response[1]["name"])
//	assert.Equal(t, photo2.SourcePath, response[1]["sourcePath"])
//	assert.Equal(t, photo2.FilesystemPath, response[1]["filesystemPath"])
//	assert.Equal(t, photo2.AlbumId, response[1]["albumId"])
//}
//
//func TestGetPhotoCount(t *testing.T) {
//	expectedBody := gin.H{
//		"photoCount": 2,
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `photos` WHERE album_id = (.+)", sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id"}).AddRow("photo1", "name1", "/dir1", "album1").AddRow("photo2", "name2", "/dir1", "album1"))
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/photos/count", map[string]string{"albumId": "album1"}, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response map[string]int
//	json.Unmarshal([]byte(w.Body.String()), &response)
//
//	assert.Equal(t, expectedBody["photoCount"], response["photoCount"])
//}
//
//func TestGetPhotoCountWhenNoneExist(t *testing.T) {
//	expectedBody := gin.H{
//		"photoCount": 0,
//	}
//
//	router := setupRouter(TestAppConfig())
//	mockEnv, mock, _ := database.MockDB(t)
//	database.ShouldReturnNotFoundErrorForQuery(mock, "SELECT (.+) FROM `photos` WHERE album_id = (.+)")
//	setDbEnv(mockEnv)
//
//	w := PerformRequest(router, "GET", "/test/photos/count", map[string]string{"albumId": "album1"}, nil)
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response map[string]int
//	json.Unmarshal([]byte(w.Body.String()), &response)
//
//	assert.Equal(t, expectedBody["photoCount"], response["photoCount"])
//}
//
//// TODO needs work to mock filesystem
////func TestStartIndexWhenIndexIsNotRunning(t *testing.T) {
////	router := setupRouter(TestAppConfig())
////	mockEnv, mock, _ := database.MockDB(t)
////	database.ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `jobs` WHERE `jobs`.`name` = (.+) ORDER BY `jobs`.`name` LIMIT 1", sqlmock.NewRows([]string{"name", "status"}).AddRow("Photo_index", "NOT_RUNNING"))
////	database.ShouldReturnRowsForUpdate(mock, "UPDATE `jobs` SET `status`=(.+),`last_run`=(.+),`updated_at`=(.+) WHERE name = (.+) AND `name` = (.+)", 1) // set to RUNNING
////	database.ShouldReturnRowsForUpdate(mock, "UPDATE `jobs` SET `status`=(.+),`last_run`=(.+),`updated_at`=(.+) WHERE name = (.+) AND `name` = (.+)", 1) // set to NOT_RUNNING
////	setDbEnv(mockEnv)
////
////	w := PerformRequest(router, "POST", "/test/photos/index", nil, nil)
////	assert.Equal(t, http.StatusOK, w.Code)
////
////	var response map[string]int
////	json.Unmarshal([]byte(w.Body.String()), &response)
////
////	//assert.Equal(t, expectedBody["photoCount"], response["photoCount"])
////}
