package routes

import (
	"os"
	authconrol "park/controller/authConrol"
	carcontrol "park/controller/carControl"
	usercontroller "park/controller/userController"
	"park/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Init(app *fiber.App) {
	plate := os.Getenv("IMAGE_URL")
	app.Static("/plate", plate)

	auth := app.Group("/api/v1/auth")
	auth.Post("/register", authconrol.Register)
	auth.Post("/login", authconrol.Login)
	auth.Post("/logout", authconrol.Logout)
	app.Use(middleware.ExtractParkNoMiddleware)

	auth.Get("/me", authconrol.Me)

	cars := app.Group("/api/v1", middleware.ExtractParkNoMiddleware)

	cars.Post("/createcar", carcontrol.CreateCar)
	cars.Get("/getallcars", carcontrol.GetCars)
	cars.Get("/getcar/:id", carcontrol.GetCar)
	cars.Get("/searchcar", carcontrol.SearchCar)
	cars.Put("/updatecar/:plate", carcontrol.UpdateCar)
	cars.Get("/ws/notification", websocket.New(carcontrol.Ws))

	admin := app.Group("/api/v1/admin")
	admin.Post("/user", usercontroller.CreateUser)
	admin.Get("/user/:id", usercontroller.GetUserByID)

}
