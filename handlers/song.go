package handlers

import (
	"fmt"
	"strconv"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/royalties"
	"github.com/arod1213/auto_ingestion/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type payInfo struct {
	Total float64
	Items []royalties.Payment
}

type SongInfo struct {
	Song     models.Song
	Share    models.Share
	Payments payInfo
}

func GetSong(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		idStr := c.Param("songID")
		songID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		var song models.Song
		if err := db.
			Joins("JOIN shares ON shares.song_id = songs.id").
			Where("songs.id = ? AND shares.user_id = ?", uint(songID), userID).
			First(&song).Error; err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		var share models.Share
		if err := db.
			Where("song_id = ? AND user_id = ?", song.ID, userID).
			First(&share).Error; err != nil {
			c.JSON(404, gin.H{"error": "share not found"})
			return
		}

		var payments []royalties.Payment
		if err := db.
			Joins("LEFT JOIN shares on shares.id = payments.share_id").
			Where("shares.song_id = ?", song.ID).
			Find(&payments).Error; err != nil {
			c.JSON(400, gin.H{"error": "failed to load payments"})
			return
		}

		total := utils.Reduce(payments, 0.0, func(acc float64, p royalties.Payment) float64 {
			return acc + p.Earnings
		})

		info := SongInfo{
			Song:  song,
			Share: share,
			Payments: payInfo{
				Total: total,
				Items: payments,
			},
		}

		c.JSON(200, gin.H{"data": info})
	}
}

func MarkRegistered(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		isrc := c.Param("isrc")
		err := db.Model(&models.Song{}).Where("isrc = ?", isrc).Update("registered", true).Error
		if err != nil {
			fmt.Println("err is ", err.Error())
			c.JSON(400, gin.H{"error": "failed to mark"})
			return
		}
		c.JSON(200, gin.H{"data": "success"})
	}
}

func FetchSongs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func DeleteShares(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		var songIDs []uint
		err = c.ShouldBindJSON(&songIDs)
		if err != nil {
			c.JSON(400, gin.H{"error": "bad request"})
			return
		}

		err = db.
			Joins("LEFT JOIN songs on songs.id = shares.song_id").
			Where("songs.id IN ? AND shares.user_id = ?", songIDs, userID).
			Delete(&models.Share{}).
			Error

		if err != nil {
			c.JSON(400, gin.H{"error": "failed to delete shares"})
			return
		}

		c.JSON(200, gin.H{"data": "successfully deleted shares"})
	}
}
