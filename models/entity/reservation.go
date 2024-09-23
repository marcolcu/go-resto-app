package entity

import (
	"time"
)

type Reservation struct {
    ID                 int                  `gorm:"primaryKey;autoIncrement" json:"id"`
    Name               string               `gorm:"type:varchar(255);not null" json:"name"`
    Email              string               `gorm:"type:varchar(255);not null" json:"email"`
    Phone              string               `gorm:"type:varchar(20);not null" json:"phone"`
    Guest              int                  `gorm:"not null" json:"guest"`
    ReserveTime        time.Time            `gorm:"not null" json:"reserve_time"`
    CreatedAt          time.Time            `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt          time.Time            `gorm:"autoUpdateTime" json:"updated_at"`
    ReservationDetails []Reservation_Detail `gorm:"foreignKey:ReservationId" json:"reservation_details"`
}

type Reservation_Detail struct {
	ReservationId int       `json:"reservation_id"`
	MenuId        int       `json:"menu_id"`
	Quantity      int       `json:"quantity"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
