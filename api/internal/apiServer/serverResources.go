package apiServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/components"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"strings"
)
import "net/http"

func defineHealthCheckResources(r *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	r.GET(fmt.Sprintf("%s/health", urlBasePath), func(c *gin.Context) {
		_, err := database.GetAllSettings()
		if err != nil {
			c.JSON(http.StatusBadGateway, err.Error())
		} else {
			c.Status(http.StatusOK)
		}
	})
}

func defineServerResources(r *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	r.POST(fmt.Sprintf("%s/index", urlBasePath), func(c *gin.Context) {
		isRunning, err := database.IsJobRunning("Photo_index")
		if err != nil {
			c.JSON(http.StatusBadRequest, "Requested setting not found")
		}
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

func defineSettingsResources(r *gin.Engine, appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	r.GET(fmt.Sprintf("%s/settings", urlBasePath), func(c *gin.Context) {
		response, err := database.GetAllSettings()
		if err != nil {
			c.JSON(http.StatusNotFound, "Requested setting not found")
		} else {
			c.JSON(http.StatusOK, response)
		}
	})

	//r.POST(fmt.Sprintf("%s/setting", urlBasePath), func(c *gin.Context) {
	//	err := database.UpdateAllSettings(c)
	//	if err != nil {
	//		c.JSON(http.StatusNotFound, err.Error())
	//	} else {
	//		c.Status(http.StatusAccepted)
	//	}
	//})
}