package entity

import "time"

type MonthlyTransaction struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Month     string    `json:"month" gorm:"type:varchar(7);index"` // Format: "YYYY-MM"
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
