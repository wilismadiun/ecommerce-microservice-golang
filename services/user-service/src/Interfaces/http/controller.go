package http

import (
	"errors"
	"net/http"
	"user-service/src/Applications/usecase"
	"user-service/src/Domains/entities"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	RegisterHandler         *usecase.Register
	LoginHandler            *usecase.Login
	AddProfileHandler       *usecase.AddProfile
	GetProfileByUserHandler *usecase.GetProfileByUser
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

		if errors.Is(err, entities.ErrUserNotFound) {
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

func (h *Handler) AddProfile(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: user context not found",
		})
		return
	}

	userCtx, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to parse user context",
		})
		return
	}

	var profile entities.Profile
	err := c.ShouldBindBodyWithJSON(&profile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	profile.UserId = userCtx

	existProfile, err := h.AddProfileHandler.Execute(profile)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, entities.ErrProfileAlreadyExists) {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile created successfully",
		"data": gin.H{
			"profile": existProfile,
		},
	})
}

func (h *Handler) GetProfileByUser(c *gin.Context) {
	userId := c.Param("id")

	profile, err := h.GetProfileByUserHandler.Execute(userId)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, entities.ErrProfileNotFound) {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    profile,
	})

}
