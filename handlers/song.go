package handlers

import (
	"fmt"
	"log"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MarkRegistered(c *gin.Context, db *gorm.DB) {
	isrc := c.Param("isrc")
	err := db.Model(&models.Song{}).Where("isrc = ?", isrc).Update("registered", true).Error
	if err != nil {
		fmt.Println("err is ", err.Error())
		c.JSON(400, gin.H{"error": "failed to mark"})
		return
	}
	c.JSON(200, gin.H{"data": "success"})
}

func FetchSongs(c *gin.Context, db *gorm.DB) {
	title := c.Query("title")
	state := c.Query("state")

	query := db
	if title != "" {
		str := "%" + title + "%"
		log.Println("searching for title: ", title)
		query = query.Where("title LIKE ? OR artist LIKE ?", str, str)
	}

	if state == "" {
		query = query.Where("registered = FALSE")
	} else {
		query = query.Where("registered = TRUE")
	}

	var songs []models.Song
	err := query.Preload("Share").Find(&songs).Order("release_date DESC").Error

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
