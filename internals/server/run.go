package server

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func Initialize() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!\nTesting the jenkins here")
	})
	err := app.Listen(":3000")
	if err != nil {
		log.Fatalln("Server error: ", err)
	}
}
