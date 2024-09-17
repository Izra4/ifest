package server

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func Initialize() {
	app := fiber.New()

	err := app.Listen(":3000")
	if err != nil {
		log.Fatalln("Server error: ", err)
	}
}
