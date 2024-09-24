package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

func UserControllerShow(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr != "" {
		// If ID is provided, fetch the user by ID
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid ID",
				"error":   err.Error(),
			})
		}

		var user entity.User
		if err := database.DB.First(&user, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success Get User",
			"user":    user,
		})
	} else {
		// If no ID is provided, fetch all users
		var users []entity.User
		if err := database.DB.Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve users",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success Get All Users",
			"users":   users,
		})
	}
}

func UserControllerUpdate(c *fiber.Ctx) error {
	// Parse the request body
	user := new(entity.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

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

	// Update the user record
	if err := database.DB.Model(&entity.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success update user",
	})
}

func UserControllerDelete(c *fiber.Ctx) error {
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

	// Delete the user record
	if err := database.DB.Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
