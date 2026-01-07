package handlers

import (
	"os"
	"strconv"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMissingSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")

		// artistId := "2877875"
		artistID, err := strconv.ParseUint(c.Query("artist_id"), 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		missingSongs, err := services.GetMissingSongs(db, uint(artistID), accessToken, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": missingSongs})
		return
	}
}
