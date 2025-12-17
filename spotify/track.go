package spotify

import "fmt"

func getTrack(id string, auth *auth) (*track, error) {
	href := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", id)
	return getModel[track](href, auth)
}
