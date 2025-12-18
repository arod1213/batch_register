package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		err := errors.New("missing user id")
		return 0, err
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		err := errors.New("invalid User ID type")
		return 0, err
	}

	return userID, nil
}
