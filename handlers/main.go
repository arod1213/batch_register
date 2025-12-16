package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"slices"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/spotify"
	"github.com/gin-gonic/gin"
)

func FetchTracks(c *gin.Context) {
	playlistID := c.Param("playlistID")
	fmt.Println("id is ", playlistID)
	tracks := spotify.PlaylistToTracks(playlistID)
	slices.SortFunc(tracks, func(x models.Song, y models.Song) int {
		return y.ReleaseDate.Compare(x.ReleaseDate)
	})
	c.JSON(200, gin.H{"data": tracks})
}

func WriteTracks(c *gin.Context) {
	var data []models.Info
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		fmt.Println("err is ", err.Error())
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

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
