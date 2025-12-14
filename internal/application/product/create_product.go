package product

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	domainProduct "github.com/leonardo849/product_supermarket/internal/domain/product"
	domainStock "github.com/leonardo849/product_supermarket/internal/domain/stock"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

type CreateProductInput struct {
    Name         string `validate:"required,min=3,max=100"`
    PriceInCents int64  `validate:"required,gt=0"`
    Category     string `validate:"required,oneof=FOOD DRINKS CLEANING HYGIENE OTHERS"`
    InitialStock int64  `validate:"gte=0"`
	Description string `validate:"required,min=10,max=300"`
	MinimumStock int64 `validate:"gte=0"`
}


type CreateProductUseCase struct {
    productRepo domainProduct.ProductRepository
    stockRepo   domainStock.StockRepository
	validator *validator.Validate
	uow         *postgres.UnitOfWork

}

func NewCreateProductUseCase(productRepo domainProduct.ProductRepository, stockRepo domainStock.StockRepository, uow *postgres.UnitOfWork) *CreateProductUseCase {
    return &CreateProductUseCase{
        productRepo: productRepo,
        stockRepo:   stockRepo,
		validator: validator.New(),
		uow: uow,
    }
}

func (uc *CreateProductUseCase) Execute(input CreateProductInput) (uuid.UUID, error) {
	if err := uc.validator.Struct(input); err != nil {
		return uuid.Nil, err
	}
	

	product, err := domainProduct.New(input.Name, input.Description, input.PriceInCents, input.Category)
	if err != nil {
		return uuid.Nil, err
	}

	// if err := uc.productRepo.Create(product); err != nil {
	// 	return uuid.Nil, err
	// }

	stock, err := domainStock.New(input.InitialStock, input.MinimumStock)
	if err != nil {
		return uuid.Nil, err
	}

	// if err := uc.stockRepo.Create(stock, product.ID); err != nil {
	// 	return  uuid.Nil, err
	// }	

	err = uc.uow.Do(func(tx *gorm.DB) error {
		pRep := uc.productRepo.(*postgres.ProductRepository).WithTx(tx)
		sRep := uc.stockRepo.(*postgres.StockRepository).WithTx(tx)
		
		if err := pRep.Create(product); err != nil {
			return err
		}

		if err := sRep.Create(stock, product.ID); err != nil {
			return err
		}

		return  nil
	})

	if err != nil {
		return  uuid.Nil, err
	}

	return  product.ID, nil

}