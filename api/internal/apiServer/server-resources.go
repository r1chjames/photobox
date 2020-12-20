package apiServer

import (
	"github.com/gin-gonic/gin"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
	"strings"
)
import "net/http"

func defineServerResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	server := router.Group(urlBasePath)
	{
		server.GET("/health", healthCheck)
		server.GET("/settings", getAllSettings)
		server.POST("/settings", updateAllSettings)
	}
}

func healthCheck(c *gin.Context) {
	_, err := dbEnv.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusBadGateway, err.Error())
	} else {
		c.Status(http.StatusOK)
	}
}

func getAllSettings(c *gin.Context) {
	response, err := dbEnv.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusNotFound, notFoundError("setting"))
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func updateAllSettings(c *gin.Context) {
	var settings Settings
	err := c.BindJSON(&settings)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, "Payload not valid")
	}

	err = dbEnv.UpdateAllSettings(&settings.Settings)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		c.Status(http.StatusAccepted)
	}
}
