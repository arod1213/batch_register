package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"slices"

	"github.com/arod1213/auto_ingestion/middleware"
	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/spotify"
	"github.com/arod1213/auto_ingestion/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FetchTracks(c *gin.Context, db *gorm.DB) {
	log.Println("FetchTracks - All keys:", c.Keys)
	var songs []models.Song

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	method := c.Query("method")

	switch method {
	case "artist":
		songs = spotify.ArtistToTracks(id)
	case "album":
		songs = spotify.AlbumToTracks(id)
	case "playlist":
		songs = spotify.PlaylistToTracks(id)
	default:
		songs = spotify.PlaylistToTracks(id)
	}

	slices.SortFunc(songs, func(x models.Song, y models.Song) int {
		return y.ReleaseDate.Compare(x.ReleaseDate)
	})

	if len(songs) == 0 {
		c.JSON(400, gin.H{"error": "no songs found: ensure your playlist is public"})
		return
	}

	tx := db.Begin()

	err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&songs).Error
	if err != nil {
		log.Println("error saving songs")
		tx.Rollback()
		c.JSON(400, gin.H{"error": "failed to save songs"})
		return
	}

	isrcs := utils.Map(songs, func(s models.Song) string {
		return s.Isrc
	})

	var insertedSongs []models.Song
	err = tx.Where("isrc IN ?", isrcs).Find(&insertedSongs).Error
	if err != nil {
		log.Println("error saving songs")
		tx.Rollback()
		c.JSON(400, gin.H{"error": "failed to save songs"})
		return
	}

	shares := make([]models.Share, len(songs))
	for i, song := range insertedSongs {
		shares[i].SongID = song.ID
		shares[i].UserID = userID
	}

	err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&shares).Error
	if err != nil {
		log.Println("error saving shares")
		tx.Rollback()
		c.JSON(400, gin.H{"error": "failed to save shares"})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{"data": shares})
}

func UpdateSongs(db *gorm.DB, songs []models.Song) error {
	return db.Save(&songs).Error
}

func WriteShares(c *gin.Context, db *gorm.DB) {
	var shares []models.Share
	err := c.ShouldBindBodyWithJSON(&shares)
	if err != nil {
		fmt.Println("err is ", err.Error())
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	go func() {
		tx := db.Begin()
		for _, share := range shares {
			share.Song.Registered = true

			err := tx.Save(&share).Error
			if err != nil {
				tx.Rollback()
			}

			// err = tx.Create(&share).Error
			// if err != nil {
			// 	tx.Rollback()
			// }
		}
		tx.Commit()
	}()

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	count := 0

	mlcFile, err := models.MLCWrite(shares)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	sxFile, err := models.SXWrite(shares)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

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
