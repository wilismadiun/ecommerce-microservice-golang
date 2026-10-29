package main

import (
	"product-service/config"
	"product-service/src/Infrastructures/database"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	dsn := config.DSN()
	database.ConnectDB(dsn)

	// drivers.Router(router, database.DB)

	router.Run(":3002")
}
