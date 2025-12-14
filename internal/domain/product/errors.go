package product

import "errors"

var (
    ErrProductNotFound = errors.New("product not found")
    ErrProductInactive = errors.New("product is inactive")
    ErrInvalidPrice    = errors.New("invalid product price")
    ErrInvalidCategory = errors.New("invalid product category")
    ErrInvalidName = errors.New("invalid product name")
    ErrInvalidDescription = errors.New(("invalid product description"))
)

