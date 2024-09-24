package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunAboutMigrate() {
	// Cek apakah tabel about sudah berisi data
	var count int64
	database.DB.Table("abouts").Count(&count)

	if count == 0 {
		// Jika tidak ada data, lakukan migrasi
		err := database.DB.AutoMigrate(&entity.About{}, &entity.Chef{})
		if err != nil {
			panic(err)
		}

		fmt.Println("About migration completed successfully")
	} else {
		// Jika sudah ada data, lewati migrasi
		fmt.Println("Migration skipped: about table already contains data")
	}
}
