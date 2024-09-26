package main

import (
	"IFEST/internals/config"
	"IFEST/internals/handlers"
	"IFEST/internals/repositories"
	"IFEST/internals/server"
	"IFEST/internals/services"
	"log"
	"os"
)

func main() {
	//err := config.LoadEnv()
	//if err != nil {
	//	return
	//}
	db, err := config.ConnectDB()
	if err != nil {
		return
	}
	log.Println("=============================\n", os.Getenv("CLIENT_ID"), "\n=============================")
	log.Println("=============================\n", os.Getenv("CLIENT_SECRET"), "\n=============================")
	log.Println("=============================\n", os.Getenv("REDIRECT_URL"), "\n=============================")
	userRepository := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepository)

	userHandler := handlers.NewUserHandler(userService)

	httpServer := server.NewServer(userHandler)

	httpServer.Initialize()
}
