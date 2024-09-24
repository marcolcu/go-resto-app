package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

// CreateAbout creates a new About record
func CreateAbout(c *fiber.Ctx) error {
	// Parse JSON request body into the About struct
	about := new(entity.About)
	if err := c.BodyParser(about); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Check if tipe_section is unique
	var existingAbout entity.About
	if err := database.DB.Where("tipe_section = ?", about.TipeSection).First(&existingAbout).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "The 'Type' must be unique. Another record with the same 'Type' already exists.",
		})
	}

	// Set timestamps for the About record
	about.CreatedAt = time.Now()
	about.UpdatedAt = time.Now()

	// Save About to the database first to get its ID
	if err := database.DB.Create(&about).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create about record",
			"error":   err.Error(),
		})
	}

	// Handle chefs if tipe_section is 'chef'
	if about.TipeSection == "chef" {
		// Validate that chefs are provided
		if about.Chefs == nil || len(about.Chefs) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Chef details are required.",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "About and Chefs created successfully",
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
		if err := database.DB.Preload("Chefs").First(&about, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "About record not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success Get About",
			"about":   about,
			"chefs":   about.Chefs, // Include chefs in the response
		})
	}

	// Retrieve all about records if no ID is provided
	var abouts []entity.About
	if err := database.DB.Preload("Chefs").Find(&abouts).Error; err != nil {
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

	// Find the existing about record with associated chefs
	if err := database.DB.Preload("Chefs").First(&about, id).Error; err != nil {
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

	// Update chefs if they are provided in the request
	if about.Chefs != nil && len(about.Chefs) > 0 {
		for _, chef := range about.Chefs {
			// Check if chef exists; if yes, update, otherwise create
			var existingChef entity.Chef
			if err := database.DB.First(&existingChef, chef.ID).Error; err == nil {
				// Update existing chef
				existingChef.ChefName = chef.ChefName
				existingChef.ChefPosition = chef.ChefPosition
				existingChef.ChefImageURL = chef.ChefImageURL
				existingChef.UpdatedAt = time.Now()
				if err := database.DB.Save(&existingChef).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": "Failed to update chef record",
						"error":   err.Error(),
					})
				}
			} else {
				// Create new chef if it doesn't exist
				chef.AboutID = about.ID
				chef.CreatedAt = time.Now()
				chef.UpdatedAt = time.Now()
				if err := database.DB.Create(&chef).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"message": "Failed to create chef record",
						"error":   err.Error(),
					})
				}
			}
		}
	}

	// Preload chefs for the response
	if err := database.DB.Model(&about).Association("Chefs").Find(&about.Chefs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve chefs",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "About updated successfully",
		"about":   about,
		"chefs":   about.Chefs, // Include chefs in the response
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
		"message":       "About record deleted successfully",
		"deleted_about": about, // Optionally include the deleted record
	})
}
