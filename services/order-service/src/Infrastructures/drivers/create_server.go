package drivers

import (
	infrastructures "order-service/src/Infrastructures"
	"order-service/src/Infrastructures/security"
	"order-service/src/Interfaces/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(router *gin.Engine, db *gorm.DB) {
	handler := infrastructures.Container(db)

	http.OrderRouter(router, handler, &security.AuthenticationTokenJWT{})
}
