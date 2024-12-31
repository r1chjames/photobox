package apiServer

import (
	"github.com/gin-gonic/gin"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

func (server *Server) defineUserResources(appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	publicRoutes := server.router.Group(urlBasePath)
	publicRoutes.POST("/login", server.login)
	publicRoutes.POST("/create", server.createUser)

	authenticatedRoutes := server.router.Group(urlBasePath).Use(authMiddleware(*server.tokenMaker))
	authenticatedRoutes.DELETE("/delete/:id", server.deleteUser)
}

func (server *Server) login(ctx *gin.Context) {
	// Request binding for login credentials
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := dbEnv.GetUserByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, req.Username)
	} else {
		ctx.JSON(http.StatusOK, user)
	}

	if user.Password == req.Password {
		// Create and send an access token
		accessToken, err := server.tokenMaker.CreateToken(req.Username, time.Minute)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := loginResponse{
			AccessToken: accessToken,
			User:        user,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusForbidden, gin.H{"error": "Incorrect password"})
	return
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	// Request binding for new user details
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign a unique ID and add the user to the list
	user.ID = strconv.Itoa(rand.Intn(1000))
	err := dbEnv.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidRequest())
	} else {
		ctx.JSON(http.StatusOK, user)
	}

	return
}

type deleteUserRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	// Request binding for user ID
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dbEnv.DeleteUser(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidRequest())
	} else {
		ctx.JSON(http.StatusOK, "")
	}

	return
}
