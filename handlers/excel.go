package handlers

import (
	"fmt"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DownloadAllShares(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := middleware.GetUser(c, db)
		if err != nil {
			c.JSON(401, gin.H{"err": "unauthorized"})
			return
		}
		var shares []models.Share
		err = db.Where("user_id = ?", user.ID).Find(&shares).Error
		if err != nil {
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}
		data, err := services.WriteShares(shares, *user)
		if err != nil {
			c.JSON(500, gin.H{"err": err.Error()})
		}
		c.Header("Content-Disposition", `attachment; filename="tracks.zip"`)
		c.Data(200, "application/zip", *data)
	}
}

func DownloadRegistrations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var shares []models.Share
		err := c.ShouldBindBodyWithJSON(&shares)
		if err != nil {
			fmt.Println("err is ", err.Error())
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}

		user, err := middleware.GetUser(c, db)
		if err != nil {
			c.JSON(401, gin.H{"err": "unauthorized"})
			return
		}

		go func() {
			tx := db.Begin()
			for _, share := range shares {
				share.Song.Registered = true

				err := tx.Save(&share.Song).Error
				if err != nil {
					fmt.Println("err is ", err.Error())
					tx.Rollback()
					return
				}
			}
			tx.Commit()
		}()

		data, err := services.WriteShares(shares, *user)
		if err != nil {
			c.JSON(500, gin.H{"err": err.Error()})
			return
		}

		c.Header("Content-Disposition", `attachment; filename="tracks.zip"`)
		c.Data(200, "application/zip", *data)
	}
}
