package product

import (
	"github.com/google/uuid"
	
)

type ProductRepository interface {
    Create(product *Product) error
    Update(id uuid.UUID, product *Product) error
    FindByID(id uuid.UUID) (*Product, error)
    FindActiveByID(id uuid.UUID) (*Product, error)
    ExistsByID(id uuid.UUID) (bool, error)
    
}
