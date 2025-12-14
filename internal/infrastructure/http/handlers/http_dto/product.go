package http_dto

type CreateProductDTOHttp struct {
	Name         string `json:"name"`
    PriceInCents int64  `json:"price"`
    Category     string `json:"category"`
	Description string `json:"description"`
	Stock struct {
		InitialStock int64  `json:"inital_stock"`
		MinimumStock int64 `json:"minimum_stock"`
	}
}