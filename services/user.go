package services

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/arod1213/auto_ingestion/genius"
	"github.com/arod1213/auto_ingestion/models"
	"gorm.io/gorm"
)

func IdentifyUser(db *gorm.DB, user models.User, songs []models.Song) (models.User, error) {
	fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	log.Printf("searching for %s\n", fullName)

	artistMap := make(map[uint]genius.Artist)
	for _, song := range songs {
		log.Printf("looking at %s by %s\n", song.Title, song.Artist)
		keyword := fmt.Sprintf("%s %s", song.Artist, song.Title)
		artists, err := genius.GeniusSearchArtists(keyword, fullName)
		if err != nil {
			continue
		}
		log.Printf("found %d artists \n", len(artists))
		for _, artist := range artists {
			artistMap[artist.ID] = artist
		}
		if len(artists) == 1 {
			break // prevent unnecessary loops
		}
	}

	allArtists := slices.Collect(maps.Values(artistMap))

	switch len(allArtists) {
	case 0:
		return user, errors.New("no matches found")
	case 1:
		{
		}
	default:
		for _, artist := range allArtists {
			log.Println("found ", artist)
		}
		return user, errors.New("too many matches please provide a different song")
	}

	user.GeniusID = &allArtists[0].ID
	err := db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
