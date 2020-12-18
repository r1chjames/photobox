package apiServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/components"
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
		photos.GET("", getPhotos)
		photos.GET("/count", getPhotoCount)
		photos.POST("/index", indexPhotos)
	}

	photo := router.Group(fmt.Sprintf("%s/photo", urlBasePath))
	{
		photo.GET("/bin", getPhoto)
		photo.GET("/thumbnail", getThumbnail)
	}
}

func indexPhotos(c *gin.Context) {
	isRunning, _ := dbEnv.IsJobRunning("Photo_index")

	if isRunning {
		c.IndentedJSON(http.StatusConflict, "Photo Index already running")
	} else {
		c.Status(http.StatusAccepted)
		dbEnv.JobStarting("Photo_index")
		defer dbEnv.JobCompleted("Photo_index")

		photoRecords := components.ScanFilesystem(config)
		dbEnv.SavePhotoRecordsToDatabase(photoRecords)
	}
}

func getPhotos(c *gin.Context) {
	albumId := c.Query("albumId")
	photoId := c.Query("photoId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	includeThumbnails, _ := strconv.ParseBool(c.DefaultQuery("include_thumbnails", "false"))

	if photoId != "" {
		resp, err := dbEnv.GetPhotoInfoById(photoId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, apiError{http.StatusNotFound, notFoundError("photo")})
		} else {
			c.IndentedJSON(http.StatusOK, resp)
		}
	} else if albumId != "" {
		resp, err := dbEnv.GetAllPhotosInfoInAlbum(albumId, page, limit, includeThumbnails)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, apiError{http.StatusNotFound, notFoundError("album")})
		} else {
			c.IndentedJSON(http.StatusOK, resp)
		}
	} else {
		resp, err := dbEnv.GetAllPhotos(page, limit, includeThumbnails)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, "No photos found")
		} else {
			c.IndentedJSON(http.StatusOK, resp)
		}
	}
}

func getPhotoCount(c *gin.Context) {
	albumId := c.Query("albumId")
	if albumId == "" {
		c.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("album ID")})
	} else {
		photoCount, _ := dbEnv.GetPhotosInAlbumCount(albumId)
		c.IndentedJSON(http.StatusOK, gin.H{
			"photoCount": photoCount,
		})
	}
}

func getPhoto(c *gin.Context) {
	photoId := c.Query("photoId")
	if photoId == "" {
		c.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("photo ID")})
	} else {

		photoInfo, err := dbEnv.GetPhotoInfoById(photoId)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, notFoundError("photo"))
		} else {
			photoPath := photoInfo.FilesystemPath
			c.File(photoPath)
		}
	}
}

func getThumbnail(c *gin.Context) {
	photoId := c.Query("photoId")
	if photoId == "" {
		c.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("photo ID")})
	} else {

		photoInfo, err := dbEnv.GetPhotoInfoById(photoId)
		if err != nil {
			c.Status(http.StatusNotFound)
		} else {
			var photoMetadata PhotoFile
			retrievedPhotoMetadata := photoInfo.Metadata
			err = json.Unmarshal(retrievedPhotoMetadata, &photoMetadata)
			if err != nil || len(photoMetadata.Thumbnail) == 0 {
				c.IndentedJSON(http.StatusNotFound, notFoundError("thumbnail"))
			} else {
				c.Data(http.StatusOK, "application/octet-stream", photoMetadata.Thumbnail)
			}
		}
	}
}