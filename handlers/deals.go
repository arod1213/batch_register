package handlers

import (
	"strconv"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteDeal(db *gorm.DB, isMaster bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		dealStr := c.Param("dealID")
		dealID, err := strconv.ParseUint(dealStr, 10, 32)
		if err != nil {
			c.JSON(440, gin.H{"error": "bad param"})
			return
		}

		if isMaster {
			err = db.Where("id = ?", uint(dealID)).Delete(&models.MasterDeal{}).Error
			if err != nil {
				c.JSON(500, gin.H{"error": "failed to delete"})
				return
			}
		} else {
			err = db.Where("id = ?", uint(dealID)).Delete(&models.PubDeal{}).Error
			if err != nil {
				c.JSON(500, gin.H{"error": "failed to delete"})
				return
			}
		}
	}
}

func bindAndSave[T models.MasterDeal | models.PubDeal](c *gin.Context, db *gorm.DB, songID uint) error {
	var vals []models.MasterDeal
	err := c.ShouldBindBodyWithJSON(&vals)
	if err != nil {
		c.JSON(440, gin.H{"error": "invalid body"})
		return err
	}
	for i := range vals {
		vals[i].SongID = uint(songID)
	}
	err = db.Save(&vals).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save"})
		return err
	}
	return nil
}

func CreateDeals(db *gorm.DB, isMaster bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		songStr := c.Param("songID")
		songID, err := strconv.ParseUint(songStr, 10, 32)
		if err != nil {
			c.JSON(440, gin.H{"error": "bad param"})
			return
		}

		if isMaster {
			err = bindAndSave[models.MasterDeal](c, db, uint(songID))
			if err != nil {
				return
			}
		} else {
			err = bindAndSave[models.PubDeal](c, db, uint(songID))
			if err != nil {
				return
			}
		}
		c.JSON(200, gin.H{"data": "success"})
	}

}
