package main

import (
	"github.com/maksimUlitin/config"
	"github.com/maksimUlitin/internal/helpers"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maksimUlitin/internal/lib"
	"github.com/maksimUlitin/internal/routes"
)

func main() {
	config.LoadConfigEnv()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8088"
	}

	serverPortFallback := os.Getenv("SERVER_PORT_FALLBACK")
	if serverPortFallback == "" {
		serverPortFallback = "8089"
	}

	router := gin.Default()

	routes.AuthRouter(router)
	logger.Info("Auth routes initialized")

	routes.UserRouter(router)
	logger.Info("User routes initialized")

	routes.APIRouter(router)
	logger.Info("API routes initialized")

	logger.Info("Starting server", "port", serverPort)
	if err := helpers.TryRunServer(router, serverPort); err != nil {
		logger.Warn("Main port is occupied, trying fallback port", slog.String("fallbackPort", serverPortFallback))
		if err := helpers.TryRunServer(router, serverPortFallback); err != nil {
			logger.Error("Server failed to start on both main and fallback ports", slog.Any("error", err))
			log.Fatal(err)
		}
	}
}
