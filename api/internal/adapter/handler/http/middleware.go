package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
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

func authMiddleware(token port.TokenService) gin.HandlerFunc {
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

		accessToken := fields[1]
		_, err := token.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access Token Not Valid"})
			return
		}

		ctx.Next()
	}
}
