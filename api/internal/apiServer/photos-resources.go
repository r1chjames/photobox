package apiServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/components"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"strconv"
	"strings"
)
import "net/http"

var config AppConfig

func definePhotosResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)
	config = appConfig

	photos := router.Group(fmt.Sprintf("%s/photos", urlBasePath))
	{
		photos.GET("", getPhotosInAlbum)
		photos.GET("/count", getPhotoCount)
		router.POST("/index", indexPhotos)
	}

	photo := router.Group(fmt.Sprintf("%s/photo", urlBasePath))
	{
		photo.GET("/bin", getPhoto)
		photo.GET("/thumbnail", getThumbnail)
	}
}

func indexPhotos(c *gin.Context) {
	isRunning, _ := database.IsJobRunning("Photo_index")

	if isRunning {
		c.JSON(http.StatusConflict, "Photo Index already running")
	} else {
		c.Status(http.StatusAccepted)
		database.JobStarting("Photo_index")
		defer database.JobCompleted("Photo_index")

		photoRecords := components.ScanFilesystem(config)
		database.SavePhotoRecordsToDatabase(photoRecords)
	}
}

func getPhotosInAlbum(c *gin.Context) {
	albumId := c.Query("albumId")
	photoId := c.Query("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	includeThumbnails, _ := strconv.ParseBool(c.DefaultQuery("include_thumbnails", "false"))

	var resp interface{}
	var err error

	if photoId != "" {
		resp, err = database.GetPhotoInfoById(albumId)
		if err != nil {
			c.JSON(http.StatusNotFound, notFoundError("photo"))
		}
	} else if albumId != "" {
		resp, err = database.GetAllPhotosInfoInAlbum(albumId, page, limit, includeThumbnails)
		if err != nil {
			c.JSON(http.StatusNotFound, notFoundError("album"))
		}
	} else {
		resp, err = database.GetAllPhotos(page, limit, includeThumbnails)
		if err != nil {
			c.JSON(http.StatusNotFound, "No photos found")
		}
	}

	c.JSON(http.StatusOK, resp)
}

func getPhotoCount(c *gin.Context) {
	albumId := c.Query("albumId")
	if albumId == "" {
		c.JSON(http.StatusBadRequest, missingQueryParam("album ID"))
	}
	response, _ := database.GetPhotosInAlbumCount(albumId)
	c.JSON(http.StatusOK, response)
}

func getPhoto(c *gin.Context) {
	photoId := c.Query("photoId")
	if photoId == "" {
		c.JSON(http.StatusBadRequest, missingQueryParam("photo ID"))
	}

	photoInfo, err := database.GetPhotoInfoById(photoId)
	if err != nil {
		c.JSON(http.StatusNotFound, notFoundError("photo"))
	} else {
		photoPath := photoInfo.FilesystemPath
		c.File(photoPath)
	}
}

func getThumbnail(c *gin.Context) {
	photoId := c.Query("photoId")
	if photoId == "" {
		c.JSON(http.StatusBadRequest, missingQueryParam("photo ID"))
	}

	photoInfo, err := database.GetPhotoInfoById(photoId)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	var photoMetadata PhotoFile
	retrievedPhotoMetadata := photoInfo.Metadata
	err = json.Unmarshal(retrievedPhotoMetadata, &photoMetadata)
	if err != nil || len(photoMetadata.Thumbnail) == 0 {
		c.JSON(http.StatusNotFound, notFoundError("thumbnail"))
	}
	c.Data(http.StatusOK, "application/octet-stream", photoMetadata.Thumbnail)
}