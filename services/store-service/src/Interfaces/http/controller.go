package http

import (
	"errors"
	"log"
	"net/http"
	"store-service/src/Applications/usecase"
	"store-service/src/Domains/entities"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	AddStoreHandler *usecase.AddStore
}

func (h *Handler) AddStore(c *gin.Context) {
	log.Println("=================================Ini controller===========================================")

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

	var store entities.Store
	err := c.ShouldBindBodyWithJSON(&store)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	existStore, err := h.AddStoreHandler.Execute(userCtx, store)
	if err != nil {
		if errors.Is(err, entities.ErrStoreAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Store Already Exist",
			})
			return
		}

		if errors.Is(err, entities.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User Not Found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Store created successfully",
		"data":    existStore,
		"membership": gin.H{
			"user_id": userCtx,
			"role":    "OWNER",
		},
	})
}
