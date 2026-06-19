package service

import (
	"context"
	"ecommerce-gin/internal/domain"
	"ecommerce-gin/internal/repository/postgres"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	product     postgres.ProductRepository
	redisClient *redis.Client
	s3Client    *s3.Client
}

func NewProductService(repo *postgres.ProductRepository, redisClient *redis.Client, s3Client *s3.Client) *ProductService {
	return &ProductService{product: *repo, redisClient: redisClient, s3Client: s3Client}
}

func (s *ProductService) CreateProduct(name string, description string, price float64, stock int, file *multipart.FileHeader) (*domain.Product, error) {

	fileOpens, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer fileOpens.Close()
	ctx := context.Background()
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("e-commerce-bucket"),
		Key:    aws.String(file.Filename),
		Body:   fileOpens,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal upload ke S3: %w", err)
	}
	product := domain.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    "http://localhost:9000/e-commerce-bucket/" + file.Filename,
	}
	result, err := s.product.CreateProduct(&product)
	if err != nil {
		return nil, err
	}

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
