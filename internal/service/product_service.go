package service

import (
	"context"
	"ecommerce-gin/internal/domain"
	"ecommerce-gin/internal/repository/postgres"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	product     postgres.ProductRepository
	redisClient *redis.Client
}

func NewProductService(repo *postgres.ProductRepository, redisClient *redis.Client) *ProductService {
	return &ProductService{product: *repo, redisClient: redisClient}
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
	ctx := context.Background()
	s.redisClient.FlushDB(ctx)
	return result, nil
}

func (s *ProductService) FindAllProduct(page, limit int, search string) ([]*domain.Product, int64, error) {
	ctx := context.Background()

	cacheKey := "products_page_" + strconv.Itoa(page) + "_limit_" + strconv.Itoa(limit) + "_search_" + search

	cachedData, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cacheResult struct {
			Products []*domain.Product
			Total    int64
		}
		json.Unmarshal([]byte(cachedData), &cacheResult)
		return cacheResult.Products, cacheResult.Total, nil
	}

	products, total, err := s.product.FindAllProduct(page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	cacheResult := struct {
		Products []*domain.Product
		Total    int64
	}{
		Products: products,
		Total:    total,
	}
	jsonData, _ := json.Marshal(cacheResult)
	s.redisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	return products, total, nil
}
