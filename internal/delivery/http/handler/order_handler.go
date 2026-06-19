package handler

import (
	"ecommerce-gin/internal/dto/request"
	"ecommerce-gin/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

type XenditCallback struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: *svc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not logged in!"})
		return
	}
	userIDStr := fmt.Sprintf("%v", userId)
	var input request.OrderRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invoiceUrl, err := h.orderService.CreateOrder(userIDStr, input.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "Order Created", "payment_url": invoiceUrl})
}

func (h *OrderHandler) WebhookPayment(c *gin.Context) {
	var callback XenditCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed ",
			"message": err.Error(),
		})
		return
	}
	orderId := callback.ExternalID
	err := h.orderService.PayOrder(orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process payment",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Payment processed successfully",
		"order_id": orderId,
	})
}
