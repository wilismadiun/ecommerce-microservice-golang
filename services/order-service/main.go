package main

import (
	"order-service/config"
	"order-service/src/Infrastructures/database"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	dsn := config.DSN()
	database.ConnectDB(dsn)

	// drivers.Router(router, database.DB)

	router.Run(":3000")
}
