package product

import "github.com/gofiber/fiber/v2"

func SetupProductRoutes(productGroup fiber.Router, productHandler *ProductHandler) {
	productGroup.Post("/", productHandler.createProduct())
}