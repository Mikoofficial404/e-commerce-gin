package domain

import "time"

type Product struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	Price       float64   `json:"price" gorm:"column:price"`
	Stock       int       `json:"stock" gorm:"column:stock"`
	ImageURL    string    `json:"imageUrl" gorm:"column:imageurl"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}
