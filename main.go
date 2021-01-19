package main

import "github.com/gofiber/fiber"

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Static("/", "./public/home")
	app.Static("/images", "./public/images")

	app.Listen(":3000")
}
