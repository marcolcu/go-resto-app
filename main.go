package main

import (
	"github.com/Fabian832/Go-Fiber/database"
	"github.com/Fabian832/Go-Fiber/database/migration"
	"github.com/Fabian832/Go-Fiber/routers"
	"github.com/gofiber/fiber/v2"
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
	app := fiber.New()

	routers.RouterApp(app)

	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
