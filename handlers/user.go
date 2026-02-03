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
		authed, err := middleware.GetUser(c, db)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		id := c.Param("id")
		userID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			fmt.Println("conv err is ", err.Error())
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		if authed.ID != uint(userID) {
			c.JSON(403, gin.H{"error": "forbidden"})
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

		// Only allow updating the authenticated user.
		// GORM Save() will overwrite every field, so we ensure we're saving against the authed row.
		user.ID = authed.ID
		user.Username = authed.Username
		// Preserve password unless the client sends the existing (hashed) value.
		// (The frontend uses GET /user then PUT /user/update/:id with the full record.)
		if user.Password == "" {
			user.Password = authed.Password
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
