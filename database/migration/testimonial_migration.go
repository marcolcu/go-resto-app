package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunTestimonialMigrate() {
	// Cek apakah tabel testimonials sudah berisi data
	var count int64
	database.DB.Table("testimonials").Count(&count)

	if count == 0 {
		// Jika tidak ada data, lakukan migrasi
		err := database.DB.AutoMigrate(&entity.Testimoni{})
		if err != nil {
			panic(err)
		}

		fmt.Println("Testimonial migration completed successfully")
	} else {
		// Jika sudah ada data, lewati migrasi
		fmt.Println("Migration skipped: testimonials table already contains data")
	}
}
