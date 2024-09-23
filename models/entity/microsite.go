package entity

import (
	"time"
)

type Microsite struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Content     string    `json:"content" form:"content" gorm:"type:text;not null"`
	Image       string    `json:"image" form:"image" gorm:"type:varchar(255);"`
	TipeSection string    `json:"tipe_section" form:"tipe_section" gorm:"type:varchar(50);not null"`
	Description string    `json:"description" form:"description" gorm:"type:text;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
