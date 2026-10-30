package drivers

import (
	infrastructures "product-service/src/Infrastructures"
	"product-service/src/Infrastructures/security"
	"product-service/src/Interfaces/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(router *gin.Engine, db *gorm.DB) {
	handler := infrastructures.Container(db)

	http.ProductRouter(router, handler, &security.AuthenticationTokenJWT{})
}
