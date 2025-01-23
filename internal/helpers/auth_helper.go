package helpers

import (
	"errors"
	"github.com/maksimUlitin/internal/lib"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	if userType != role {
		logger.Warn("Unauthorized access attempt",
			"required_role", role,
			"user_role", userType,
			"ip", c.ClientIP(),
			"path", c.Request.URL.Path)
		return errors.New("Unauthorized to access this resource")
	}

	logger.Info("User type check passed",
		"role", role,
		"ip", c.ClientIP(),
		"path", c.Request.URL.Path)
	return nil
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")

	if userType == "USER" && uid != userId {
		logger.Warn("Unauthorized access attempt: user ID mismatch",
			"user_type", userType,
			"request_user_id", userId,
			"session_user_id", uid,
			"ip", c.ClientIP(),
			"path", c.Request.URL.Path)
		return errors.New("Unauthorized to access this resource")
	}

	err = CheckUserType(c, userType)
	if err == nil {
		logger.Info("User type and ID match verified",
			"user_type", userType,
			"user_id", userId,
			"ip", c.ClientIP(),
			"path", c.Request.URL.Path)
	}
	return err
}
