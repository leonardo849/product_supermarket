package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/config"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/auth"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/pact"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/product"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/user"
	domainError "github.com/leonardo849/product_supermarket/internal/domain/error"
)

func SetupApp(productHandler *product.ProductHandler, userHandler *user.UserHandler, errorCache domainError.ErrorCache) *fiber.App {
	app := fiber.New()
	productGroup := app.Group("/product")
	userGroup := app.Group("/user")
	jwtParser := auth.NewParser(config.Load().SecretJWT)
	product.SetupProductRoutes(productGroup, productHandler, jwtParser)
	user.SetupUserRoutes(userGroup, userHandler, jwtParser)
	pactProvider := pact.NewProviderStateHandler(errorCache)
	app.Post("/_pact/provider-states", pactProvider.Handle)
	return app
}
