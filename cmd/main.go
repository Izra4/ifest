package main

import (
	"IFEST/internals/config"
	"IFEST/internals/handlers"
	"IFEST/internals/repositories"
	"IFEST/internals/server"
	"IFEST/internals/services"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		return
	}
	db, err := config.ConnectDB()
	if err != nil {
		return
	}

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	docsRepository := repositories.NewDocsRepository(db)
	docsService := services.NewDocsService(docsRepository)
	docsHandler := handlers.NewDocHandler(docsService)

	httpServer := server.NewServer(userHandler, docsHandler)

	httpServer.Initialize()
}
