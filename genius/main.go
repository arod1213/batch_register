package genius

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/arod1213/auto_ingestion/utils"
)

type wrapper[T any] struct {
	Response T `json:"response"`
}

type ArtistSongs struct {
	Songs []Song `json:"songs"`
}

type Song struct {
	ApiPath               string   `json:"api_path"`
	ArtistNames           string   `json:"artist_names"`
	FullTitle             string   `json:"full_title"`
	ID                    int      `json:"id"`
	PrimaryArtistNames    string   `json:"primary_artist_names"`
	ReleaseDateForDisplay string   `json:"release_date_for_display"`
	Title                 string   `json:"title"`
	Url                   string   `json:"url"`
	FeaturedArtists       []Artist `json:"featured_artists"`
	PrimaryArtist         Artist   `json:"primary_artist"`
	PrimaryArtists        []Artist `json:"primary_artists"`
	HeaderImageUrl        string   `json:"header_image_url"`
	Missing               bool     `json:"missing"`
	ProducerArtists       []Artist `json:"producer_artists"`
	WriterArtists         []Artist `json:"writer_artists"`

	// Stats                 Stats    `json:"stats"`
	// AnnotationCount                           int                   `json:"annotation_count"`
	// HeaderImageThumbnailUrl                   string                `json:"header_image_thumbnail_url"`
	// LyricsOwnerId                             int                   `json:"lyrics_owner_id"`
	// LyricsState                               string                `json:"lyrics_state"`
	// Path                                      string                `json:"path"`
	// PyongsCount                               *int                  `json:"pyongs_count"`
	// RelationshipsIndexUrl                     string                `json:"relationships_index_url"`
	// ReleaseDateComponents                     ReleaseDateComponents `json:"release_date_components"`
	// ReleaseDateWithAbbreviatedMonthForDisplay string                `json:"release_date_with_abbreviated_month_for_display"`
	// SongArtImageThumbnailUrl                  string                `json:"song_art_image_thumbnail_url"`
	// SongArtImageUrl                           string                `json:"song_art_image_url"`
	// TitleWithFeatured                         string                `json:"title_with_featured"`
}

// type ReleaseDateComponents struct {
// 	Year  int `json:"year"`
// 	Month int `json:"month"`
// 	Day   int `json:"day"`
// }

// type Stats struct {
// 	UnreviewedAnnotations int  `json:"unreviewed_annotations"`
// 	Concurrents           *int `json:"concurrents"`
// 	Hot                   bool `json:"hot"`
// 	Pageviews             int  `json:"pageviews"`
// }

type Artist struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
	ApiPath    string `json:"api_path"`
	Iq         int    `json:"iq"`
	ImageUrl   string `json:"image_url"`
	// Url        string `json:"url"`
	// HeaderImageUrl string `json:"header_image_url"`
	// IsMemeVerified bool   `json:"is_meme_verified"`
}

func GetArtistSongs(artistId uint, accessToken string) ([]Song, error) {
	href := fmt.Sprintf("https://api.genius.com/artists/%d/songs", artistId)

	body, err := utils.FetchBody(href, accessToken, "GET", url.Values{})
	if err != nil {
		log.Println("error fetching body: ", err)
		return nil, err
	}

	var res wrapper[ArtistSongs]
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println("error unmarshalling response: ", err)
		return nil, err
	}
	return res.Response.Songs, nil
}
