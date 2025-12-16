package handlers

import (
	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchSongs(c *gin.Context, db *gorm.DB) {
	title := c.Query("title")

	query := db
	if title != "" {
		query = query.Where("title LIKE ?", title)
	}

	var songs []models.Song
	err := query.Find(&songs).Order("release_date DESC").Error

	if err != nil {
		c.JSON(400, gin.H{"error": "could not find any songs"})
		return
	}
	c.JSON(200, gin.H{"data": songs})
}

func DeleteSongs(c *gin.Context, db *gorm.DB) {
	var songs []models.Song
	err := db.Where("title IS NOT NULL").Delete(&songs).Error
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to delete songs"})
		return
	}
	c.JSON(200, gin.H{"data": "successfully reset songs"})
}
