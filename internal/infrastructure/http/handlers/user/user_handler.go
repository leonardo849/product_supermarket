package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/leonardo849/product_supermarket/internal/application/user"
)

type UserHandler struct {
	findIfUserIsInErrors *user.FindIfUserIsInErrors
}

func NewUserHandler(findIfUserIsInErrors *user.FindIfUserIsInErrors) *UserHandler {
	return  &UserHandler{
		findIfUserIsInErrors: findIfUserIsInErrors,
	}
}


func (u *UserHandler) FindIfUserIsInErrors() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		claimsRaw := ctx.Locals("user")


		claims, ok := claimsRaw.(jwt.MapClaims)
		if !ok {
			return ctx.Status(500).JSON(fiber.Map{
				"error": "invalid jwt claims type",
			})
		}

		authId, ok := claims["id"].(string)
		if !ok {
			return ctx.Status(500).JSON(fiber.Map{
				"error": "invalid token: id missing",
			})
		}

		issuedAtFloat, ok := claims["iat"].(float64)
		if !ok {
			return ctx.Status(500).JSON(fiber.Map{
				"error": "invalid token: iat missing",
			})
		}

		targetId := ctx.Params("auth_id")

		has, err := u.findIfUserIsInErrors.Execute(
			authId,
			issuedAtFloat,
			targetId,
		)

		if err != nil {
			return ctx.Status(200).JSON(fiber.Map{
				"allowed": false,
				"error":   err.Error(),
			})
		}

		if has {
			return ctx.Status(200).JSON(fiber.Map{
				"allowed": false,
				"error":   "he is in errors",
			})
		}

		return ctx.Status(200).JSON(fiber.Map{
			"allowed": true,
			"error":   nil,
		})
	}
}
