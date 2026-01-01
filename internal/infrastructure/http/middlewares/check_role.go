package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thoas/go-funk"
)

func CheckRole(roles []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// cfg := config.Load()

		
		// if cfg.PactMode == "true" {
		// 	return ctx.Next()
		// }

		claims, ok := ctx.Locals("user").(jwt.MapClaims)
		if !ok {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "user not authenticated",
			})
		}

		roleValue, ok := claims["role"]
		if !ok {
			return ctx.Status(403).JSON(fiber.Map{
				"error": "role not found",
			})
		}

		role, ok := roleValue.(string)
		if !ok {
			return ctx.Status(403).JSON(fiber.Map{
				"error": "invalid role format",
			})
		}

		if !funk.Contains(roles, role) {
			return ctx.Status(403).JSON(fiber.Map{
				"error": "you don't have role to do that",
			})
		}

		return ctx.Next()
	}
}
