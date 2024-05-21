package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Setenv("IP_LIMIT", os.Getenv("IP_LIMIT"))
	os.Setenv("TOKEN_LIMIT", os.Getenv("TOKEN_LIMIT"))
}
