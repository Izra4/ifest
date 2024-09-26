package server

import (
	"IFEST/internals/core/middleware"
	"IFEST/internals/handlers"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Server struct {
	userHandler handlers.UserHandler
}

func NewServer(userHandler handlers.UserHandler) *Server {
	return &Server{
		userHandler: userHandler,
	}
}

func (s *Server) Initialize() {
	app := fiber.New()

	middleware.Cors(app)
	user := app.Group("/api/user")

	user.Get("/login/google", s.userHandler.GoogleLogin)
	user.Get("/profile", middleware.Authentication(), s.userHandler.Profile)
	user.Get("/auth/google/callback", s.userHandler.GoogleCallback)
	user.Post("/register", s.userHandler.Create)
	user.Post("/login", s.userHandler.Login)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!\nTesting the jenkins here")
	})
	err := app.Listen(":3000")
	if err != nil {
		log.Fatalln("Server error: ", err)
	}
}
