package stock

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID uuid.UUID 
	Quantity int64 
	ProductID string
	Minimum int64 
	CreatedAt time.Time 
	UpdatedAt time.Time 
}

func New(Quantity int64, Minimum int64) (*Stock, error) {
	now := time.Now().UTC()
	if Quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if Minimum < 0 {
		return nil, ErrInvalidQuantity
	}
	if Minimum > Quantity {
		return nil, ErrInvalidMinimum
	}
	return  &Stock{
		ID: uuid.New(),
		Quantity: Quantity,
		Minimum: Minimum,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}