package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunReservationDetailMigrate() {
	// Cek apakah tabel reservation sudah berisi data
	var count int64
	database.DB.Table("reservation_details").Count(&count)

	if count == 0 {
		// Jika tidak ada data, lakukan migrasi
		err := database.DB.AutoMigrate(&entity.Reservation_Detail{})
		if err != nil {
			panic(err)
		}

		fmt.Println("Reservation migration completed successfully")
	} else {
		// Jika sudah ada data, lewati migrasi
		fmt.Println("Migration skipped: reservation table already contains data")
	}
}
