package apiServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"strings"
)
import "net/http"

func definePhotosResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	router.GET(fmt.Sprintf("%s/photos", urlBasePath), func(c *gin.Context) {
		albumId := c.Query("albumId")
		photoId := c.Query("id")

		if photoId != "" {
			response, err := database.GetPhotoInfoById(appConfig, albumId)
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
			} else {
				c.JSON(http.StatusOK, response)
			}
		} else if albumId != "" {
			response, err := database.GetAllPhotosInfoInAlbum(appConfig, albumId)
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
			} else {
				c.JSON(http.StatusOK, response)
			}
		} else {
			response, err := database.GetAllPhotos(appConfig)
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
			} else {
				c.JSON(http.StatusOK, response)
			}
		}
	})

	router.GET(fmt.Sprintf("%s/photos/count", urlBasePath), func(c *gin.Context) {
		albumId := c.Query("albumId")
		response := database.GetPhotosInAlbumCount(appConfig, albumId)
		c.JSON(http.StatusOK, response)
	})

	router.GET(fmt.Sprintf("%s/photo/bin", urlBasePath), func(c *gin.Context) {
		photoId := c.Query("photoId")
		photoInfo, _ := database.GetPhotoInfoById(appConfig, photoId)
		photoPath := photoInfo.FilesystemPath
		c.File(photoPath)
	})

}