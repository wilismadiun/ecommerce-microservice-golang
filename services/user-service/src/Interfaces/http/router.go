package http

import (
	"user-service/src/Interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, h *Handler, validator middleware.TokenValidator) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	// profile
	profile := router.Group("/profile")
	profile.Use(middleware.AuthMiddleware(validator))
	{
		profile.POST("", h.AddProfile)
	}
}
