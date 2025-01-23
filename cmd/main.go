package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/maksimUlitin/internal/lib"
	"github.com/maksimUlitin/internal/routes"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logger.Warn("Error loading .env file", "error", err)
	} else {
		logger.Info(".env file loaded successfully")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("No PORT environment variable found, using default", "port", port)
	} else {
		logger.Info("Server port set from environment variable", "port", port)
	}

	router := gin.Default()

	routes.AuthRouter(router)
	logger.Info("Auth routes initialized")

	routes.UserRouter(router)
	logger.Info("User routes initialized")

	routes.APIRouter(router)
	logger.Info("API routes initialized")

	logger.Info("Starting server", "port", port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
