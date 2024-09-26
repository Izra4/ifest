package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000",
		AllowMethods:     "*",
		AllowHeaders:     "Content-Type, Origin, Accept, X-Requested-With, Authorization, X-CSRF-Token, Access-Control-Allow-Origin",
		MaxAge:           86400,
		AllowCredentials: true,
	}))
}
