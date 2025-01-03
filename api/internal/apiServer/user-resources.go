package apiServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	. "gitlab.com/r1chjames/photobox/api/internal/token"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"net/http"
	"strings"
	"time"
)

type registrationRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type CreationOrUpdateRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role" binding:"required"`
	Approved bool     `json:"approved"`
	Password string   `json:"password" binding:"required"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func (server *Server) defineUserResources(appConfig AppConfig) {
	urlBasePath := strings.TrimSpace(appConfig.ApiBasePath)

	publicRoutes := server.router.Group(fmt.Sprintf("%s/user", urlBasePath))
	publicRoutes.POST("/login", server.login)
	publicRoutes.POST("/register", server.registerUser)

	authenticatedRoutes := server.router.Group(fmt.Sprintf("%s/user", urlBasePath)).Use(authMiddleware(*server.tokenMaker))
	authenticatedRoutes.POST("/create", server.createUser)
	authenticatedRoutes.POST("/update", server.updateUser)

}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := dbEnv.GetApprovedUserByUsername(req.Username)
	if user.ID == "" || err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	valid, err := ComparePasswordAndHash(req.Password, user.Password)
	if !valid || err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	// Create and send an access token
	accessToken, err := server.tokenMaker.CreateToken(req.Username, time.Hour)
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

func CreateUser(req CreationOrUpdateRequest, approved bool) error {
	hashSalt, err := CreateHash(req.Password, DefaultArgon2idHash())
	if err != nil {
		return errors.New("Error creating hash")
	}

	var user = UserResponse{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Approved: approved,
		Password: hashSalt,
	}

	err = dbEnv.CreateUser(user)
	if err != nil {
		return errors.New("Unable to save user")
	}
	return nil
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreationOrUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := CreateUser(req, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidRequest())
	} else {
		ctx.Status(http.StatusOK)
	}

	return
}

func (server *Server) verifyAuthorizedUserRole(ctx *gin.Context) (string, UserRole, error) {
	authHeader := strings.Fields(GetAuthHeader(ctx))[1]
	payload, err := server.tokenMaker.ParseToken(authHeader)
	if err != nil {
		return "", "", err
	}
	tokenUserObject, err := dbEnv.GetUserByUsername(payload.Username)
	if err != nil {
		return payload.Username, "", err
	}
	return payload.Username, tokenUserObject.Role, nil
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req CreationOrUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userToUpdate UserResponse
	var username, role, _ = server.verifyAuthorizedUserRole(ctx)
	if role == ADMINISTRATOR {
		userToUpdate = UserResponse{
			Username: req.Username,
			Email:    req.Email,
			Role:     req.Role,
			Approved: req.Approved,
		}
	} else if username == req.Username {
		userToUpdate = UserResponse{
			Email:    req.Email,
			Password: req.Password,
		}
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	err := dbEnv.UpdateUser(userToUpdate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidRequest())
	} else {
		ctx.Status(http.StatusOK)
	}

	return
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req registrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashSalt, err := CreateHash(req.Password, DefaultArgon2idHash())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user = UserResponse{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		Role:     VIEWER,
		Approved: false,
		Password: hashSalt,
	}

	err = dbEnv.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, invalidRequest())
	} else {
		ctx.Status(http.StatusOK)
	}

	return
}
