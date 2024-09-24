package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunUserMigrate() {
	// Check if the 'users' table already has data
	var count int64
	database.DB.Table("users").Count(&count)

	if count == 0 {
		// If no data exists, perform the migration
		err := database.DB.AutoMigrate(&entity.User{})
		if err != nil {
			panic(err)
		}

		fmt.Println("User migration completed successfully")
	} else {
		// Skip migration if data exists
		fmt.Println("Migration skipped: users table already contains data")
	}

	// Optional: Modify the email column to `longtext`
	err := database.DB.Exec("ALTER TABLE users MODIFY COLUMN email LONGTEXT NOT NULL").Error
	if err != nil {
		fmt.Println("Error modifying column:", err.Error())
	} else {
		fmt.Println("Email column modified successfully")
	}
}
