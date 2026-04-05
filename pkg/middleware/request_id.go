package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		// set ke context
		c.Set("request_id", requestID)

		// optional: set ke header response
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
