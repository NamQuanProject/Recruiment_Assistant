package initializers

import (
	"log"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load("./internal/apigateway/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
} 