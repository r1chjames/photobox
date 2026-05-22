package http

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		accessToken := ""

		// Try Authorization header first (standard HTTP requests)
		authHeader := GetAuthHeader(ctx)
		if authHeader != "" {
			fields := strings.Fields(authHeader)
			if len(fields) == 2 && strings.ToLower(fields[0]) == authorizationHeaderBearerType {
				accessToken = fields[1]
			}
		}

		// Fallback to query parameter (browser WebSocket API cannot set custom headers)
		if accessToken == "" {
			accessToken = ctx.Query("token")
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
			return
		}

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

type ipLimiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// IPRateLimiter manages rate limiters for each IP address
type IPRateLimiter struct {
	limiters map[string]*ipLimiterEntry
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
	ttl      time.Duration
}

// NewIPRateLimiter creates a new IP-based rate limiter
// rate: maximum requests per second
// burst: maximum burst size
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters: make(map[string]*ipLimiterEntry),
		rate:     r,
		burst:    b,
		ttl:      time.Hour,
	}

	// Start cleanup goroutine to prevent memory leaks
	go limiter.cleanupStaleEntries()

	return limiter
}

// GetLimiter returns the rate limiter for the given IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	entry, exists := i.limiters[ip]
	if !exists {
		entry = &ipLimiterEntry{
			limiter:    rate.NewLimiter(i.rate, i.burst),
			lastAccess: time.Now(),
		}
		i.limiters[ip] = entry
	} else {
		entry.lastAccess = time.Now()
	}

	return entry.limiter
}

// cleanupStaleEntries removes rate limiters that haven't been used within the TTL
func (i *IPRateLimiter) cleanupStaleEntries() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		cutoff := time.Now().Add(-i.ttl)
		for ip, entry := range i.limiters {
			if entry.lastAccess.Before(cutoff) {
				delete(i.limiters, ip)
			}
		}
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

const requestIDHeader = "X-Request-ID"

// requestIDMiddleware generates a unique request ID for each request
// and adds it to both the response headers and the gin context
func requestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Set(requestIDHeader, requestID)
		ctx.Header(requestIDHeader, requestID)
		ctx.Next()
	}
}

// contentTypeMiddleware validates that POST and PUT requests include
// a Content-Type: application/json header
func contentTypeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodPost || ctx.Request.Method == http.MethodPut {
			contentType := ctx.ContentType()
			if contentType != "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				return
			}
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

// securityHeadersMiddleware adds security-related HTTP headers to all responses
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		ctx.Next()
	}
}

// maxBodySizeMiddleware limits the request body size to prevent OOM attacks
func maxBodySizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check Content-Length header early before reading the body
		if ctx.Request.ContentLength > maxSize {
			ctx.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			return
		}
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)
		ctx.Next()
	}
}
