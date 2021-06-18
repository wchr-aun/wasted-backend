package endpoints

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var USERS_TABLE string

func init() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	USERS_TABLE = os.Getenv("USERS_TABLE")
}
