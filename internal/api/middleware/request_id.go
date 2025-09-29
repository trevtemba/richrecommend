package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Adds a request id for logging/bug fix purposes
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.New().String()
		c.Set("request_id", reqID)

		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()
	}
}
