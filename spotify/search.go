package spotify

import (
	"fmt"
	"net/url"
	"slices"
)

func SearchRelated(auth *auth, t track) (*[]track, error) {
	var artist string
	if len(t.Artists) != 0 {
		artist = t.Artists[0].Name
	}

	query := url.QueryEscape(fmt.Sprintf("track:%s artist:%s&type=track&limit=5", t.Name, artist))
	href := fmt.Sprintf("https://api.spotify.com/v1/search?%s", query)
	tracks, error := getModel[[]track](href, auth)

	if error != nil {
		return nil, error
	}

	related := slices.DeleteFunc(*tracks, func(x track) bool {
		a := x.ExternalIds.ISRC
		b := t.ExternalIds.ISRC
		if a == nil || b == nil {
			return false
		}
		return *a == *b
	})

	return &related, nil
}
