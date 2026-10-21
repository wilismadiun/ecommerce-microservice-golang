package main

import (
	"log"
	"user-service/config"
	"user-service/src/Infrastructures/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env from the root directory")
	}

	dsn := config.DSN()
	database.ConnectDB(dsn)

	router.Run(":3000")
}
