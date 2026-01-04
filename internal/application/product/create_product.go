package product

import (
	// "time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/leonardo849/product_supermarket/internal/application/common"

	// eventsUser "github.com/leonardo849/product_supermarket/internal/domain/events/user"
	// eventsProduct "github.com/leonardo849/product_supermarket/internal/domain/events/product"
	applicationUser "github.com/leonardo849/product_supermarket/internal/application/user"
	domainProduct "github.com/leonardo849/product_supermarket/internal/domain/product"
	domainStock "github.com/leonardo849/product_supermarket/internal/domain/stock"
	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
	"github.com/leonardo849/product_supermarket/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

type CreateProductInput struct {
	Name         string `validate:"required,min=3,max=100"`
	PriceInCents int64  `validate:"required,gt=0"`
	Category     string `validate:"required,oneof=FOOD DRINKS CLEANING HYGIENE OTHERS"`
	InitialStock int64  `validate:"gte=0"`
	Description  string `validate:"required,min=10,max=300"`
	MinimumStock int64  `validate:"gte=0"`
}

type CreateProductUseCase struct {
	productRepo domainProduct.ProductRepository
	stockRepo   domainStock.StockRepository
	// userRepo    domainUser.UserRepository
	validator   *validator.Validate
	uow         *postgres.UnitOfWork
	publisher   common.EventPublisher
	findUserUc *applicationUser.FindUserUseByAuthIdCase
}

func NewCreateProductUseCase(productRepo domainProduct.ProductRepository, stockRepo domainStock.StockRepository, uow *postgres.UnitOfWork, findUserUc *applicationUser.FindUserUseByAuthIdCase, publisher common.EventPublisher) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepo: productRepo,
		stockRepo:   stockRepo,
		findUserUc: findUserUc,
		validator:   validator.New(),
		uow:         uow,
		publisher:   publisher,
	}
}

func (uc *CreateProductUseCase) Execute(input CreateProductInput, authId string, issuedAt float64) (uuid.UUID, error) {
	if err := uc.validator.Struct(input); err != nil {
		return uuid.Nil, err
	}

	user, err := uc.findUserUc.Execute(authId)
	if err != nil {
		return  uuid.Nil, err
	}

	if user.UserWasUpdatedAfterToken(issuedAt) {
		return uuid.Nil, domainUser.ErrUserWasUpdatedAfterToken
	}

	if !user.CanUserCreateOrEditAProduct() {
		return uuid.Nil, domainUser.ErrUserCannotCreateAProduct
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

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	// go func() {
	// 	event := eventsProduct.ProductCreated{
	// 		UserID: user.ID.String(),
	// 		ProductID: product.ID.String(),
	// 		ProductName: product.Name,
	// 		OccurredAt: time.Now(),
	// 	}
	// 	uc.publisher.Publish(event)
	// }()
	return product.ID, nil

}
