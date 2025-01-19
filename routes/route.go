package routes

import (
	"os"
	carcontrol "park/controller/carControl"
	usercontrol "park/controller/userControl"
	"park/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Init(app *fiber.App) {
	plate := os.Getenv("IMAGE_URL")
	app.Static("/plate", plate)

	user := app.Group("/api/v1/auth")
	user.Post("/register", usercontrol.Register)
	user.Post("/login", usercontrol.Login)
	user.Post("/logout", usercontrol.Logout)
	user.Get("/users", usercontrol.ListUsers)
	app.Use(middleware.ExtractParkNoMiddleware)

	user.Get("me", usercontrol.Me)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	cars := app.Group("/api/v1", middleware.ExtractParkNoMiddleware)

	cars.Post("/createcar", carcontrol.CreateCar)
	cars.Get("/getallcars", carcontrol.GetCars)
	cars.Get("/getcar/:id", carcontrol.GetCar)
	cars.Get("/searchcar", carcontrol.SearchCar)
	cars.Put("/updatecar/:plate", carcontrol.UpdateCar)
	cars.Get("/ws/notification", websocket.New(carcontrol.Ws))
}
