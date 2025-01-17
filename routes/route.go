package routes

import (
	carcontrol "park/controller/carControl"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Init(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	cars := app.Group("/api/v1")
	cars.Post("/createcar", carcontrol.CreateCar)
	cars.Get("/getallcars", carcontrol.GetCars)
	cars.Get("/getcar/:id", carcontrol.GetCar)
	cars.Get("/searchcar", carcontrol.SearchCar)
	cars.Put("/updatecar/:plate", carcontrol.UpdateCar)
	cars.Get("/ws/notification", websocket.New(carcontrol.Ws))
}
