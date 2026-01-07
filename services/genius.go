package services

import (
	"github.com/arod1213/auto_ingestion/genius"
	"github.com/arod1213/auto_ingestion/models"
	"gorm.io/gorm"
)

func GetMissingSongs(db *gorm.DB, artistId uint, accessToken string, userID uint) ([]genius.Song, error) {
	songs, err := genius.GetArtistSongs(artistId, accessToken)
	if err != nil {
		return nil, err
	}

	var missingSongs []genius.Song
	for _, song := range songs {
		var share models.Share
		err := db.
			Joins("INNER JOIN songs on songs.id = shares.song_id").
			Where("shares.user_id = ?", userID).
			Where("songs.title LIKE ? AND songs.artist LIKE ?", "%"+song.Title+"%", "%"+song.PrimaryArtist.Name+"%").
			First(&share).
			Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				missingSongs = append(missingSongs, song)
			} else {
				continue
			}
		}
	}

	return missingSongs, nil
}
