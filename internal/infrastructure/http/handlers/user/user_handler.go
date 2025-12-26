package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/application/user"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	findIfUserIsInErrors *user.FindIfUserIsInErrors
}

func NewUserHandler(findIfUserIsInErrors *user.FindIfUserIsInErrors) *UserHandler {
	return  &UserHandler{
		findIfUserIsInErrors: findIfUserIsInErrors,
	}
}

func(u *UserHandler) FindIfUserIsInErrors() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		authId := user["id"].(string)
		issuedAt, ok := user["iat"].(float64) 
		targetId := ctx.Params("auth_id")
		if !ok {
			return ctx.Status(500).JSON(fiber.Map{"error": "invalid token, 'iat' field not found or has incorrect type"})
		}

		has, err := u.findIfUserIsInErrors.Execute(authId, issuedAt, targetId)
		if err != nil {
			return ctx.Status(200).JSON(fiber.Map{"allowed": false, "error": err.Error()})
		}
		if has  {
			return ctx.Status(200).JSON(fiber.Map{"allowed": false, "error": "he is in errors"})
		}
		
		return  ctx.Status(200).JSON(fiber.Map{"allowed": true, "error": nil})
	}
}