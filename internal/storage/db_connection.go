package storage

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/maksimUlitin/internal/lib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}
	logger.Info(".env file loaded successfully")

	mongoDB := os.Getenv("MONGODB_URL")
	if mongoDB == "" {
		logger.Error("MONGODB_URL environment variable is not set")
		os.Exit(1)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDB))
	if err != nil {
		logger.Error("Failed to create MongoDB client", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Error("Failed to ping MongoDB", "error", err)
		os.Exit(1)
	}

	logger.Info("Connected to MongoDB successfully")

	return client
}

func OpenCollection(client mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("cluster0").Collection(collectionName)
	logger.Info("Opened MongoDB collection", "collection", collectionName)
	return collection
}
