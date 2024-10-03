package server

import (
	"IFEST/internals/core/middleware"
	"IFEST/internals/handlers"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	userHandler     handlers.UserHandler
	docsHandler     handlers.DocHandler
	userDocsHandler handlers.UserDocHandler
	cronJob         handlers.CronJob
	reportHandler   handlers.ReportHandler
}

func NewServer(
	userHandler handlers.UserHandler,
	docsHanlder handlers.DocHandler,
	userDocsHandler handlers.UserDocHandler,
	cronJob handlers.CronJob,
	reportHandler handlers.ReportHandler,
) *Server {
	return &Server{
		userHandler:     userHandler,
		docsHandler:     docsHanlder,
		userDocsHandler: userDocsHandler,
		cronJob:         cronJob,
		reportHandler:   reportHandler,
	}
}

func (s *Server) Initialize() {
	app := fiber.New()

	middleware.Cors(app)
	s.cronJob.Start()
	defer s.cronJob.Stop()

	user := app.Group("/api/user")
	docs := app.Group("/api/document")
	access := docs.Group("/access")
	reports := app.Group("/api/reports")

	user.Get("/login/google", s.userHandler.GoogleLogin)
	user.Get("/profile", middleware.Authentication(), s.userHandler.Profile)
	user.Get("/auth/google/callback", s.userHandler.GoogleCallback)
	user.Post("/register", s.userHandler.Create)
	user.Post("/login", s.userHandler.Login)

	docs.Post("/upload", middleware.Authentication(), s.docsHandler.Upload)
	docs.Get("/download", s.userDocsHandler.Download)
	docs.Get("/detail/:id", middleware.Authentication(), s.docsHandler.GetByID)
	docs.Get("/all", middleware.Authentication(), s.docsHandler.GetAll)
	docs.Get("/history", middleware.Authentication(), s.userDocsHandler.GetHistoryByUserID)

	access.Post("/:id", middleware.Authentication(), s.userDocsHandler.Create)
	access.Delete("/delete", middleware.Authentication(), s.userDocsHandler.DeleteAccess)

	reports.Post("/create", middleware.Authentication(), s.reportHandler.CreateReport)
	reports.Get("/", s.reportHandler.GetReportsByUserID)
	reports.Get("/:id", s.reportHandler.GetReportByID)
	reports.Put("/:id", s.reportHandler.UpdateReport)
	reports.Delete("/:id", s.reportHandler.DeleteReport)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!\nTesting the jenkins here")
	})
	app.Get("/email", s.userDocsHandler.TestEmail)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalln("Server error: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := app.Shutdown(); err != nil {
		log.Fatalln("Server Shutdown: ", err)
	}
}
