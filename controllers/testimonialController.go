package controllers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
)

var deliciousKeywords = []string{
    "enak",
    "lezat",
    "luar biasa",
    "memuaskan",
    "nikmat",
    "sangat",
    "terbaik",
    "cita rasa",
    "wajib coba",
}

var negativeKeywords = []string{
    "gaenak",
    "buruk",
    "tidak enak",
}

// UpdateTestimonialActivity updates the active status of testimonials based on keywords
func UpdateTestimonialActivity(c *fiber.Ctx) error {
    var testimonies []entity.Testimoni

    // Retrieve all testimonials from the database
    if err := database.DB.Find(&testimonies).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to retrieve testimonials",
            "error":   err.Error(),
        })
    }

    // Keep track of how many convincing testimonials have been marked active
    activeCount := 0

    for i := range testimonies {
        // Initialize the active field as false
        active := false

        // Check for delicious keywords in the description
        for _, keyword := range deliciousKeywords {
            if strings.Contains(strings.ToLower(testimonies[i].Description), strings.ToLower(keyword)) {
                active = true
                break // Stop checking after the first match
            }
        }

        // Check for negative keywords in the description
        for _, keyword := range negativeKeywords {
            if strings.Contains(strings.ToLower(testimonies[i].Description), strings.ToLower(keyword)) {
                active = false // If a negative keyword is found, set active to false
                break // Stop checking after the first match
            }
        }

        // Update the testimony's active status based on the checks
        if active {
            // Only mark as active if it's not already active and we haven't reached the limit
            if !testimonies[i].Active && activeCount < 2 {
                testimonies[i].Active = true
                activeCount++
            }
        } else {
            // If active was set to false due to negative keywords or no delicious keywords, mark it inactive
            testimonies[i].Active = false
        }

        // Save the updated testimony back to the database
        if err := database.DB.Save(&testimonies[i]).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Failed to update testimony",
                "error":   err.Error(),
            })
        }
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message":      "Successfully updated testimonials activity status",
        "testimonials": testimonies, // Return all testimonials
    })
}

func CreateTestimoni(c *fiber.Ctx) error {
	testimoni := new(entity.Testimoni)

	// Parse the request body into the Testimoni struct
	if err := c.BodyParser(testimoni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	// Set the timestamps for created and updated times
	testimoni.CreatedAt = time.Now()
	testimoni.UpdatedAt = time.Now()

	// Save the testimoni to the database
	if err := database.DB.Create(&testimoni).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create testimoni",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Testimoni created successfully",
		"testimoni": testimoni,
	})
}

// GetAllTestimoni retrieves all testimonials from the database
func GetAllTestimoni(c *fiber.Ctx) error {
	var testimonies []entity.Testimoni

	// Retrieve all testimonials from the database
	if err := database.DB.Find(&testimonies).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve testimonials",
			"error":   err.Error(),
		})
	}

	// If no testimonials found, return not found response
	if len(testimonies) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No testimonials found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Successfully retrieved all testimonials",
		"testimonials": testimonies,
	})
}