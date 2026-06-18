package postgres

import (
	"context"
	"ecommerce-gin/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	dbGorm *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{dbGorm: db}
}

func (r *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	result := r.dbGorm.Create(user)
	err := result.Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email *domain.User) (*domain.User, error) {
	ctx := context.Background()
	var result domain.User

	err := r.dbGorm.WithContext(ctx).
		Where("email = ?", email.Email).
		First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
