package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/auth"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/middlewares"
)

func SetupProductRoutes(productGroup fiber.Router, productHandler *ProductHandler, jwtParser *auth.Parser) {
	productGroup.Post("/", middlewares.VerifyJwt(jwtParser), middlewares.CheckRole([]string{"WORKER", "MANAGER"}), productHandler.createProduct())
}