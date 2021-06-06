package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "serviceAccountKey.json")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectFirestore() *firestore.Client {
	PROJECT_ID := os.Getenv("PROJECT_ID")
	client, err := firestore.NewClient(context.Background(), PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create a Firestore Client: %v", err)
	}

	return client
}
