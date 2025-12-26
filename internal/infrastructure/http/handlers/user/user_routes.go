package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/auth"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/middlewares"
)

func SetupUserRoutes(userGroup fiber.Router, userHandler *UserHandler,  jwtParser *auth.Parser) {
	userGroup.Get("/:auth_id/permissions/errors", middlewares.VerifyJwt(jwtParser), middlewares.CheckRole([]string{"MANAGER"}), userHandler.FindIfUserIsInErrors())
}