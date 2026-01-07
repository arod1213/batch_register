package services

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/arod1213/auto_ingestion/genius"
	"github.com/arod1213/auto_ingestion/models"
	"gorm.io/gorm"
)

func IdentifyUser(db *gorm.DB, user models.User, songs []models.Song) (models.User, error) {
	artistMap := make(map[uint]genius.Artist)
	for _, song := range songs {
		keyword := fmt.Sprintf("%s %s", song.Artist, song.Title)
		fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		artists, err := genius.GeniusSearchArtists(keyword, fullName)
		if err != nil {
			continue
		}
		for _, artist := range artists {
			artistMap[artist.ID] = artist
		}
		if len(artists) == 1 {
			break // prevent unnecessary loops
		}
	}
	allArtists := slices.Collect(maps.Values(artistMap))
	if len(allArtists) != 1 {
		return user, errors.New("too many matches please provide a different song")
	}
	user.GeniusID = &allArtists[0].ID
	err := db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
