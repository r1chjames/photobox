package apiServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"strconv"
	"strings"
)
import "net/http"

func definePhotosResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	router.GET(fmt.Sprintf("%s/photos", urlBasePath), func(c *gin.Context) {
		albumId := c.Query("albumId")
		photoId := c.Query("id")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		if photoId != "" {
			response, err := database.GetPhotoInfoById(albumId)
			if err != nil {
				c.JSON(http.StatusNotFound, "Requested photo not found")
			} else {
				c.JSON(http.StatusOK, response)
			}
		} else if albumId != "" {
			response, err := database.GetAllPhotosInfoInAlbum(albumId, page, limit)
			if err != nil {
				c.JSON(http.StatusNotFound, "Requested album not found")
			} else {
				c.JSON(http.StatusOK, response)
			}
		} else {
			response, err := database.GetAllPhotos(page, limit)
			if err != nil {
				c.JSON(http.StatusNotFound, "No photos found")
			} else {
				c.JSON(http.StatusOK, response)
			}
		}
	})

	router.GET(fmt.Sprintf("%s/photos/count", urlBasePath), func(c *gin.Context) {
		albumId := c.Query("albumId")
		response, _ := database.GetPhotosInAlbumCount(albumId)
		c.JSON(http.StatusOK, response)
	})

	router.GET(fmt.Sprintf("%s/photo/bin", urlBasePath), func(c *gin.Context) {
		photoId := c.Query("photoId")
		if photoId == "" {
			c.JSON(http.StatusBadRequest, "Missing photo ID")
		}
		photoInfo, err := database.GetPhotoInfoById(photoId)
		if err != nil {
			c.Status(http.StatusNotFound)
		} else {
			photoPath := photoInfo.FilesystemPath
			c.File(photoPath)
		}
	})

	router.GET(fmt.Sprintf("%s/photo/thumbnail", urlBasePath), func(c *gin.Context) {
		photoId := c.Query("photoId")
		if photoId == "" {
			c.JSON(http.StatusBadRequest, "Missing photo ID")
		}
		photoInfo, err := database.GetPhotoInfoById(photoId)
		if err != nil {
			c.Status(http.StatusNotFound)
		} else {
			var photoMetadata PhotoFile
			retrievedPhotoMetadata := photoInfo.Metadata
			err := json.Unmarshal(retrievedPhotoMetadata, &photoMetadata)
			if err != nil {
				c.JSON(http.StatusNotFound, "No thumbnail found")
			}
			c.Data(http.StatusOK, "application/octet-stream" ,photoMetadata.Thumbnail)
		}
	})

}