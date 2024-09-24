package controllers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/marcolcu/go-resto-app/database"
	"github.com/marcolcu/go-resto-app/models/entity"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type LoginResponse struct {
	Message string      `json:"message"`
	Token   string      `json:"token"`
	User    entity.User `json:"user"`
}

func Login(c *fiber.Ctx) error {
	credentials := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := c.BodyParser(credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	// Fetch user from database
	var user entity.User
	if err := database.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Compare passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	// Generate JWT token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Send response with token and user data
	response := LoginResponse{
		Message: "Login successful",
		Token:   tokenString,
		User:    user,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func Register(c *fiber.Ctx) error {
	user := new(entity.User)

	// Parsing Body
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validation
	validate := validator.New()
	errValidate := validate.Struct(user)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to validate",
			"error":   errValidate.Error(),
		})
	}

	// Hashing Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error while hashing password",
		})
	}

	// Creating new user
	newUser := entity.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(hashedPassword), // Assign the hashed password
		Phone:     user.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Saving to database
	if err := database.DB.Debug().Create(&newUser).Error; err != nil {
		fmt.Println("Database error:", err.Error()) // Log for debugging
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
			"error":   err.Error(),
		})
	}

	// Generate JWT token
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Email: newUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate token",
		})
	}

	// Prepare response with message, token, and user data
	response := LoginResponse{
		Message: "User created successfully",
		Token:   tokenString,
		User:    newUser,
	}

	// Sending response with token and user data
	return c.Status(fiber.StatusCreated).JSON(response)
}
