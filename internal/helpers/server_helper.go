package helpers

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

func TryRunServer(router *gin.Engine, port string) error {
	if IsPortAvailable(port) {
		return router.Run(":" + port)
	}
	return fmt.Errorf("port %s is not available", port)
}

func IsPortAvailable(port string) bool {
	addr := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	_ = listener.Close()
	return true
}
