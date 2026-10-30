package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"product-service/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func migrations(db *gorm.DB) error {
	var sqlDB *sql.DB

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database connection: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	migrationDir := "./migrations"

	if err := goose.Up(sqlDB, migrationDir); err != nil {
		return fmt.Errorf("run goose migrations: %w", err)
	}

	return nil
}

func main() {
	db, err := connectDB(config.DSN())
	if err != nil {
		log.Printf("Migration setup failed: %v", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Get database connection failed: %v", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	if err := migrations(db); err != nil {
		log.Printf("Migration failed: %v", err)
		os.Exit(1)
	}

	log.Println("Migration completed successfully")
}
