package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/buiminhduc234/audit-log-api/internal/utils"
)

// InjectContext middleware injects gin.Context into context.Context
func InjectContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), utils.GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
