package handlers

import (
	"fmt"

	"github.com/arod1213/auto_ingestion/middleware"
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
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	title := c.Query("title")
	state := c.Query("state")

	query := db.Table("shares").
		Select("shares.*").
		Joins("INNER JOIN songs ON songs.id = shares.song_id").
		Where("shares.user_id = ?", userID)

	if title != "" {
		str := "%" + title + "%"
		query = query.Where("songs.title LIKE ? OR songs.artist LIKE ?", str, str)
	}

	if state == "" {
		query = query.Where("songs.registered = FALSE")
	} else {
		query = query.Where("songs.registered = TRUE")
	}

	var shares []models.Share
	err = query.Preload("Song").Find(&shares).Error

	if err != nil {
		c.JSON(400, gin.H{"error": "could not find any songs"})
		return
	}
	c.JSON(200, gin.H{"data": shares})
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
