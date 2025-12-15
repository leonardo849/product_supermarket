package postgres

import (
	"time"

	"github.com/google/uuid"
)

type ProductModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	Price       int64
	Category    string
	Active      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}


type StockModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;constraint:OnUpdate:OnDelete:CASCADE;"`
	Quantity int64
	Minimum  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	AuthId string
	Role string 
	CreatedAt time.Time
	UpdatedAt time.Time
}