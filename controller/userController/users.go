package usercontroller

import (
	"park/database"
	modelsuser "park/models/modelsUser"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      Create User
// @Description  Creates a new user and stores their hashed password.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        user body modelsuser.User true "User Registration Data"
// @Success      201 {object} map[string]string "message: User Created"
// @Failure      400 {object} map[string]string "message: Bad Request"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /admin/user [post]
func CreateUser(c *fiber.Ctx) error {
	var user modelsuser.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Bad Request", "error": err.Error()})
	}

	var existingUser modelsuser.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Username already exists"})
	}

	if len(user.Password) < 8 {
		return c.Status(400).JSON(fiber.Map{"message": "Password must be at least 8 characters long"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error hashing password"})
	}
	user.Password = string(hashedPassword)

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error", "error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User Created"})
}

// @Summary      Get User by ID
// @Description  Retrieves user details by their ID.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} modelsuser.User "User details"
// @Failure      404 {object} map[string]string "message: User not found"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /admin/user/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID format",
		})
	}

	var user modelsuser.User
	if err := database.DB.Where("id = ?", uint(id)).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(404).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(user)
}
