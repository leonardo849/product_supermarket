package product

import "time"

type ProductCreated struct {
	UserID      string
	ProductID   string
	ProductName string
	OccurredAt  time.Time
}
