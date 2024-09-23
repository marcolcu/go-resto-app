package entity

import (
	"time"
)

type Menu struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" form:"name" validate:"required" gorm:"not null"`
	Description string    `json:"description" form:"description" validate:"required" gorm:"not null"`
	Price       float64   `json:"price" form:"price" validate:"required,gt=0" gorm:"not null"`
	Category    string    `json:"category" form:"category" validate:"required" gorm:"not null"`
	Image       string    `json:"image" form:"image" validate:"required" gorm:"not null"`
	Stock       int       `json:"stock" form:"stock" validate:"required,gte=0" gorm:"not null"`
	Disable     bool      `json:"disable" form:"disable" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
