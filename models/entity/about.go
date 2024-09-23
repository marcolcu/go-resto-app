package entity

import (
	"time"
)

type About struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Description string    `gorm:"type:text;not null" json:"description"`
	TipeSection string    `gorm:"type:varchar(255);not null" json:"tipe_section"`
	Image       string    `gorm:"type:varchar(255)" json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
