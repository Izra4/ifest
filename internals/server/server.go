package server

import (
	"IFEST/internals/core/middleware"
	"IFEST/internals/handlers"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Server struct {
	userHandler handlers.UserHandler
	docsHandler handlers.DocHandler
}

func NewServer(userHandler handlers.UserHandler, docsHanlder handlers.DocHandler) *Server {
	return &Server{
		userHandler: userHandler,
		docsHandler: docsHanlder,
	}
}

func (s *Server) Initialize() {
	app := fiber.New()

	middleware.Cors(app)
	user := app.Group("/api/user")
	docs := app.Group("/api/document")

	user.Get("/login/google", s.userHandler.GoogleLogin)
	user.Get("/profile", middleware.Authentication(), s.userHandler.Profile)
	user.Get("/auth/google/callback", s.userHandler.GoogleCallback)
	user.Post("/register", s.userHandler.Create)
	user.Post("/login", s.userHandler.Login)

	docs.Post("/upload", middleware.Authentication(), s.docsHandler.Upload)
	docs.Get("/download/:id", s.docsHandler.Download)
	docs.Get("/all", middleware.Authentication(), s.docsHandler.GetAll)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!\nTesting the jenkins here")
	})
	err := app.Listen(":3000")
	if err != nil {
		log.Fatalln("Server error: ", err)
	}
}
