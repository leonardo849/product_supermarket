package postgres

import (
	"errors"

	"github.com/google/uuid"
	domainProduct "github.com/leonardo849/product_supermarket/internal/domain/product"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func toDomainProduct(model *ProductModel) *domainProduct.Product {
	return &domainProduct.Product{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price,
		Category:    model.Category,
		Active:      model.Active,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (r *ProductRepository) WithTx(tx *gorm.DB) *ProductRepository {
	return &ProductRepository{db: tx}
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *domainProduct.Product) error {
	model := ProductModel{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		Active:      product.Active,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	return r.db.Create(&model).Error
}

func (r *ProductRepository) FindByID(id uuid.UUID) (*domainProduct.Product, error) {
	var model ProductModel

	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainProduct.ErrProductNotFound
		}
		return nil, err
	}

	return toDomainProduct(&model), nil
}


func (r *ProductRepository) ExistsByID(id uuid.UUID) (bool, error) {
	var count int64

	err := r.db.Model(&ProductModel{}).
		Where("id = ?", id).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ProductRepository) FindActiveByID(id uuid.UUID) (*domainProduct.Product, error) {
	var model ProductModel

	err := r.db.
		Where("id = ? AND active = true", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainProduct.ErrProductNotFound
		}
		return nil, err
	}

	return toDomainProduct(&model), nil
}

func (r *ProductRepository) Update(id uuid.UUID,product *domainProduct.Product) error {
	updates := make(map[string]interface{})

	if product.Name != "" {
		updates["name"] = product.Name
	}

	if product.Description != "" {
		updates["description"] = product.Description
	}

	if product.Price != 0 {
		updates["price"] = product.Price
	}

	if product.Category != "" {
		updates["category"] = product.Category
	}


	result := r.db.
		Model(&ProductModel{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainProduct.ErrProductNotFound
	}

	return nil
}

