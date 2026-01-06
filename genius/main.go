package genius

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/arod1213/auto_ingestion/utils"
)

type ArtistSongs struct {
	Response Response `json:"response"`
}

type Response struct {
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

	// Stats                 Stats    `json:"stats"`
	// AnnotationCount                           int                   `json:"annotation_count"`
	// HeaderImageThumbnailUrl                   string                `json:"header_image_thumbnail_url"`
	// HeaderImageUrl                            string                `json:"header_image_url"`
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
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
	ApiPath    string `json:"api_path"`
	Iq         int    `json:"iq"`
	// Url        string `json:"url"`
	// ImageUrl       string `json:"image_url"`
	// HeaderImageUrl string `json:"header_image_url"`
	// IsMemeVerified bool   `json:"is_meme_verified"`
}

func GetArtistSongs(artistId string) ([]Song, error) {
	href := fmt.Sprintf("https://api.genius.com/artists/%v/songs", artistId)
	accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")

	body, err := utils.FetchBody(href, accessToken, "GET", url.Values{})
	if err != nil {
		log.Println("error fetching body: ", err)
		return nil, err
	}

	var songs ArtistSongs
	err = json.Unmarshal(body, &songs)
	if err != nil {
		log.Println("error unmarshalling response: ", err)
		return nil, err
	}
	return songs.Response.Songs, nil
}
