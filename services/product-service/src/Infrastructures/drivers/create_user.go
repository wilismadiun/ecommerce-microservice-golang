package drivers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(router *gin.Engine, db *gorm.DB) {
	// handler := infrastructures.Container(db)

	// http.StoreRouter(router, handler, &security.AuthenticationTokenJWT{})
}
