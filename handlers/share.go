package handlers

import (
	"strconv"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SaveShare(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	var share models.Share
	err = c.ShouldBindJSON(&share)
	if err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	if share.ID != uint(id) {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	tx := db.Begin()
	err = tx.Save(&share).Error
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed to save"})
		return
	}

	if share.Song.Iswc != nil {
		err = tx.Save(&share.Song).Error
		if err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "failed to save song"})
			return
		}
	}

	tx.Commit()
	c.JSON(200, gin.H{"data": "success"})
}
