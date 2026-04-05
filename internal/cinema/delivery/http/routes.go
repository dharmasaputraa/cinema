package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *CinemaHandler) {
	cinemas := r.Group("/cinemas")
	{
		cinemas.POST("", h.Create)
		cinemas.GET("", h.List)
		cinemas.GET("/:id", h.GetByID)
		cinemas.PUT("/:id", h.Update)
		cinemas.DELETE("/:id", h.Delete)

		cinemas.POST("/:id/screens", h.AddScreen)
		cinemas.GET("/:id/screens", h.GetScreens)

		cinemas.POST("/screens/:screen_id/seats", h.AddSeats)
	}
}
