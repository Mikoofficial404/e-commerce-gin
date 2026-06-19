package postgres

import (
	"context"
	"ecommerce-gin/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	dbGorm *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{dbGorm: db}
}

func (r *OrderRepository) CreateOrder(order *domain.Order, items []*domain.OrderItem) error {
	return r.dbGorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, item := range items {
			item.OrderID = order.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OrderRepository) UpdateStatus(orderID string, newStatus string) error {
	ctx := context.Background()

	err := r.dbGorm.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("status", newStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) FindById(id string) (*domain.Order, error) {
	var result domain.Order
	err := r.dbGorm.Where("id = ?", id).First(&result).Error
	return &result, err
}
