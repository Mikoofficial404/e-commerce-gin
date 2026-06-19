package service

import (
	"context"
	"ecommerce-gin/internal/domain"
	"ecommerce-gin/internal/dto/request"
	"ecommerce-gin/internal/pkg/rabbitmq"
	"ecommerce-gin/internal/repository/postgres"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	xendit "github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

type OrderService struct {
	orderRepo   postgres.OrderRepository
	productRepo postgres.ProductRepository
	userRepo    postgres.UserRepository
	rabbitConn  *amqp.Connection
}

func NewOrderService(orderRepo *postgres.OrderRepository, productRepo *postgres.ProductRepository, userRepo *postgres.UserRepository, rabbitConn *amqp.Connection) *OrderService {
	return &OrderService{orderRepo: *orderRepo, productRepo: *productRepo, userRepo: *userRepo, rabbitConn: rabbitConn}
}

func (s *OrderService) CreateOrder(userID string, items []request.OrderItemRequest) (string, error) {
	var totalPrice float64
	var orderItems []*domain.OrderItem
	for _, item := range items {
		product, err := s.productRepo.FindById(item.ProductID)
		if err != nil {
			return "", err
		}

		if product.Stock < item.Quantity {
			return "", fmt.Errorf("insufficient stock for product %s", product.Name)
		}
		product.Stock -= item.Quantity

		err = s.productRepo.UpdateStok(product.ID, product.Stock)
		if err != nil {
			return "", err
		}
		subTotal := int(product.Price) * item.Quantity
		totalPrice += float64(subTotal)
		orderItems = append(orderItems, &domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     float64(subTotal),
		})
	}

	order := domain.Order{
		UserID:     userID,
		Status:     "PENDING",
		TotalPrice: totalPrice,
	}

	err := s.orderRepo.CreateOrder(&order, orderItems)
	if err != nil {
		return "", err
	}

	keyXendit := os.Getenv("XENDIT_SECRET_KEY")
	xenditClient := xendit.NewClient(keyXendit)
	req := invoice.CreateInvoiceRequest{
		ExternalId: order.ID,
		Amount:     totalPrice,
	}
	resp, _, errXendit := xenditClient.InvoiceApi.CreateInvoice(context.Background()).CreateInvoiceRequest(req).Execute()
	if errXendit != nil {
		return "", errXendit
	}

	fmt.Fprintf(os.Stdout, "Response from `InvoiceApi.GetInvoiceById`: %v\n", resp)
	return resp.InvoiceUrl, nil
}

func (s *OrderService) PayOrder(orderID string) error {
	err := s.orderRepo.UpdateStatus(orderID, "PAID")
	if err != nil {
		return err
	}

	order, err := s.orderRepo.FindById(orderID)
	if err == nil {
		user, err := s.userRepo.FindById(order.UserID)
		if err == nil {
			rabbitmq.PublishMessage(s.rabbitConn, "invoice_queue", user.Email)
		}
	}
	return nil
}
