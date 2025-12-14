package postgres

import (
	"errors"

	"github.com/google/uuid"
	domainStock "github.com/leonardo849/product_supermarket/internal/domain/stock"
	"gorm.io/gorm"
)

type StockRepository struct {
	db *gorm.DB
}

func toDomainStock(model *StockModel) *domainStock.Stock {
	return &domainStock.Stock{
		ID:        model.ID,
		ProductID: model.ProductID.String(),
		Quantity:  model.Quantity,
		Minimum:   model.Minimum,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}


func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) WithTx(tx *gorm.DB) *StockRepository {
	return &StockRepository{db: tx}
}


func (r *StockRepository) Create(stock *domainStock.Stock, productID uuid.UUID) error {
	model := StockModel{
		ID:        stock.ID,
		ProductID: productID,
		Quantity:  stock.Quantity,
		Minimum:   stock.Minimum,
		CreatedAt: stock.CreatedAt,
		UpdatedAt: stock.UpdatedAt,
	}

	return r.db.Create(&model).Error
}

func (r *StockRepository) FindByProductID(productID uuid.UUID) (*domainStock.Stock, error) {
	var model StockModel

	err := r.db.
		Where("product_id = ?", productID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainStock.ErrStockNotFound
		}
		return nil, err
	}

	return toDomainStock(&model), nil
}
