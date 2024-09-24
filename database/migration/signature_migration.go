package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunSignatureMigrate() {
	// Check if the 'signature' table already exists
	if database.DB.Migrator().HasTable(&entity.Signature{}) {
		fmt.Println("Signature table already exists. Skipping migration.")
		return
	}

	// Run AutoMigrate for entity Signature
	err := database.DB.AutoMigrate(&entity.Signature{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Signature table migrated successfully")
}
