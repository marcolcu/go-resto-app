package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunMenuMigrate() {
	// Cek apakah tabel menu sudah berisi data
	var count int64
	database.DB.Table("menus").Count(&count)

	if count == 0 {
		// Jika tidak ada data, lakukan migrasi
		err := database.DB.AutoMigrate(&entity.Menu{})
		if err != nil {
			panic(err)
		}

		fmt.Println("Menu migration completed successfully")
	} else {
		// Jika sudah ada data, lewati migrasi
		fmt.Println("Migration skipped: menu table already contains data")
	}
}
