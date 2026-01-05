package pact

import (
	// "context"

	"github.com/gofiber/fiber/v2"
	domainError "github.com/leonardo849/product_supermarket/internal/domain/error"
)

type ProviderStateHandler struct {
	errorCache domainError.ErrorCache
}

func NewProviderStateHandler(errorCache domainError.ErrorCache) *ProviderStateHandler {
	return &ProviderStateHandler{
		errorCache: errorCache,
	}
}

func (h *ProviderStateHandler) Handle(c *fiber.Ctx) error {
	var req struct {
		State string `json:"state"`
	}

	_ = c.BodyParser(&req)

	switch req.State {
	case "user exists and permissions were evaluated":
		
	}

	return c.SendStatus(fiber.StatusOK)
}
