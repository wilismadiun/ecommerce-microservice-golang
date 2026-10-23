package http

import (
	"store-service/src/Interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

func StoreRouter(router *gin.Engine, h *Handler, validator middleware.TokenValidator) {

	store := router.Group("/stores")
	store.Use(middleware.AuthMiddleware(validator))
	{
		store.POST("", h.AddStore)
	}
}
