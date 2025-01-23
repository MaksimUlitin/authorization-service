package middleware

import (
	"fmt"
	"github.com/maksimUlitin/internal/helpers"
	"github.com/maksimUlitin/internal/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			logger.Error("Authentication failed: no Authorization header provided",
				"ip", c.ClientIP(),
				"path", c.Request.URL.Path)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("no Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			logger.Error("Authentication failed: token validation error",
				"error", err,
				"ip", c.ClientIP(),
				"path", c.Request.URL.Path)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		logger.Info("Authentication successful",
			"email", claims.Email,
			"user_id", claims.Id,
			"user_type", claims.UserType,
			"ip", c.ClientIP(),
			"path", c.Request.URL.Path)

		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.Id)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
