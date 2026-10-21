package http

import "github.com/gin-gonic/gin"

func UserRouter(router *gin.Engine, h *Handler) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}
