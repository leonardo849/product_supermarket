package product

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID uuid.UUID 
	Name string 
	Description string 
	Price int64 
	Category string 
	Active bool
	CreatedAt time.Time 
	UpdatedAt time.Time 
}

func New(Name string, Description string, Price int64, category string) (*Product, error){
	if Price <= 0 {
		return nil, ErrInvalidPrice
	}
	if !Category(category).isValid() {
		return nil, ErrInvalidCategory
	}
	if strings.TrimSpace(Name) == "" {
        return nil, ErrInvalidName
    }

    if len(strings.TrimSpace(Description)) < 10 {
        return nil, ErrInvalidDescription
    }


	now := time.Now()
	return  &Product{
		ID: uuid.New(),
		Name: Name,
		Description: Description,
		Price: Price,
		Category: category,
		CreatedAt: now,
		Active: true,
		UpdatedAt: now,
	}, nil
}