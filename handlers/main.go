package handlers

import (
	"log"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchAndSaveTracks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		id := c.Param("id")
		method := c.Query("method")
		songs := services.GetSpotifyTracks(method, id)

		if len(songs) == 0 {
			c.JSON(400, gin.H{"error": "no songs found: ensure your playlist is public"})
			return
		}

		shares, err := services.SaveSongs(db, userID, songs)
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to save songs"})
			return
		}

		c.JSON(200, gin.H{"data": shares})

	}
}

func SaveTracks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		var songs []models.Song
		err = c.ShouldBindJSON(&songs)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}
		shares, err := services.SaveSongs(db, userID, songs)
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to save songs"})
			return
		}

		c.JSON(200, gin.H{"data": shares})
	}
}

func FetchTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("FetchTracks - All keys:", c.Keys)

		id := c.Param("id")
		method := c.Query("method")
		songs := services.GetSpotifyTracks(method, id)

		if len(songs) == 0 {
			c.JSON(400, gin.H{"error": "no songs found: ensure your playlist is public"})
			return
		}

		c.JSON(200, gin.H{"data": songs})
	}
}

func UpdateSongs(db *gorm.DB, songs []models.Song) error {
	return db.Save(&songs).Error
}
