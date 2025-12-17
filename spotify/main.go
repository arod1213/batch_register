package spotify

import (
	"fmt"
	"log"
	"sync"

	"github.com/arod1213/auto_ingestion/models"
)

func AlbumToTracks(id string) []models.Song {
	var x []models.Song
	auth, err := getAuth()
	if err != nil {
		fmt.Println("auth error: ", err)
		return x
	}

	a, err := getAlbumById(id, auth)
	if err != nil {
		fmt.Println("album error: ", err)
		return x
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Println("found items of", len(a.Tracks.Items))
	tracksInfo, err := getAlbumTracks(id, auth)
	if err != nil {
		return x
	}

	for _, item := range tracksInfo.Tracks {
		i := item
		wg.Add(1)
		go func() {
			defer wg.Done()
			song := i.toSong()
			a.updateSong(&song)

			mu.Lock()
			x = append(x, song)
			mu.Unlock()
		}()
	}

	wg.Wait()
	return x
}

func ArtistToTracks(id string) []models.Song {
	var x []models.Song
	auth, err := getAuth()
	if err != nil {
		fmt.Println("auth error: ", err)
		return x
	}
	a, err := getArtistAlbums(id, auth)
	if err != nil {
		fmt.Println("artist album error: ", err)
		return x
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // prevent data race on save songs

	for _, item := range a.Items {
		wg.Add(1)

		a, err := getAlbumById(item.Id, auth)
		if err != nil {
			continue
		}
		go func() {
			tracksInfo, err := getAlbumTracks(item.Id, auth)
			if err != nil {
				log.Println("album tracks error: ", err.Error())
				return
			}
			for _, item := range tracksInfo.Tracks {
				song := item.toSong()
				a.updateSong(&song)

				mu.Lock()
				x = append(x, song)
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return x
}

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
				return
			}
			album.updateSong(&song)
			x = append(x, song)
		}()
	}
	wg.Wait()

	return x
}
