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

func defineAlbumsResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	albums := router.Group(fmt.Sprintf("%s/albums", urlBasePath))
	{
		albums.GET("/", func(c *gin.Context) {
			albumId := c.Query("albumId")
			if albumId != "" {
				album, err := database.GetAlbumById(albumId)
				if err != nil {
					c.Status(http.StatusNotFound)
				}
				c.JSON(http.StatusOK, album)
			} else {
				page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
				limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
				albums, err := database.GetAllAlbums(page, limit)
				if err != nil {
					c.Status(http.StatusNotFound)
				}
				c.JSON(http.StatusOK, albums)
			}
		})

		albums.GET("/count", func(c *gin.Context) {
			albumCount, _ := database.GetAlbumCount()
			c.JSON(http.StatusOK, gin.H{
				"albumCount": albumCount,
			})
		})
	}
}
