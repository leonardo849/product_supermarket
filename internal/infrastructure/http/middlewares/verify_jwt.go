package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/auth"
)

func VerifyJwt(jwtParser *auth.Parser) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return  ctx.Status(401).JSON(fiber.Map{"error": "there isn't token"})
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return  ctx.Status(401).JSON(fiber.Map{"error": "your token is wrong"})
		}
		if parts[0] != "Bearer" {
			return  ctx.Status(401).JSON(fiber.Map{"error": "the token is without the prefix 'bearer'"})
		}
		tokenString := parts[1]
		mapClaims, err := jwtParser.ParseJWT(tokenString)
		if err != nil {
			return  ctx.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		ctx.Locals("user", *mapClaims)
		return ctx.Next()
	}
}