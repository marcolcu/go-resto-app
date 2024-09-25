package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
	"time"
)

// CalculateTotalTransactionsCurrentMonth calculates total transactions for the current month without saving to the database
func CalculateTotalTransactionsCurrentMonth(c *fiber.Ctx) error {
	// Get the current year and month
	now := time.Now()
	currentMonth := now.Format("2006-01") // Format as "YYYY-MM"
	lastMonth := now.AddDate(0, -1, 0).Format("2006-01") // Last month

	// Parse start and end dates for the current month
	startDate, _ := time.Parse("2006-01", currentMonth)
	endDate := startDate.AddDate(0, 1, 0) // Add one month

	var reservations []entity.Reservation

	// Fetch all reservations within the current month
	if err := database.DB.Preload("ReservationDetails").
		Where("reserve_time >= ? AND reserve_time < ?", startDate, endDate).
		Find(&reservations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve reservations",
			"error":   err.Error(),
		})
	}

	// Calculate the total price for the current month
	var totalCurrentMonth float64
	for _, reservation := range reservations {
		for _, detail := range reservation.ReservationDetails {
			var menu entity.Menu
			if err := database.DB.First(&menu, detail.MenuId).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to retrieve menu details",
					"error":   err.Error(),
				})
			}
			totalCurrentMonth += float64(detail.Quantity) * menu.Price
		}
	}

	// Calculate total transactions for the last month
	var totalLastMonth float64
	if err := database.DB.Model(&entity.MonthlyTransaction{}).
		Where("month = ?", lastMonth).
		Select("total").
		Scan(&totalLastMonth).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve last month's total",
			"error":   err.Error(),
		})
	}

	// Check if there were transactions last month
	if totalLastMonth == 0 {
		totalLastMonth = 0 // Or leave it, and add null in JSON
	}

	// Count upcoming reservations
	var upcomingCount int64
	if err := database.DB.Model(&entity.Reservation{}).
		Where("reserve_time > ?", now).
		Count(&upcomingCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count upcoming reservations",
			"error":   err.Error(),
		})
	}

	// Create response JSON
	response := fiber.Map{
		"message":                "Total transactions calculated",
		"total_current_month":    totalCurrentMonth,
		"total_last_month":       totalLastMonth,
		"upcoming_reservations":  upcomingCount,
	}

	if totalLastMonth == 0 {
		response["total_last_month"] = nil // Set to null if there were no transactions
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
