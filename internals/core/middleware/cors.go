package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://aman-in.vercel.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Origin, Accept, X-Requested-With, Authorization, X-CSRF-Token",
		MaxAge:           86400,
		AllowCredentials: true,
	}))
}
