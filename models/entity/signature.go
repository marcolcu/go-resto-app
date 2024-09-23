package entity

import (
	"time"
)

type Signature struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" form:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" form:"description" gorm:"type:text;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
