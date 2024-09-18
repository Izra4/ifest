package main

import (
	"IFEST/internals/config"
	"IFEST/internals/server"
)

func main() {
	//err := config.LoadEnv()
	//if err != nil {
	//	return
	//}
	_, err := config.ConnectDB()
	if err != nil {
		return
	}

	server.Initialize()
}
