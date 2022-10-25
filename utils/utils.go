package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error oading .env file")
	}

	return os.Getenv(key)
}