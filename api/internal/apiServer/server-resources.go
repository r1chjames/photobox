package apiServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/components"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
	"strings"
)
import "net/http"

func defineHealthCheckResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	router.GET(fmt.Sprintf("%s/health", urlBasePath), func(c *gin.Context) {
		_, err := database.GetAllSettings()
		if err != nil {
			c.JSON(http.StatusBadGateway, err.Error())
		} else {
			c.Status(http.StatusOK)
		}
	})
}

func defineServerResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	router.POST(fmt.Sprintf("%s/index", urlBasePath), func(c *gin.Context) {
		isRunning, _ := database.IsJobRunning("Photo_index")

		if isRunning {
			c.JSON(http.StatusConflict, "Photo Index already running")
		} else {
			c.Status(http.StatusAccepted)
			database.JobStarting("Photo_index")
			defer database.JobCompleted("Photo_index")

			photoRecords := components.ScanFilesystem(appConfig)
			database.SavePhotoRecordsToDatabase(photoRecords)
			database.JobCompleted("Photo_index")
		}
	})
}

func defineSettingsResources(router *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	router.GET(fmt.Sprintf("%s/settings", urlBasePath), func(c *gin.Context) {
		response, err := database.GetAllSettings()
		if err != nil {
			c.JSON(http.StatusNotFound, "Requested setting not found")
		} else {
			c.JSON(http.StatusOK, response)
		}
	})

	router.POST(fmt.Sprintf("%s/setting", urlBasePath), func(c *gin.Context) {
		var settings Settings
		err := c.BindJSON(&settings)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, "Payload not valid")
		}

		err = database.UpdateAllSettings(&settings.Settings)
		if err != nil {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.Status(http.StatusAccepted)
		}
	})
}
