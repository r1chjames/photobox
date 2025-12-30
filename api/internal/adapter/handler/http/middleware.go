package http

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"golang.org/x/time/rate"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
	authorizationPayloadKey       = "authorization_payload"
)

func GetAuthHeader(ctx *gin.Context) string {
	return ctx.GetHeader(authorizationHeaderKey)
}

// GetAuthPayload retrieves the authenticated user's payload from the context
func GetAuthPayload(ctx *gin.Context) *domain.TokenPayload {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil
	}
	return payload.(*domain.TokenPayload)
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
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access Token Not Valid"})
			return
		}

		// Store the payload in the context for downstream handlers
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// requireRole creates middleware that enforces role-based access control
// roles parameter accepts one or more roles that are allowed to access the endpoint
func requireRole(roles ...domain.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := GetAuthPayload(ctx)
		if payload == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// Check if user's role matches any of the allowed roles
		hasPermission := false
		for _, role := range roles {
			if payload.Role == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		ctx.Next()
	}
}

// IPRateLimiter manages rate limiters for each IP address
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter creates a new IP-based rate limiter
// rate: maximum requests per second
// burst: maximum burst size
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}

	// Start cleanup goroutine to prevent memory leaks
	go limiter.cleanupStaleEntries()

	return limiter
}

// GetLimiter returns the rate limiter for the given IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.limiters[ip] = limiter
	}

	return limiter
}

// cleanupStaleEntries removes rate limiters that haven't been used recently
func (i *IPRateLimiter) cleanupStaleEntries() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		// Clear all limiters periodically to prevent unbounded memory growth
		// This is a simple approach - in production, you might want to track last access time
		i.limiters = make(map[string]*rate.Limiter)
		i.mu.Unlock()
	}
}

// rateLimitMiddleware creates middleware that limits requests per IP
func rateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get client IP
		ip := ctx.ClientIP()

		// Get rate limiter for this IP
		ipLimiter := limiter.GetLimiter(ip)

		// Check if request is allowed
		if !ipLimiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		ctx.Next()
	}
}

// timeoutMiddleware creates middleware that enforces a timeout on requests
// This prevents long-running requests from tying up resources
func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Create a context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
		defer cancel()

		// Replace the request context with the timeout context
		ctx.Request = ctx.Request.WithContext(timeoutCtx)

		// Channel to signal when the handler completes
		done := make(chan struct{})

		go func() {
			ctx.Next()
			close(done)
		}()

		select {
		case <-done:
			// Handler completed successfully
			return
		case <-timeoutCtx.Done():
			// Timeout occurred
			ctx.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timeout exceeded",
			})
		}
	}
}
