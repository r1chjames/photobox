package apiServer

import (
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
		photo.POST("", addPhoto)
		photo.GET("/:id", getPhoto)
		photo.GET("/:id/thumbnail", getThumbnail)
		photo.GET("/:id/bin", getPhotoBin)
	}
}

func indexPhotos(c *gin.Context) {
	isRunning, _ := dbEnv.IsJobRunning("Photo_index")

	if isRunning {
		c.IndentedJSON(http.StatusConflict, "Photo Index already running")
	} else {
		c.Status(http.StatusAccepted)
		go func() {
			components.PerformPhotoIndex(config, dbEnv)
		}()
	}
}

func getPhotos(c *gin.Context) {
	albumId := c.Query("albumId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	includeThumbnail, _ := strconv.ParseBool(c.DefaultQuery("thumbnail", "false"))

	if albumId != "" {
		resp, err := dbEnv.GetAllPhotosInfoInAlbum(albumId, page, limit, includeThumbnail)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, apiError{http.StatusNotFound, notFoundError("album")})
		} else {
			c.IndentedJSON(http.StatusOK, resp)
		}
	} else {
		resp, err := dbEnv.GetAllPhotos(page, limit, includeThumbnail)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, "No photos found")
		} else {
			c.IndentedJSON(http.StatusOK, resp)
		}
	}
}

func getPhoto(c *gin.Context) {
	photoId := c.Param("id")

	resp, err := dbEnv.GetPhotoInfoById(photoId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, apiError{http.StatusNotFound, notFoundError("photo")})
	} else {
		c.IndentedJSON(http.StatusOK, resp)
	}
}

func addPhoto(c *gin.Context) {
	var photo PhotoUpload
	err := c.BindJSON(&photo)

	photoFile := components.WriteFileToFilesystem(dbEnv, photo)
	dbEnv.SavePhotoRecordsToDatabase([]PhotoFile{photoFile})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, invalidRequest()})
	} else {
		c.Status(http.StatusCreated)
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

func getPhotoBin(c *gin.Context) {
	photoId := c.Param("id")
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
	photoId := c.Param("id")
	if photoId == "" {
		c.IndentedJSON(http.StatusBadRequest, apiError{http.StatusBadRequest, missingQueryParam("photo ID")})
	} else {
		photoInfo, err := dbEnv.GetPhotoInfoById(photoId)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, notFoundError("thumbnail"))
		} else {
			c.Data(http.StatusOK, "application/octet-stream", photoInfo.Thumbnail)
		}
	}
}
