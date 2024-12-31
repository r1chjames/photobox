package apiServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"strconv"
	"strings"
)
import "net/http"

func (server *Server) defineAlbumsResources(appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	server.router.GET(fmt.Sprintf("%s/album/:id", urlBasePath), getAlbumById)
	albums := server.router.Group(fmt.Sprintf("%s/albums", urlBasePath)).Use(authMiddleware(*server.tokenMaker))
	{
		albums.GET("", getAllAlbums)
		albums.GET("/count", countAllAlbums)
	}

}

func countAllAlbums(c *gin.Context) {
	albumCount, _ := dbEnv.GetAlbumCount()
	c.IndentedJSON(http.StatusOK, gin.H{
		"albumCount": albumCount,
	})
}

func getAllAlbums(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))

	resp, err := dbEnv.GetAllAlbums(page, limit)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, &apiError{http.StatusNotFound, notFoundError("album")})
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func getAlbumById(c *gin.Context) {
	albumId := c.Param("id")
	if albumId == "" {
		c.IndentedJSON(http.StatusBadRequest, &apiError{http.StatusNotFound, missingQueryParam("album ID")})
	} else {
		resp, err := dbEnv.GetAlbumById(albumId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, &apiError{http.StatusNotFound, notFoundError("album")})
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}
