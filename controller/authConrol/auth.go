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
// @Router       /auth/register [post]
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
//
//	@Param        credentials body LoginInput true "User Login Data" {
//	  "username": "exampleUser",
//	  "password": "examplePassword",
//	  "parkno": "P4"
//	}
//
// @Success      200 {object} map[string]string "message: Login successful"
// @Failure      400 {object} map[string]string "message: Bad Request"
// @Failure      401 {object} map[string]string "message: Invalid credentials"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router       /auth/login [post]
func Login(c *fiber.Ctx) error {
	var loginInput struct {
		LoginInput
		ParkNo string `json:"parkno" validate:"required"`
	}
	if err := c.BodyParser(&loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	var user modelsuser.User
	if err := database.DB.Where("username = ?", loginInput.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}
	if !user.IsActive {
		return c.Status(401).JSON("is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	token, err := util.CreateJWT(user.Id, user.Username, user.Role, loginInput.ParkNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating JWT",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   86400,
	})

	return c.JSON(fiber.Map{"message": "Login successful"})
}

// @Summary      Logout User
// @Description  Ends the session of a logged-in user by deleting the JWT token cookie.
// @Tags         User
// @Produce      json
// @Success      200 {object} map[string]string "message: Logout successful"
// @Failure      500 {object} map[string]string "message: Internal Server Error"
// @Router      /auth/logout [post]
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
// @Router       /auth/users [get]
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

// @Summary      Get current user information
// @Description  Retrieves the current user's username, role, and user ID from the JWT token.
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "Returns user information"
// @Failure      400 {object} map[string]string "message: Bad Request - Missing data from middleware"
// @Failure      401 {object} map[string]string "message: Unauthorized - Invalid token"
// @Failure      500 {object} map[string]string "message: Internal Server Error - Missing data from middleware"
// @Router       /auth/me [get]
func Me(c *fiber.Ctx) error {
	usernameVal := c.Locals("username")
	roleVal := c.Locals("role")
	userIDVal := c.Locals("user_id")
	parkno := c.Locals("parkno")

	if usernameVal == nil || roleVal == nil || userIDVal == nil || parkno == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Missing data from middleware",
		})
	}

	username, ok := usernameVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid username type",
		})
	}
	park, ok := parkno.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid parkno type",
		})
	}
	role, ok := roleVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid role type",
		})
	}

	userID, ok := userIDVal.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error - Invalid user ID type",
		})
	}

	return c.JSON(fiber.Map{
		"username": username,
		"role":     role,
		"user_id":  userID,
		"parkno":   park,
	})
}
