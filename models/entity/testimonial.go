package entity

import "time"

type Testimoni struct {
	ID          uint      `json:"id" gorm:"primary_key;autoIncrement"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Customer    string    `json:"customer" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Active      bool      `json:"active" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
