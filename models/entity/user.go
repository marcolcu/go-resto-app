package entity

import (
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" form:"name" validate:"gte=6,lte=32" gorm:"not null"`
	Email     string    `json:"email" form:"email" validate:"required,email" gorm:"type:varchar(255);unique;not null"`
	Password  string    `json:"password" form:"password" validate:"required,gte=8" gorm:"not null"`
	Phone     string    `json:"phone" form:"phone" validate:"required,number,min=12" gorm:"not null"`
	IsActive  int       `json:"isActive" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
