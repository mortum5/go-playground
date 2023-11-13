package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	PublicRoutes(app)

	app.Listen(":5000")
}
