package apiServer

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

var dbEnv *database.Env

func Start(appConfig AppConfig, env *database.Env) {
	dbEnv = env

	router := setupRouter(appConfig)

	err := router.Run()
	if err != nil {
		log.Fatalf("Error starting API Server, %s", err)
	}
}

func setupRouter(appConfig AppConfig) *gin.Engine {
	router := gin.Default()
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://photobox"}
	//router.Use(cors.New(config))
	router.Use(cors.Default())

	defineAlbumsResources(router, appConfig)
	definePhotosResources(router, appConfig)
	defineServerResources(router, appConfig)
	return router
}

func setDbEnv(env *database.Env) {
	dbEnv = env
}
