package handlers

import (
	"fmt"
	"strconv"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IdentifyUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := middleware.GetUser(c, db)
		if err != nil {
			c.JSON(404, gin.H{"error": "unauthorized"})
			return
		}

		if user.GeniusID != nil {
			c.JSON(200, gin.H{"data": "already identified"})
			return
		}

		var songs []models.Song
		err = db.
			Joins("LEFT JOIN shares on shares.song_id = songs.id").
			Where("shares.user_id = ?", user.ID).
			Find(&songs).
			Error

		if err != nil {
			c.JSON(500, gin.H{"error": "could not find any songs"})
			return
		}

		_, err = services.IdentifyUser(db, *user, songs)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to identify"})
			return
		}
		c.JSON(200, gin.H{"data": "success"})
	}
}

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
