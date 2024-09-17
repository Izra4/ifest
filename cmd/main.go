package main

import (
	"IFEST/internals/config"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		return
	}
	_, err = config.ConnectDB()
	if err != nil {
		return
	}
}
