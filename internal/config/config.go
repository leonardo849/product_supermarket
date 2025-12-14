package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	HTTPPort    string
	DatabaseURL string
}

func Load() Config {
	env := getEnv("APP_ENV", "development")

	if env != "PROD" {
		if err := godotenv.Load(); err != nil {
			log.Println(".env wasn't found. using enviroment variables from system")
		}
	}

	if env == "PROD" {
		return Config{
			Env:         env,
			HTTPPort:    mustGetEnv("PORT"),
			DatabaseURL: mustGetEnv("DATABASE_URI"),
		}
	}

	
	return Config{
		Env:         env,
		HTTPPort:    getEnv("PORT", "8080"),
		DatabaseURL: getEnv(
			"DATABASE_URI",
			"postgres://postgres:postgres@localhost:5432/product_db?sslmode=disable",
		),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("env %s wasn't defined", key)
	}
	return value
}
