package spotify

import (
	"fmt"
)

func getPlaylist(playlistId string, auth *auth) (*playlist, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%v", playlistId)
	return getModel[playlist](url, auth)
}
