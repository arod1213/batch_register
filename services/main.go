package services

import (
	"slices"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/spotify"
	"github.com/arod1213/auto_ingestion/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetSpotifyTracks(method string, id string) []models.Song {
	var songs []models.Song
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
	return songs
}

func SaveSongs(db *gorm.DB, userID uint, songs []models.Song) ([]models.Share, error) {
	tx := db.Begin()
	shares := make([]models.Share, len(songs))

	err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&songs).Error
	if err != nil {
		tx.Rollback()
		return shares, err
	}

	isrcs := utils.Map(songs, func(s models.Song) string {
		return s.Isrc
	})

	var insertedSongs []models.Song
	err = tx.Where("isrc IN ?", isrcs).Find(&insertedSongs).Error
	if err != nil {
		tx.Rollback()
		return shares, err
	}

	for i, song := range insertedSongs {
		shares[i].SongID = song.ID
		shares[i].UserID = userID
	}

	err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&shares).Error
	if err != nil {
		tx.Rollback()
		return shares, err
	}

	tx.Commit()
	return shares, nil
}
