package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)
		count := 0

		mlcFile, err := models.MLCWrite(shares, *user)
		if err == nil {
			f2, err := zipWriter.Create("mlc.xlsx")
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			_, err = f2.Write(mlcFile.Bytes())
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			count++
		}

		sxFile, err := models.SXWrite(shares)
		if err == nil {
			f1, err := zipWriter.Create("sx.xlsx")
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			_, err = f1.Write(sxFile.Bytes())
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			count++
		}

		err = zipWriter.Close()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Close failed: %v", err)})
			return
		}

		if count == 0 {
			c.JSON(500, gin.H{"error": "No files added"})
			return
		}

		c.Header("Content-Disposition", `attachment; filename="tracks.zip"`)
		c.Data(200, "application/zip", buf.Bytes())
	}
}
