package config

import (
	"log"
	"os"
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
}

func Load() Config {
	env := getEnv("APP_ENV", "development")

	if env != "PROD" {
		if err := godotenv.Load(); err != nil {
			log.Println(".env wasn't found. using enviroment variables from system")
		}
	}

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
