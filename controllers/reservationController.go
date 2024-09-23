package controllers

import (
	"fmt"
	"time"

	"github.com/Fabian832/Go-Fiber/database"
	"github.com/Fabian832/Go-Fiber/models/entity"
	"github.com/gofiber/fiber/v2"
)

// CreateReservation handles creation of a reservation
func CreateReservation(c *fiber.Ctx) error {
	type MenuInput struct {
		MenuId   int `json:"menu_id"`
		Quantity int `json:"quantity"`
	}

	type ReservationInput struct {
		Name        string      `json:"name"`
		Email       string      `json:"email"`
		Phone       string      `json:"phone"`
		Guest       int         `json:"guest"` // Changed from Number to Guest
		ReserveTime string      `json:"reserve_time"`
		Menus       []MenuInput `json:"menus"`
	}

	input := new(ReservationInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Time validator
	fmt.Println("Received ReserveTime:", input.ReserveTime)

	if input.ReserveTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ReserveTime cannot be empty",
		})
	}

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, input.ReserveTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date format",
			"error":   err.Error(),
		})
	}

	// Guest Validator
	if input.Guest <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Guest count must be greater than zero",
		})
	}

	reservation := &entity.Reservation{
		Name:        input.Name,
		Email:       input.Email,
		Phone:       input.Phone,
		Guest:       input.Guest,
		ReserveTime: parsedTime,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&reservation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create reservation",
			"error":   err.Error(),
		})
	}

	for _, menu := range input.Menus {
		// Create a new reservation_detail entry for each menu
		reservationDetails := &entity.Reservation_Detail{
			ReservationId: reservation.ID,
			MenuId:        menu.MenuId,
			Quantity:      menu.Quantity,
		}

		var menuEntity entity.Menu
		if err := database.DB.First(&menuEntity, reservationDetails.MenuId).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Menu not found",
				"error":   err.Error(),
			})
		}

		if menuEntity.Stock < reservationDetails.Quantity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Not enough stock for " + menuEntity.Name,
			})
		}

		menuEntity.Stock -= reservationDetails.Quantity

		if err := database.DB.Save(&menuEntity).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update menu stock",
				"error":   err.Error(),
			})
		}

		// Save the reservation_detail to the database
		if err := database.DB.Create(&reservationDetails).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save reservation details",
				"error":   err.Error(),
			})
		}
	}

	if err := database.DB.Preload("ReservationDetails").First(&reservation, reservation.ID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Reservation not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Reservation created successfully",
		"reservation": reservation,
	})
}

// GetReservations retrieves all reservations or a single reservation by ID
func GetReservations(c *fiber.Ctx) error {
	id := c.Query("id")

	if id != "" {
		// If an id is provided, get the reservation by ID
		var reservation entity.Reservation

		// Use Preload to load associated reservation details
		if err := database.DB.Preload("ReservationDetails").First(&reservation, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Reservation not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":     "Success Get Reservation",
			"reservation": reservation,
		})
	}

	// If no id is provided, get all reservations
	var reservations []entity.Reservation

	// Use Preload to load associated reservation details
	if err := database.DB.Preload("ReservationDetails").Find(&reservations).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Reservations not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Success Get All Reservations",
		"reservations": reservations,
	})
}

// UpdateReservation handles updating a reservation
func UpdateReservation(c *fiber.Ctx) error {
	id := c.Query("id")
	var reservation entity.Reservation

	if err := database.DB.First(&reservation, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Reservation not found",
			"error":   err.Error(),
		})
	}

	type MenuInputDetail struct {
		MenuId   int `json:"menu_id"`
		Quantity int `json:"quantity"`
	}

	type ReservationInput struct {
		Name        string            `json:"name"`
		Email       string            `json:"email"`
		Phone       string            `json:"phone"`
		Guest       int               `json:"guest"`
		ReserveTime string            `json:"reserve_time"`
		Menus       []MenuInputDetail `json:"menus"`
	}

	input := new(ReservationInput)

	// Time validator
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	layout := "2006-01-02 15:04:05"
	if input.ReserveTime != "" {
		parsedTime, err := time.Parse(layout, input.ReserveTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid date format",
				"error":   err.Error(),
			})
		}
		reservation.ReserveTime = parsedTime
	}

	// Guest Validator
	if input.Guest <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Guest count must be greater than zero",
		})
	}

	// Update other fields
	reservation.Name = input.Name
	reservation.Email = input.Email
	reservation.Phone = input.Phone
	reservation.Guest = input.Guest
	reservation.UpdatedAt = time.Now()

	if err := database.DB.Save(&reservation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update reservation",
			"error":   err.Error(),
		})
	}

	// Delete existing reservation details
	if err := database.DB.Where("reservation_id = ?", reservation.ID).Delete(&entity.Reservation_Detail{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete existing reservation details",
			"error":   err.Error(),
		})
	}

	// Add new reservation details
	for _, menu := range input.Menus {
		reservationDetails := &entity.Reservation_Detail{
			ReservationId: reservation.ID,
			MenuId:        menu.MenuId,
			Quantity:      menu.Quantity,
		}

		var menuEntity entity.Menu
		if err := database.DB.First(&menuEntity, reservationDetails.MenuId).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Menu not found",
				"error":   err.Error(),
			})
		}

		if menuEntity.Stock < reservationDetails.Quantity {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Not enough stock for " + menuEntity.Name,
			})
		}

		menuEntity.Stock -= reservationDetails.Quantity

		if err := database.DB.Save(&menuEntity).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update menu stock",
				"error":   err.Error(),
			})
		}

		if err := database.DB.Create(&reservationDetails).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create reservation details",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Reservation updated successfully",
		"reservation": reservation,
	})
}

// DeleteReservation handles deleting a reservation
func DeleteReservation(c *fiber.Ctx) error {
	id := c.Query("id")
	var reservation entity.Reservation

	if err := database.DB.First(&reservation, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Reservation not found",
			"error":   err.Error(),
		})
	}

	// Delete related reservation details
	if err := database.DB.Where("reservation_id = ?", reservation.ID).Delete(&entity.Reservation_Detail{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete reservation details",
			"error":   err.Error(),
		})
	}

	// Delete reservation
	if err := database.DB.Delete(&reservation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete reservation",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reservation deleted successfully",
	})
}
