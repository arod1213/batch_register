package spotify

import "fmt"

func getAlbum(href string, auth *auth) (*album, error) {
	return getModel[album](href, auth)
}

func getAlbumById(id string, auth *auth) (*album, error) {
	href := fmt.Sprintf("https://api.spotify.com/v1/albums/%s", id)
	return getModel[album](href, auth)
}

func getAlbumTracks(id string, auth *auth, detailed bool) (*albumItems, error) {
	href := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks", id)

	items, err := getModel[albumItems](href, auth)
	if err != nil {
		return nil, err
	}

	if !detailed {
		return items, nil
	}

	// modify in place
	for i, track := range items.Tracks {
		t, err := getTrack(track.ID, auth)
		if err != nil {
			fmt.Println("failed to enrich ISRC for ", track.Name, track.ID)
			continue
		}
		items.Tracks[i].ExternalIds.ISRC = t.ExternalIds.ISRC
	}

	return items, nil
}
