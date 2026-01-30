package handlers

import (
	"os"

	"github.com/arod1213/auto_ingestion/genius"
	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GeniusSearchArtistIDs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		artistName := c.Query("artistName")
		if artistName == "" {
			c.JSON(400, gin.H{"error": "artistName parameter is required"})
			return
		}
		keyword := c.Query("q")
		if keyword == "" {
			c.JSON(400, gin.H{"error": "q parameter is required"})
			return
		}
		artists, err := genius.GeniusSearchArtists(keyword, artistName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": artists})
	}
}

func GeniusSearch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("q")
		if keyword == "" {
			c.JSON(400, gin.H{"error": "q parameter is required"})
			return
		}
		hits, err := genius.GeniusSearch(keyword)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": hits})
	}
}

func GetMissingSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := middleware.GetUser(c, db)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")

		var artistID uint
		if user.GeniusID != nil {
			artistID = *user.GeniusID
		} else {
			c.JSON(400, gin.H{"error": "artist not found"})
			return
		}

		missingSongs, err := services.GetMissingSongs(db, artistID, accessToken, user.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": missingSongs})
	}
}
