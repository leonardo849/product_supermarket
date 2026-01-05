package testutils

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/bootstrap"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupTestApp() *fiber.App {
	cfg := config.Load()
	container, err := bootstrap.BuildApp(&cfg, false)
	db = container.DB
	if err != nil {
		log.Fatal(err.Error())
	}
	return  container.App
}

func GetDB() *gorm.DB {
	return db
}