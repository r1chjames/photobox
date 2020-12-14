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

	photos := router.Group(fmt.Sprintf("%s/photos", urlBasePath))
	{
		photos.GET("", func(c *gin.Context) {
			albumId := c.Query("albumId")
			photoId := c.Query("id")
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
			includeThumbnails, _ := strconv.ParseBool(c.DefaultQuery("include_thumbnails", "false"))

			if photoId != "" {
				response, err := database.GetPhotoInfoById(albumId)
				if err != nil {
					c.JSON(http.StatusNotFound, "Requested photo not found")
				} else {
					c.JSON(http.StatusOK, response)
				}
			} else if albumId != "" {
				response, err := database.GetAllPhotosInfoInAlbum(albumId, page, limit, includeThumbnails)
				if err != nil {
					c.JSON(http.StatusNotFound, "Requested album not found")
				} else {
					c.JSON(http.StatusOK, response)
				}
			} else {
				response, err := database.GetAllPhotos(page, limit, includeThumbnails)
				if err != nil {
					c.JSON(http.StatusNotFound, "No photos found")
				} else {
					c.JSON(http.StatusOK, response)
				}
			}
		})

		photos.GET("/count", func(c *gin.Context) {
			albumId := c.Query("albumId")
			response, _ := database.GetPhotosInAlbumCount(albumId)
			c.JSON(http.StatusOK, response)
		})
	}

	photo := router.Group(fmt.Sprintf("%s/photo", urlBasePath))
	{
		photo.GET("/bin", func(c *gin.Context) {
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

		photo.GET("/thumbnail", func(c *gin.Context) {
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
				if err != nil || len(photoMetadata.Thumbnail) == 0 {
					c.JSON(http.StatusNotFound, "No thumbnail found")
				}
				c.Data(http.StatusOK, "application/octet-stream", photoMetadata.Thumbnail)
			}
		})
	}
}
