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

func defineAlbumsResources(r *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	r.GET(fmt.Sprintf("%s/albums", urlBasePath), func(c *gin.Context) {
		albumId := c.Query("albumId")
		if albumId != "" {
			album, err := database.GetAlbumById(albumId)
			if err != nil {
				c.JSON(http.StatusOK, album)
			}
			c.Status(http.StatusNotFound)
		} else {
			pageNumber, _ := strconv.Atoi(c.Query("page_number"))
			pageSize, _ := strconv.Atoi(c.Query("page_size"))
			albums, err := database.GetAllAlbums(pageNumber, pageSize)
			if err != nil {
				c.JSON(http.StatusOK, albums)
			}
			c.Status(http.StatusNotFound)
		}
	})

	r.GET(fmt.Sprintf("%s/albums/count", urlBasePath), func(c *gin.Context) {
		albumCount, _ := database.GetAlbumCount()
		c.JSON(http.StatusOK, gin.H{
			"albumCount": albumCount,
		})
	})
}
