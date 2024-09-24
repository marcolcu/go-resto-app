package entity

import "time"

type About struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	TipeSection string    `json:"tipe_section"`
	Chefs       []Chef    `json:"chefs"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Chef struct {
	ID           uint      `json:"id"`
	ChefName     string    `json:"chef_name"`
	ChefPosition string    `json:"chef_position"`
	ChefImageURL string    `json:"chef_image_url"`
	AboutID      uint      `json:"about_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
