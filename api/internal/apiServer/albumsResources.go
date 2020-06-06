package apiServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"strings"
)
import "net/http"

func defineAlbumsResources(r *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	r.GET(fmt.Sprintf("%s/albums", urlBasePath), func(c *gin.Context) {
		albumId := c.Query("albumId")
		var response interface{}
		if albumId != "" {
			response = database.GetAlbumById(appConfig, albumId)
		} else {
			response = database.GetAllAlbums(appConfig)
		}
		c.JSON(http.StatusOK, response)
	})

	r.GET(fmt.Sprintf("%s/albums/count", urlBasePath), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"albumCount": database.GetAlbumCount(appConfig),
		})
	})
}
