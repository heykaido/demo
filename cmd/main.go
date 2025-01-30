package main

import (
	"ev-api/config"
	"ev-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()

	app := fiber.New()

	app.Post("/login", handlers.Login)
	app.Get("/data-report", handlers.DataReport)

	app.Listen(":8080")
}
