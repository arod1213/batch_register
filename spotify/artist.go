package spotify

import (
	"fmt"
)

func getArtistAlbums(id string, auth *auth) (*artistAlbums, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%v/albums?limit=5", id)
	return getModel[artistAlbums](url, auth)
}
