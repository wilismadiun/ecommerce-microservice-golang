package http

import (
	"order-service/src/Interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRouter(router *gin.Engine, h *Handler, validator middleware.TokenValidator) {
	order := router.Group("/orders")
	order.Use(middleware.AuthMiddleware(validator))
	{
		order.POST("/:id", h.CreateOrder)
	}
}
