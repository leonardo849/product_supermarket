package stock

import "errors"

var (
    ErrStockNotFound      = errors.New("stock not found")
    ErrInsufficientStock = errors.New("insufficient stock")
    ErrInvalidQuantity   = errors.New("invalid stock quantity")
    ErrInvalidMinimum = errors.New("invalid minimum")
)
