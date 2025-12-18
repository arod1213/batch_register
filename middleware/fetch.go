package middleware

import (
	"errors"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func GetUser(c *gin.Context, db *gorm.DB) (*models.User, error) {
	userID, err := GetUserID(c)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = db.Where("id = ?", userID).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
