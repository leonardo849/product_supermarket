package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/auth"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/user"
)

func SetupApp(productHandler *product.ProductHandler, userHandler *user.UserHandler) *fiber.App {
	app := fiber.New()
	productGroup := app.Group("/product")
	userGroup := app.Group("/user")
	jwtParser := auth.NewParser(config.Load().SecretJWT)
	product.SetupProductRoutes(productGroup, productHandler, jwtParser)
	user.SetupUserRoutes(userGroup, userHandler, jwtParser)
	return app
}
