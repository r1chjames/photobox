package apiServer

import (
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
		pageNumber, _ := strconv.Atoi(c.Query("page_number"))
		pageSize, _ := strconv.Atoi(c.Query("page_size"))

		if photoId != "" {
			response, err := database.GetPhotoInfoById(albumId)
			if err != nil {
				c.JSON(http.StatusNotFound, "Requested photo not found")
			} else {
				c.JSON(http.StatusOK, response)
			}
		} else if albumId != "" {
			response, err := database.GetAllPhotosInfoInAlbum(albumId, pageNumber, pageSize)
			if err != nil {
				c.JSON(http.StatusNotFound, "Requested album not found")
			} else {
				c.JSON(http.StatusOK, response)
			}
		} else {
			response, err := database.GetAllPhotos(pageNumber, pageSize)
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
		photoInfo, _ := database.GetPhotoInfoById(photoId)
		photoPath := photoInfo.FilesystemPath
		c.File(photoPath)
	})

}