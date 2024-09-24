package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

// CreateMicrosite handles the creation of a new microsite
func CreateMicrosite(c *fiber.Ctx) error {
	// Define a struct to capture the request body including points
	type CreateMicrositeRequest struct {
		Content     string              `json:"content"`
		Description string              `json:"description"`
		Image       string              `json:"image"`
		TipeSection string              `json:"tipe_section"`
		Points      []map[string]string `json:"points"` // Capture points as a slice of maps
	}

	// Create an instance of the request struct
	var req CreateMicrositeRequest

	// Parse the request body into the struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Create and populate the microsite entity
	microsite := entity.Microsite{
		Content:     req.Content,
		Description: req.Description,
		Image:       req.Image,
		TipeSection: req.TipeSection,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Validate fields
	if microsite.Content == "" || microsite.Description == "" || microsite.TipeSection == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Content, Description, and TipeSection are required fields",
		})
	}

	// Check if a microsite with the same TipeSection already exists
	var existingMicrosite entity.Microsite
	if err := database.DB.Where("tipe_section = ?", microsite.TipeSection).First(&existingMicrosite).Error; err == nil {
		// If an entry exists, return an error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Max 1 description per TipeSection",
		})
	}

	// Save the microsite to the database
	if err := database.DB.Create(&microsite).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create microsite",
			"error":   err.Error(),
		})
	}

	// If TipeSection is 'signature', insert points into Signature table
	if microsite.TipeSection == "signature" {
		for _, point := range req.Points {
			signature := entity.Signature{
				Title:       point["title"],
				Description: point["description"],
			}
			if err := database.DB.Create(&signature).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to create signature",
					"error":   err.Error(),
				})
			}
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Microsite created successfully",
		"microsite": microsite,
	})
}

// GetMicrosite handles retrieving a microsite by ID
func GetMicrosite(c *fiber.Ctx) error {
	id := c.Query("id")
	category := c.Query("category")

	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid ID",
				"error":   err.Error(),
			})
		}

		var microsite entity.Microsite
		if err := database.DB.First(&microsite, idInt).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Microsite not found",
			})
		}

		// Retrieve all signatures directly
		var signatures []entity.Signature
		if err := database.DB.Find(&signatures).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve signatures",
				"error":   err.Error(),
			})
		}

		// Format data points from signature
		points := make([]map[string]string, len(signatures))
		for i, signature := range signatures {
			points[i] = map[string]string{
				"title":       signature.Title,
				"description": signature.Description,
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Microsite and signature retrieved successfully",
			"microsite": map[string]interface{}{
				"id":           microsite.ID, // Include ID
				"content":      microsite.Content,
				"description":  microsite.Description,
				"tipe_section": microsite.TipeSection,
				"image":        microsite.Image,
				"points":       points,
			},
		})
	}

	var microsites []entity.Microsite
	query := database.DB

	if category != "" {
		query = query.Where("tipe_section = ?", category)
	}

	if err := query.Find(&microsites).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve microsites",
			"error":   err.Error(),
		})
	}

	// Fetch all signatures directly
	var signatures []entity.Signature
	if err := database.DB.Find(&signatures).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve signatures",
			"error":   err.Error(),
		})
	}

	// Prepare the result
	var result []map[string]interface{}
	for _, microsite := range microsites {
		item := map[string]interface{}{
			"id":           microsite.ID, // Include ID
			"content":      microsite.Content,
			"description":  microsite.Description,
			"tipe_section": microsite.TipeSection,
			"image":        microsite.Image,
			"points":       []map[string]string{}, // Default to empty array
		}

		if microsite.TipeSection == "signature" {
			points := make([]map[string]string, len(signatures))
			for i, signature := range signatures {
				points[i] = map[string]string{
					"title":       signature.Title,
					"description": signature.Description,
				}
			}
			item["points"] = points
		}

		result = append(result, item)
	}

	// Return an empty array if no microsites were found
	if len(result) == 0 {
		result = []map[string]interface{}{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Microsites retrieved successfully",
		"microsites": result,
	})
}

// UpdateMicrosite handles updating an existing microsite by ID
func UpdateMicrosite(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"error":   err.Error(),
		})
	}

	// Parse request body
	var requestBody struct {
		Content     string              `json:"content"`
		Description string              `json:"description"`
		Image       string              `json:"image"`
		TipeSection string              `json:"tipe_section"`
		Points      []map[string]string `json:"points"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Fetch the microsite to check its TipeSection
	var microsite entity.Microsite
	if err := database.DB.First(&microsite, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Microsite not found",
		})
	}

	// Update the microsite record
	if err := database.DB.Model(&microsite).Updates(entity.Microsite{
		Content:     requestBody.Content,
		Description: requestBody.Description,
		Image:       requestBody.Image,
		TipeSection: requestBody.TipeSection,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update microsite",
			"error":   err.Error(),
		})
	}

	// If TipeSection is 'signature', handle signatures
	if microsite.TipeSection == "signature" {
		// Delete all existing signatures for the microsite
		if err := database.DB.Where("1 = 1").Delete(&entity.Signature{}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to delete existing signatures",
				"error":   err.Error(),
			})
		}

		// Create new signatures from the payload
		for _, point := range requestBody.Points {
			newSignature := entity.Signature{
				Title:       point["title"],
				Description: point["description"],
			}
			if err := database.DB.Create(&newSignature).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to create signature",
					"error":   err.Error(),
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Microsite and related signatures updated successfully",
	})
}

// DeleteMicrosite handles deleting a microsite by ID
func DeleteMicrosite(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID is required",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"error":   err.Error(),
		})
	}

	// Perform deletion of all signatures directly
	if err := database.DB.Unscoped().Where("1 = 1").Delete(&entity.Signature{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete signatures",
			"error":   err.Error(),
		})
	}

	// Perform deletion of the specified microsite
	if err := database.DB.Delete(&entity.Microsite{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete microsite",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Microsite and all signatures deleted successfully",
	})
}
