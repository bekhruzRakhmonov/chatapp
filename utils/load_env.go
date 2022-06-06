package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVariable(key string) string {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Println(err)
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
