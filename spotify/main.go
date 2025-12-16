package spotify

import (
	"fmt"
	"sync"

	"github.com/arod1213/auto_ingestion/models"
)

func PlaylistToTracks(playlist string) []models.Song {
	var x []models.Song
	auth, err := getAuth()
	if err != nil {
		fmt.Println("err is ", err)
		return x
	}
	p, err := getPlaylist(playlist, auth)
	if err != nil {
		fmt.Println("err is ", err)
		return x
	}

	var wg sync.WaitGroup
	for _, t := range p.TracksInfo.Items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			song := t.Track.toSong()
			album, err := getAlbum(t.Track.Album.Href, auth)
			if err != nil {
				fmt.Println("err is ", err)
				x = append(x, song)
				return
			}
			album.updateSong(&song)
			x = append(x, song)
		}()
	}
	wg.Wait()

	return x
}
