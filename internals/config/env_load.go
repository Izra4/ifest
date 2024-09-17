package config

import "github.com/joho/godotenv"

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	return nil
}
