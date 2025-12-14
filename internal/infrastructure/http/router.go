package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
)

func SetupApp(productHandler *product.ProductHandler) *fiber.App {
	app := fiber.New()
	productGroup := app.Group("/product")
	product.SetupProductRoutes(productGroup, productHandler)
	return app
}
