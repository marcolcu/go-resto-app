package routers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/marcolcu/go-resto-app/controllers"
	"github.com/marcolcu/go-resto-app/middleware"
)

func RouterApp(app *fiber.App) {
	// Konfigurasi CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://demo-resto-app.vercel.app/",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Middleware Helmet - untuk perlindungan dari beberapa ancaman keamanan HTTP
	app.Use(helmet.New())

	// Middleware Rate Limiting - Membatasi jumlah request dari IP dalam interval tertentu
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	// Route untuk register (tanpa middleware)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/register", controllers.Register)
	app.Get("/api/microsites", controllers.GetMicrosite)
	app.Get("/api/about", controllers.GetAbout)
	app.Get("/api/menus", controllers.GetMenu)
	app.Post("/api/reservations", controllers.CreateReservation)
	app.Get("/api/reservations", controllers.GetReservations)
	app.Post("/api/testimonial", controllers.CreateTestimoni)
	app.Get("/api/testimonial", controllers.UpdateTestimonialActivity)

	// Grup route yang memerlukan otentikasi
	protected := app.Group("/api", middleware.AuthMiddleware)

	// Route dalam grup yang memerlukan otentikasi
	protected.Get("/users", controllers.UserControllerShow)
	protected.Post("/users/update", controllers.UserControllerUpdate)
	protected.Post("/users/delete", controllers.UserControllerDelete)

	protected.Post("/microsites", controllers.CreateMicrosite)
	protected.Post("/microsites/update", controllers.UpdateMicrosite)
	protected.Post("/microsites/delete", controllers.DeleteMicrosite)

	protected.Post("/menus", controllers.CreateMenu)
	protected.Post("/menus/update", controllers.UpdateMenu)
	protected.Post("/menus/delete", controllers.DeleteMenu)

	protected.Post("/about", controllers.CreateAbout)
	protected.Post("/about/update", controllers.UpdateAbout)
	protected.Post("/about/delete", controllers.DeleteAbout)

	protected.Post("/reservations/update", controllers.UpdateReservation)
	protected.Post("/reservations/delete", controllers.DeleteReservation)

	protected.Get("/all-testimonial", controllers.GetAllTestimoni)
}
