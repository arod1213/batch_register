package genius

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/arod1213/auto_ingestion/utils"
)

type SongResponse struct {
	Song Song `json:"song"`
}

func GetSong(songID uint) (*Song, error) {
	accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")
	href := fmt.Sprintf("https://api.genius.com/songs/%d", songID)
	body, err := utils.FetchBody(href, accessToken, "GET", url.Values{})
	if err != nil {
		return nil, err
	}
	var res wrapper[SongResponse]
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res.Response.Song, nil
}

func CollectContributors(songID uint, matchArtist func(Song) (bool, *Artist)) ([]*Artist, error) {
	song, err := GetSong(songID)
	if err != nil {
		return nil, err
	}
	// spotify.Pretty(song)
	artists := []*Artist{}

	found, artist := matchArtist(*song)
	if found {
		artists = append(artists, artist)
	}

	return artists, nil
}
