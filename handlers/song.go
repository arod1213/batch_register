package handlers

import (
	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchSongs(c *gin.Context, db *gorm.DB) {
	var songs []models.Song
	err := db.Find(&songs).Error
	if err != nil {
		c.JSON(400, gin.H{"error": "could not find any songs"})
		return
	}
	c.JSON(200, gin.H{"data": songs})
}
