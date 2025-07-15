package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		visitors: make(map[string]*rate.Limiter),
	}
}

func (m *RateLimitMiddleware) getVisitor(identifier string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	limiter, exists := m.visitors[identifier]
	if !exists {
		// Create a new rate limiter for this visitor
		// Allow 100 requests per minute with a burst of 10
		limiter = rate.NewLimiter(rate.Every(time.Minute/100), 10)
		m.visitors[identifier] = limiter
	}

	return limiter
}

func (m *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get identifier from JWT token or API key
		identifier := c.GetString("tenant_id")
		if identifier == "" {
			identifier = c.ClientIP() // Fallback to IP address
		}

		limiter := m.getVisitor(identifier)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
