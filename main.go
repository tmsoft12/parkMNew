package main

import (
	carcontrol "park/controller/carControl"
	"park/database"
	_ "park/docs"
	"park/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title Airline API
// @host 192.168.100.192:3000
// @BasePath /api/v1
func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)

	go carcontrol.HandleMessages()

	routes.Init(app)

	app.Listen(":3000")
}
