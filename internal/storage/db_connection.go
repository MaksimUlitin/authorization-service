package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/maksimUlitin/config"
	"github.com/maksimUlitin/internal/lib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
	config.LoadConfigEnv()
	logger.Info(".env file loaded successfully")

	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")

	if mongoUser == "" || mongoPassword == "" || mongoHost == "" || mongoPort == "" {
		logger.Error("One or more MongoDB environment variables are not set")
		os.Exit(1)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPassword, mongoHost, mongoPort)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
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
