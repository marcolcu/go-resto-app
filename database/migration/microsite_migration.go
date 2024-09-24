package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunMicrositeMigrate() {
	// Check if the 'microsite' table already exists
	if database.DB.Migrator().HasTable(&entity.Microsite{}) {
		fmt.Println("Microsite table already exists. Skipping migration.")
		return
	}

	// Run AutoMigrate for entity Microsite
	err := database.DB.AutoMigrate(&entity.Microsite{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Microsite table migrated successfully")
}
