package apiServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"strconv"
	"strings"
)
import "net/http"

func defineAlbumsResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)
	albums := router.Group(fmt.Sprintf("%s/albums", urlBasePath))
	{
		albums.GET("", getAlbumById)
		albums.GET("/count", countAllAlbums)
	}
}

func countAllAlbums(c *gin.Context) {
	albumCount, _ := dbEnv.GetAlbumCount()
	c.JSON(http.StatusOK, gin.H{
		"albumCount": albumCount,
	})
}

func getAlbumById(c *gin.Context) {
	albumId := c.Query("albumId")
	if albumId == "" {
		c.JSON(http.StatusBadRequest, missingQueryParam("album ID"))
	}

	var resp interface{}
	var err error
	if albumId != "" {
		resp, err = dbEnv.GetAlbumById(albumId)
	} else {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		resp, err = dbEnv.GetAllAlbums(page, limit)
	}

	if err != nil {
		c.JSON(http.StatusNotFound, notFoundError("album"))
	}
	c.JSON(http.StatusOK, resp)

}