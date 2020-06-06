package apiServer

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

func Start(appConfig AppConfig) {
	r := gin.Default()
	r.Use(cors.Default())

	defineAlbumsResources(r, appConfig)
	definePhotosResources(r, appConfig)
	defineServerResources(r, appConfig)
	defineSettingsResources(r, appConfig)

	err := r.Run()
	if err!= nil {
		log.Fatal("Error starting API Server")
	}
}
