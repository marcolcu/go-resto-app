package controllers

import (
	"strconv"
	"time"

	"github.com/Fabian832/Go-Fiber/database"
	"github.com/Fabian832/Go-Fiber/models/entity"
	"github.com/gofiber/fiber/v2"
)

// CreateAbout creates a new About record
func CreateAbout(c *fiber.Ctx) error {
	about := new(entity.About)

	// Parse JSON request body into the About struct
	if err := c.BodyParser(about); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Check if the tipe_section is unique
	var existingAbout entity.About
	if err := database.DB.Where("tipe_section = ?", about.TipeSection).First(&existingAbout).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "The 'Type' must be unique. Another record with the same 'Type' already exists.",
		})
	}

	// Set timestamps
	about.CreatedAt = time.Now()
	about.UpdatedAt = time.Now()

	// Save to database
	if err := database.DB.Create(&about).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create about record",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "About created successfully",
		"about":   about,
	})
}

// GetAbout retrieves About records (either all or by ID)
func GetAbout(c *fiber.Ctx) error {
	// Retrieve ID from query parameters
	idStr := c.Query("id")

	if idStr != "" {
		// Convert ID to integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid ID",
				"error":   err.Error(),
			})
		}

		// Find about record by ID
		var about entity.About
		if err := database.DB.First(&about, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "About record not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success Get About",
			"about":   about,
		})
	}

	// Retrieve all about records if no ID is provided
	var abouts []entity.About
	if err := database.DB.Find(&abouts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve about records",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Get Abouts",
		"abouts":  abouts,
	})
}

// UpdateAbout updates an existing About record by ID
func UpdateAbout(c *fiber.Ctx) error {
	// Retrieve ID from query parameters
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required to update the about record",
		})
	}

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"error":   err.Error(),
		})
	}

	var about entity.About

	// Find the existing about record
	if err := database.DB.First(&about, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "About record not found",
			"error":   err.Error(),
		})
	}

	// Parse the updated data from the request body
	if err := c.BodyParser(&about); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Update timestamps
	about.UpdatedAt = time.Now()

	// Save updated about record to database
	if err := database.DB.Save(&about).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update about record",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "About updated successfully",
		"about":   about,
	})
}

// DeleteAbout deletes an About record by ID
func DeleteAbout(c *fiber.Ctx) error {
	// Retrieve ID from query parameters
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required to delete the about record",
		})
	}

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"error":   err.Error(),
		})
	}

	var about entity.About

	// Find about record by ID
	if err := database.DB.First(&about, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "About record not found",
			"error":   err.Error(),
		})
	}

	// Delete the about record
	if err := database.DB.Delete(&about).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete about record",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "About record deleted successfully",
	})
}
