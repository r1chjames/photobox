package apiServer

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	"gitlab.com/r1chjames/photobox/api/internal/token"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

var dbEnv *database.Env

func NewServer(appConfig AppConfig, env *database.Env) *Server {
	tokenMaker, err := token.NewPaseto("abcdefghijkl12345678901234567890")
	if err != nil {
		log.Fatalf("Couldn't create token maker: %w", err)
	}

	server := &Server{
		tokenMaker: tokenMaker,
	}
	dbEnv = env
	server.setupRouter(appConfig)
	err = server.router.Run()
	if err != nil {
		log.Fatalf("Error starting API Server, %s", err)
	}

	return server
}

func (server *Server) setupRouter(appConfig AppConfig) {
	router := gin.Default()
	router.Use(cors.Default())
	server.router = router

	server.defineAlbumsResources(appConfig)
	server.definePhotosResources(appConfig)
	server.defineServerResources(appConfig)
	server.defineUserResources(appConfig)
}

func setDbEnv(env *database.Env) {
	dbEnv = env
}
