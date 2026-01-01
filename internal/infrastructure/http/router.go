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
	// app.Use(func (ctx *fiber.Ctx) error {
	// 	cfg := config.Load()

	// 	if cfg.PactMode != "true" {
	// 		return ctx.Next()
	// 	}

	// 	path := ctx.Path()

		
	// 	if path == "/_pact/provider-states" {
	// 		return ctx.Next()
	// 	}

	// 	if path == "/user/user-123/permissions/errors" {
	// 		return ctx.Next()
	// 	}
	// 	if path == "/health" {
	// 		return ctx.Next()
	// 	}

	// 	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"error": "route disabled during pact verification",
	// 	})
	// })
	productGroup := app.Group("/product")
	userGroup := app.Group("/user")
	jwtParser := auth.NewParser(config.Load().SecretJWT)
	product.SetupProductRoutes(productGroup, productHandler, jwtParser)
	user.SetupUserRoutes(userGroup, userHandler, jwtParser)
	pactProvider := pact.NewProviderStateHandler(errorCache)
	app.Get("/health", func(ctx*fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{"message": "ok"})
	})
	app.Post("/_pact/provider-states", pactProvider.Handle)
	return app
}
