package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/leonardo849/product_supermarket/internal/application/product"
	appProduct "github.com/leonardo849/product_supermarket/internal/application/product"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/http/handlers/http_dto"
)

type ProductHandler struct {
	createProductUC *product.CreateProductUseCase
}

func (p *ProductHandler) createProduct() fiber.Handler{
	return  func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		authId := user["id"].(string)
		issuedAt, ok := user["iat"].(float64) 
		if !ok {
			return ctx.Status(500).JSON(fiber.Map{"error": "invalid token, 'iat' field not found or has incorrect type"})
		}

		var input http_dto.CreateProductDTOHttp
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		uuid, err := p.createProductUC.Execute(appProduct.CreateProductInput{
			Name: input.Name,
			PriceInCents: input.PriceInCents,
			Category: input.Category,
			InitialStock: input.Stock.InitialStock,
			Description: input.Description,
			MinimumStock: input.Stock.MinimumStock,
		}, authId, issuedAt)
		if err != nil {
			return  ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return  ctx.Status(200).JSON(fiber.Map{"message": "product was created", "id": uuid})
	}
}



func NewProductHandler(createProductUC *product.CreateProductUseCase) *ProductHandler {
	return &ProductHandler{
		createProductUC: createProductUC,
	}
}