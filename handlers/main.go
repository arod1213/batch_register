package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"slices"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/spotify"
	"github.com/arod1213/auto_ingestion/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchTracks(c *gin.Context, db *gorm.DB) {
	var songs []models.Song

	id := c.Param("id")
	method := c.Query("method")

	switch method {
	case "artist":
		songs = spotify.ArtistToTracks(id)
	case "album":
		fmt.Println("searching album ", id)
		songs = spotify.AlbumToTracks(id)
	case "playlist":
		songs = spotify.PlaylistToTracks(id)
		fmt.Println("found ", len(songs))
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

	go func() {
		err := SaveSongs(db, songs) // async call
		if err != nil {
			log.Println("error saving songs")
		}
	}()

	c.JSON(200, gin.H{"data": songs})
}

func SaveSongs(db *gorm.DB, songs []models.Song) error {
	return db.Save(&songs).Error
}

func UpdateSongs(db *gorm.DB, songs []models.Song) error {
	return db.Save(&songs).Error
}

func WriteTracks(c *gin.Context, db *gorm.DB) {
	var data []models.Info
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		fmt.Println("err is ", err.Error())
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	go func() {
		songs := utils.Map(data, func(info models.Info) models.Song {
			s := info.Song
			s.Registered = true
			return s
		})
		err := SaveSongs(db, songs)
		if err != nil {
			log.Println("error saving songs")
		}
	}()

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	count := 0

	mlcFile, err := models.MLCWrite(data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	sxFile, err := models.SXWrite(data)
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
