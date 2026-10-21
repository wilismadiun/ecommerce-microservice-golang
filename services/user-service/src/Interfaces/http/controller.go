package http

import (
	"errors"
	"net/http"
	"user-service/src/Applications/usecase"
	"user-service/src/Domains/entities"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	RegisterHandler *usecase.Register
	LoginHandler    *usecase.Login
}

func (h *Handler) Register(c *gin.Context) {
	var user entities.User

	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userSuccess, err := h.RegisterHandler.Execute(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil menambahkan user",
		"data":    userSuccess,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var login entities.Login

	err := c.ShouldBindBodyWithJSON(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := h.LoginHandler.Execute(login)
	if err != nil {
		statusCode := http.StatusBadRequest

		if errors.Is(err, entities.ErrNotFound) {
			statusCode = http.StatusNotFound
		} else if errors.Is(err, entities.ErrInvalidPassword) {
			statusCode = http.StatusUnauthorized
		}

		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data":    token,
	})
}
