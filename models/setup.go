package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() {
	connectionString := os.Getenv("MONGODB_URI")
	if connectionString == "" {
		// Fallback for local development
		connectionString = "mongodb+srv://auth-api:oPrJhEuCvRF*Y9@cluster0.64wve.mongodb.net/"
	}
	clientOptions := options.Client().ApplyURI(connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB: ", err)
	}

	fmt.Println("Connected to MongoDB!")
	DB = client.Database("gin_jwt_db")
}

func DBMigrate() {
	// MongoDB is schema-less, but we can ensure indexes here if needed.
	fmt.Println("Database migration (schema-less) completed.")
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
