package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	HTTPPort    string
	DatabaseURL string
	SecretJWT string
	RabbitURL string
	RedisUri string
	RedisPassword string
	RedisDatabase int
	PactAddress string
	PactMode string
}

func Load() Config {
	env := getEnv("APP_ENV", "DEV")

	loadEnv()

	redisDatabase, err := strconv.Atoi(getEnv("REDIS_DATABASE", "0"))
	if err != nil {
		log.Fatal(err.Error())
	}

	if env == "PROD" {
		
		return Config{
			Env:         env,
			HTTPPort:    mustGetEnv("PORT"),
			DatabaseURL: mustGetEnv("DATABASE_URI"),
			SecretJWT: mustGetEnv("SECRETJWT"),
			RabbitURL: mustGetEnv("RABBIT_URI"),
			RedisUri: mustGetEnv("REDIS_URI"),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
			RedisDatabase: redisDatabase,
			PactAddress: mustGetEnv("PACT_BROKER_BASE_URL"),
			PactMode: getEnv("PACT_MODE", "false"),
		}
	}

	
	return Config{
		Env:         env,
		HTTPPort:    getEnv("PORT", "8080"),
		DatabaseURL: getEnv(
			"DATABASE_URI",
			"postgres://postgres:postgres@localhost:5432/product_db?sslmode=disable",
		),
		SecretJWT:mustGetEnv("SECRETJWT"),
		RabbitURL: mustGetEnv("RABBIT_URI"),
		RedisUri: mustGetEnv("REDIS_URI"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDatabase: redisDatabase,
		PactAddress: getEnv("PACT_BROKER_BASE_URL", "http://localhost:9292"),
		PactMode: getEnv("PACT_MODE", "false"),
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

func loadEnv() {
	env := getEnv("APP_ENV", "DEV")

	if env == "PROD" {
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}
