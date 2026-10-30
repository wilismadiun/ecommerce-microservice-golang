package http

import (
	"product-service/src/Interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRouter(router *gin.Engine, h *Handler, validator middleware.TokenValidator) {
	store := router.Group("/products")
	store.Use(middleware.AuthMiddleware(validator))
	{
		store.POST("", h.CreateProduct)
	}
}
