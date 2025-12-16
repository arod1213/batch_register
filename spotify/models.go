package spotify

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arod1213/auto_ingestion/models"
)

type auth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
	Error       string `json:"error"`
}

type playlist struct {
	TracksInfo tracksInfo `json:"tracks"`
}
type tracksInfo struct {
	Next  string         `json:"next"`
	Total uint           `json:"total"`
	Items []playlistItem `json:"items"`
}

type playlistItem struct {
	Track track `json:"track"`
}

type trackAlbum struct {
	Href string `json:"href"`
}

type album struct {
	ExternalIds externalIDs `json:"external_ids"`
	Label       string      `json:"label"`
	ReleaseDate PartialDate `json:"release_date"`
}

func (a album) updateSong(s *models.Song) {
	var upc uint64
	if a.ExternalIds.UPC != nil {
		x, err := strconv.ParseUint(*(a.ExternalIds.UPC), 10, 64)
		if err == nil {
			upc = x
		}
	}
	s.ReleaseDate = a.ReleaseDate.Time
	s.Upc = upc
	s.Label = a.Label
}

type track struct {
	Album        trackAlbum   `json:"album"`
	Artists      []artist     `json:"artists"`
	Name         string       `json:"name"`
	ExternalIds  externalIDs  `json:"external_ids"`
	Duration     uint32       `json:"duration_ms"`
	ExternalUrls externalUrls `json:"external_urls"`
}

func (t track) toSong() models.Song {
	var artist string
	if len(t.Artists) > 0 {
		artist = t.Artists[0].Name
	}
	var isrc string
	if t.ExternalIds.ISRC != nil {
		isrc = *(t.ExternalIds.ISRC)
	}

	dStr := fmt.Sprintf("%vms", t.Duration)
	duration, _ := time.ParseDuration(dStr)

	return models.Song{
		Url:      t.ExternalUrls.Spotify,
		Duration: duration,
		Title:    t.Name,
		Artist:   artist,
		Isrc:     isrc,
	}
}

type externalUrls struct {
	Spotify string `json:"spotify"`
}

type externalIDs struct {
	UPC  *string `json:"upc"`
	ISRC *string `json:"isrc"`
}

type artist struct {
	Name string `json:"name"`
}
