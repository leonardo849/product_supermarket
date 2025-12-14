package stock

import "github.com/google/uuid"

type StockRepository interface {
    Create(stock *Stock, productID uuid.UUID) error
    FindByProductID(productID uuid.UUID) (*Stock, error)
    // Increase(productID uuid.UUID, quantity int64) error
    // Decrease(productID uuid.UUID, quantity int64) error
}
