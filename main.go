package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/database/migration"
	"github.com/marcolcu/go-resto-app/routers"
)

func main() {
	database.ConnectDB()
	migration.RunUserMigrate()
	migration.RunMicrositeMigrate()
	migration.RunMenuMigrate()
	migration.RunReservationMigrate()
	migration.RunReservationDetailMigrate()
	migration.RunSignatureMigrate()
	migration.RunAboutMigrate()
	migration.RunTestimonialMigrate()
	migration.RunMonthlyTransactionMigrate()
	app := fiber.New()

	routers.RouterApp(app)

	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
