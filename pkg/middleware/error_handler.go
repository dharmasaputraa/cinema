package middleware

import (
	"net/http"

	"github.com/dharmasaputraa/cinema-api/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				requestID, _ := c.Get("request_id")

				log.Error("panic recovered",
					zap.Any("error", rec),
					zap.String("path", c.Request.URL.Path),
					zap.Any("request_id", requestID),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"errors": []gin.H{
						{
							"code":    "INTERNAL_ERROR",
							"message": "internal server error",
							"type":    "panic",
						},
					},
				})
				return
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			requestID, _ := c.Get("request_id")

			log.Error("handled error",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Any("request_id", requestID),
			)

			response.Error(c, err)
			c.Abort()
			return
		}
	}
}
