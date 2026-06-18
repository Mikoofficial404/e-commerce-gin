package domain

import (
	"time"
)

type User struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email       string    `json:"email" gorm:"column:email"`
	Password    string    `json:"password" gorm:"column:password"`
	FullName    string    `json:"full_name" gorm:"column:full_name"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number"`
	Role        string    `json:"role" gorm:"column:role"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}
