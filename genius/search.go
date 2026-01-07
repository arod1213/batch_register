package genius

import (
	"encoding/json"
	"log"
	"maps"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/arod1213/auto_ingestion/utils"
)

type searchResponse struct {
	Hits []GeniusSearchHit `json:"hits"`
}

type GeniusSearchHit struct {
	Index  string `json:"index"`
	Type   string `json:"type"`
	Result Song   `json:"result"`
}

func GeniusSearch(keyword string) ([]GeniusSearchHit, error) {
	accessToken := os.Getenv("GENIUS_ACCESS_TOKEN")
	href := url.URL{
		Scheme: "https",
		Host:   "api.genius.com",
		Path:   "/search",
		RawQuery: url.Values{
			"q": {keyword},
		}.Encode(),
	}

	body, err := utils.FetchBody(href.String(), accessToken, "GET", url.Values{})
	if err != nil {
		return nil, err
	}

	var response wrapper[searchResponse]
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response.Response.Hits, nil
}

func GeniusSearchArtists(keyword string, artistName string) ([]*Artist, error) {
	hits, err := GeniusSearch(keyword)
	if err != nil {
		return nil, err
	}

	matchArtist := func(song Song) (bool, *Artist) {
		lowerArtistName := strings.ToLower(artistName)
		for _, writer := range song.WriterArtists {
			if strings.ToLower(writer.Name) == lowerArtistName {
				return true, &writer
			}
		}
		for _, producer := range song.ProducerArtists {
			if strings.ToLower(producer.Name) == lowerArtistName {
				return true, &producer
			}
		}
		if strings.ToLower(song.PrimaryArtist.Name) == lowerArtistName {
			return true, &song.PrimaryArtist
		}
		return false, nil
	}

	wg := sync.WaitGroup{}
	ch := make(chan []*Artist, len(hits))
	for _, hit := range hits {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			artists, err := CollectContributors(uint(hit.Result.ID), matchArtist)
			if err != nil {
				log.Println("error collecting contributors: ", err)
				return
			}
			ch <- artists
		}(hit.Result.ID)
	}
	wg.Wait()
	close(ch)

	artistMap := make(map[int]*Artist)
	for match := range ch {
		for _, artist := range match {
			artistMap[artist.ID] = artist
		}
	}

	artists := slices.Collect(maps.Values(artistMap))
	return artists, nil
}
