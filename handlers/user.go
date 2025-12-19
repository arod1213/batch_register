package handlers

import (
	"fmt"
	"strconv"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMe(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := middleware.GetUser(c, db)
		if err != nil {
			fmt.Println("err is ", err.Error())
			c.JSON(401, gin.H{"error": "user not found"})
			return
		}

		c.JSON(200, gin.H{"data": user})
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			fmt.Println("conv err is ", err.Error())
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		var user models.User
		err = c.ShouldBindJSON(&user)
		if err != nil {
			fmt.Println("bind err is ", err.Error())
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		if user.ID != uint(userID) {
			fmt.Println("id err is ", err.Error())
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		err = db.Save(&user).Error
		if err != nil {
			fmt.Println("save err is ", err.Error())
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.JSON(200, gin.H{"data": user})
	}
}
