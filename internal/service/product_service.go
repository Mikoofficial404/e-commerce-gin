package service

import (
	"ecommerce-gin/internal/domain"
	"ecommerce-gin/internal/repository/postgres"
)

type ProductService struct {
	product postgres.ProductRepository
}

func NewProductService(repo *postgres.ProductRepository) *ProductService {
	return &ProductService{product: *repo}
}

func (s *ProductService) CreateProduct(name string, description string, price float64, stock int, imageURL string) (*domain.Product, error) {
	product := domain.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    imageURL,
	}
	result, err := s.product.CreateProduct(&product)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProductService) FindAllProduct() ([]*domain.Product, error) {
	result, err := s.product.FindAllProduct()
	if err != nil {
		return nil, err
	}
	return result, nil
}
