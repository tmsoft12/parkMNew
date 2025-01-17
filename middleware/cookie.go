package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("user")

	if cookie == "" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	return c.Next()
}

func SetCookie(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "user",
		Value:    "john",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})
	return c.SendString("Cookie is set")
}
