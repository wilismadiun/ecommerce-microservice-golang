package config

// import (
// 	"fmt"
// 	"os"
// )

// type Config struct {
// 	AppEnv  string
// 	AppPort string
// 	Database DatabaseConfig
// }

// func Load() (*Config, error) {
// 	cfg := &Config{
// 		AppEnv:  getEnv("APP_ENV", "development"),
// 		AppPort: getEnv("APP_PORT", "8080"),
// 		Database: DatabaseConfig{
// 			Host:     getEnv("DATABASE_HOST", "localhost"),
// 			Port:     getEnv("DATABASE_PORT", "5432"),
// 			User:     getEnv("DATABASE_USER", "postgres"),
// 			Password: os.Getenv("DATABASE_PASSWORD"),
// 			Name:     getEnv("DATABASE_NAME", "my_project"),
// 			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
// 		},
// 	}

// 	if cfg.Database.Password == "" {
// 		return nil, fmt.Errorf(
// 			"DATABASE_PASSWORD wajib diisi",
// 		)
// 	}

// 	return cfg, nil
// }

// func getEnv(key, fallback string) string {
// 	if value, ok := os.LookupEnv(key); ok {
// 		return value
// 	}
// 	return fallback
// }
