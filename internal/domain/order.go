package domain

import "time"

type Order struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"user_id" gorm:"column:user_id"`
	Status     string    `json:"status" gorm:"column:status"`
	TotalPrice float64   `json:"total_price" gorm:"column:total_price"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type OrderItem struct {
	ID        string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID   string  `json:"order_id" gorm:"column:order_id"`
	ProductID string  `json:"product_id" gorm:"column:product_id"`
	Quantity  int     `json:"quantity" gorm:"column:quantity"`
	Price     float64 `json:"price" gorm:"column:price"`
}
