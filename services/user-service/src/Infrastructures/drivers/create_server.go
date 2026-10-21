package drivers

import (
	infrastructures "user-service/src/Infrastructures"
	"user-service/src/Interfaces/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(router *gin.Engine, db *gorm.DB) {
	handler := infrastructures.Container(db)

	// router
	http.UserRouter(router, handler)
}
