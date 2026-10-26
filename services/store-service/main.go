package main

import (
	"store-service/config"
	"store-service/src/Infrastructures/database"
	"store-service/src/Infrastructures/drivers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	dsn := config.DSN()
	database.ConnectDB(dsn)

	drivers.Router(router, database.DB)

	router.Run(":3001")
}
