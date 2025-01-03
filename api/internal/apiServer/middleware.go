package apiServer

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/token"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
)

func GetAuthHeader(ctx *gin.Context) string {
	return ctx.GetHeader(authorizationHeaderKey)
}

func authMiddleware(maker token.PasetoMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := GetAuthHeader(ctx)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No header was passed"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Missing Bearer Token"})
			return
		}

		authType := fields[0]
		if strings.ToLower(authType) != authorizationHeaderBearerType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization Type Not Supported"})
			return
		}

		token := fields[1]
		_, err := maker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access Token Not Valid"})
			return
		}

		ctx.Next()
	}
}
