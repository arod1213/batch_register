package spotify

import "fmt"

func getAlbum(href string, auth *auth) (*album, error) {
	return getModel[album](href, auth)
}

func getAlbumById(id string, auth *auth) (*album, error) {
	href := fmt.Sprintf("https://api.spotify.com/v1/albums/%s", id)
	return getModel[album](href, auth)
}

func getAlbumTracks(id string, auth *auth) (*albumItems, error) {
	href := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks", id)
	return getModel[albumItems](href, auth)
}
