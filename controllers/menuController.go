package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

// Create Menu
func CreateMenu(c *fiber.Ctx) error {
	menu := new(entity.Menu)

	// Parsing body
	if err := c.BodyParser(menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Set timestamps
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()

	// Save to database
	if err := database.DB.Create(&menu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create menu",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Menu created successfully",
		"menu":    menu,
	})
}

// Get Menu by ID
func GetMenu(c *fiber.Ctx) error {
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

		// Find menu by ID
		var menu entity.Menu
		if err := database.DB.First(&menu, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Menu not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success Get Menu",
			"menu":    menu,
		})
	}

	// Retrieve all menus if no ID is provided
	var menus []entity.Menu
	if err := database.DB.Find(&menus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve menus",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Get Menus",
		"menus":   menus,
	})
}

// Update Menu
func UpdateMenu(c *fiber.Ctx) error {
	// Retrieve ID from query parameters
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required",
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

	var menu entity.Menu

	// Find the existing menu
	if err := database.DB.First(&menu, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Menu not found",
			"error":   err.Error(),
		})
	}

	// Parse request body
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Update timestamps
	menu.UpdatedAt = time.Now()

	// Save updated menu to database
	if err := database.DB.Save(&menu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update menu",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu updated successfully",
		"menu":    menu,
	})
}

// Delete Menu
func DeleteMenu(c *fiber.Ctx) error {
	// Retrieve ID from query parameters
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required",
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

	var menu entity.Menu

	// Find menu by ID
	if err := database.DB.First(&menu, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Menu not found",
			"error":   err.Error(),
		})
	}

	// Delete the menu
	if err := database.DB.Delete(&menu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete menu",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu deleted successfully",
	})
}
