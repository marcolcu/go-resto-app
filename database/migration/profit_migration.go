package migration

import (
	"fmt"

	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func RunMonthlyTransactionMigrate() {
	// Check if the table already contains data
	var count int64
	database.DB.Table("monthly_transactions").Count(&count)

	if count == 0 {
		// If no data exists, perform migration
		err := database.DB.AutoMigrate(&entity.MonthlyTransaction{})
		if err != nil {
			panic(err)
		}

		fmt.Println("MonthlyTransaction migration completed successfully")
	} else {
		// If data already exists, skip the migration
		fmt.Println("Migration skipped: monthly_transactions table already contains data")
	}
}
