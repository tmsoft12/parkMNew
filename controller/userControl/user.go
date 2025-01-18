package usercontrol

import (
	"park/database"
	modelsuser "park/models/modelsUser"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Swagger için genel API bilgileri
// @title           User Management API
// @version         1.0
// @description     API for user registration, login, and logout.
// @contact.name    API Support
// @contact.email   support@park.com
// @host            localhost:3000
// @BasePath        /api

var store = session.New()

// Register a new user
// @Summary      Register User
// @Description  Creates a new user and stores their hashed password.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body modelsuser.User true "User Registration Data"
// @Success      201 {object} map[string]string "message: User Created"
// @Failure      400 {object} map[string]string "message: Bad Request"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /register [post]
func Register(c *fiber.Ctx) error {
	var user modelsuser.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Bad Request", "error": err.Error()})
	}

	// Password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error hashing password"})
	}
	user.Password = string(hashedPassword)

	// Save user in the database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User Created"})
}

// Login a user
// @Summary      Login User
// @Description  Authenticates a user and starts a session.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        credentials body LoginInput true "User Login Data"
// @Success      200 {object} map[string]string "message: Login successful"
// @Failure      400 {object} map[string]string "message: Bad Request"
// @Failure      401 {object} map[string]string "message: Invalid credentials"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /login [post]
func Login(c *fiber.Ctx) error {

	var loginInput LoginInput
	if err := c.BodyParser(&loginInput); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Bad Request", "error": err.Error()})
	}

	if loginInput.Username == "" || loginInput.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Username and password are required"})
	}

	var user modelsuser.User
	if err := database.DB.Where("username = ?", loginInput.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"message": "Login successful"})
}

// Logout a user
// @Summary      Logout User
// @Description  Ends the session of a logged-in user.
// @Tags         User
// @Produce      json
// @Success      200 {object} map[string]string "message: Logout successful"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /logout [post]
func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error destroying session"})
	}

	return c.JSON(fiber.Map{"message": "Logout successful"})
}
