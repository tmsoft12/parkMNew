package usercontrol

import (
	"park/database"
	modelsuser "park/models/modelsUser"
	"park/util"

	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	IsActive  bool      `json:"is_active"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	user.IsActive = false
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

// @Summary      Login User
// @Description  Authenticates a user and returns a JWT token in a cookie.
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

	var user modelsuser.User
	if err := database.DB.Where("username = ?", loginInput.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	token, err := util.CreateJWT(user.Id, user.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error creating JWT", "error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   0,
	})

	return c.JSON(fiber.Map{"message": "Login successful"})
}

// @Summary      Logout User
// @Description  Ends the session of a logged-in user by deleting the JWT token cookie.
// @Tags         User
// @Produce      json
// @Success      200 {object} map[string]string "message: Logout successful"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /logout [post]
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		HTTPOnly: true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   -1,
	})

	return c.JSON(fiber.Map{"message": "Logout successful"})
}

// @Summary      List Users
// @Description  Retrieves a list of all users.
// @Tags         User
// @Produce      json
// @Success      200 {array} UserResponse
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /users [get]
func ListUsers(c *fiber.Ctx) error {
	var users []modelsuser.User
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Error retrieving users", "error": err.Error()})
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.Id,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			IsActive:  user.IsActive,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return c.JSON(userResponses)
}
