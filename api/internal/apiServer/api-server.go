package apiServer

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

func Start(appConfig AppConfig) {
	router := gin.Default()
	router.Use(cors.Default())

	defineAlbumsResources(router, appConfig)
	definePhotosResources(router, appConfig)
	defineServerResources(router, appConfig)

	err := router.Run()
	if err != nil {
		log.Fatal("Error starting API Server")
	}
}
