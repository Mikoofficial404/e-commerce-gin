package postgres

import (
	"context"
	"ecommerce-gin/internal/domain"

	"gorm.io/gorm"
)

type ProductRepository struct {
	dbGorm *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{dbGorm: db}
}

func (r *ProductRepository) CreateProduct(product *domain.Product) (*domain.Product, error) {
	isCreate := r.dbGorm.Create(product)
	err := isCreate.Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) FindAllProduct(page, limit int, search string) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	query := r.dbGorm.Model(&domain.Product{})
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepository) FindById(id string) (*domain.Product, error) {
	ctx := context.Background()
	var result domain.Product

	err := r.dbGorm.WithContext(ctx).
		Where("id = ?", id).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ProductRepository) UpdateStok(id string, newStock int) error {
	ctx := context.Background()

	err := r.dbGorm.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", id).
		Update("stock", newStock).Error
	if err != nil {
		return err
	}
	return nil
}
