package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dionvu/spogo/components"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/utils"
	"github.com/dionvu/spogo/views"
	"github.com/ktr0731/go-fuzzyfinder"
)

func (p *Program) View() string {
	switch p.CurrentView {
	case PLAYER_VIEW:
		return p.Player.View(p.Terminal)

	case PLAYLIST_VIEW:
		return p.Playlist.View(p.Player, p.Terminal)

	case HELP_VIEW:
		// return HelpString()
		return ""

	case PLAYLIST_TRACK_VIEW:
		playlist := p.Playlist.GetSelectedPlaylist()

		if playlist == nil {
			p.CurrentView = PLAYLIST_VIEW
			return ""
		}

		t, err := spotify.PlaylistTracks(p.session, playlist.ID)
		if t == nil || err != nil || len(*t) < 1 {
			p.CurrentView = PLAYLIST_VIEW
			return ""
		}

		tracks := *t
		asciis := make([]string, len(tracks))

		p.CurrentView = PLAYER_VIEW

		cd, _ := os.UserCacheDir()
		imagePath := filepath.Join(cd, config.APPNAME, playlist.ID, "temp")
		os.MkdirAll(imagePath, os.ModePerm)
		for i, track := range tracks {
			image := components.Image{FilePath: fmt.Sprint(imagePath, i, ".jpeg")}
			image.Update(track.Album.Images[0].Url)
			asciis[i] = image.AsciiSmall().String()
		}

		// Fzf tracks from the playlist currently slected
		// and plays the selected track.
		idx, err := fuzzyfinder.Find(
			tracks,
			func(i int) string {
				return tracks[i].Name + " - " + tracks[i].Artists[0].Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}

				mins, secs := utils.MsToMinutesAndSeconds(tracks[i].DurationMs)

				return fmt.Sprintf("Track: %s \nArtist: %s\nAlbum: %s\nDuration: %sm:%ss\n\n%s",
					tracks[i].Name,
					tracks[i].Artists[0].Name,
					tracks[i].Album.Name,
					mins,
					secs,
					asciis[i],
				)
			}))

		contextUri := p.Playlist.GetSelectedPlaylist().URI

		// Prevents user pressing Esc from playing the first track.
		if err == nil {
			p.player.Play(contextUri, tracks[idx].Uri, p.session)
		}

		HideCursor()

		if err == nil {
			p.CurrentView = PLAYER_VIEW
		} else {
			p.CurrentView = PLAYLIST_VIEW
		}

		return ""

	case ALBUM_TRACK_VIEW:
		if p.Player.State == nil || p.Player.State.Track == nil {
			p.CurrentView = PLAYER_VIEW
			return ""
		}

		album := &p.Player.State.Track.Album

		t, _ := spotify.AlbumTracks(p.session, album.ID)
		tracks := *t

		// Fzf tracks from the album currently playing
		// and plays the selected track.
		idx, err := FzfAlbumTracks(t)
		if err == nil {
			p.player.Play(album.Uri, tracks[idx].Uri, p.session)

			go func() {
				// After playing spotify takes a moment to update the state
				// to match the newly played song.
				time.Sleep(time.Second)

				// Syncs state to new song.
				p.Player.UpdateStateSync()
			}()
		}

		p.CurrentView = PLAYER_VIEW

		return ""

	case REFRESH_VIEW:
		return "Refreshing..."

	case TERMINAL_WARNING_VIEW:
		return p.Terminal.WarningString()

		// case SEARCH_TYPE_VIEW:
		// 	return p.SearchType.View(p.Views.Player, p.Terminal)

		// case SEARCH_QUERY_VIEW:
		// 	return p.Squery.View()

		// case SEARCH_RESULT_TRACK:
		// 	fmt.Println(p.SearchType.ListModel.list.SelectedItem())
		// 	fmt.Println(p.SearchType.itemsMap[p.Views.SearchType.ListModel.list.SelectedItem()])
		// 	p.SearchResult = NewSearchResultView(p.Views.Squery.Query(), p.Views.SearchType.itemsMap[p.Views.SearchType.ListModel.list.SelectedItem()], p.session)
		// 	return p.SearchResult.items.Tracks[0].Name

	case views.SEARCH_VIEW_QUERY, views.SEARCH_VIEW_TYPE, views.SEARCH_VIEW_RESULTS:
		return p.Search.View(p.Terminal, p.CurrentView)

	case DEVICE_VIEW:
		return p.Device.View(p.Terminal, p.Player.State.Device)
	}

	return "UNREACHABLE"
}

func FzfAlbumTracks(albumTracks *[]spotify.AlbumTrack) (int, error) {
	tracks := *albumTracks

	idx, err := fuzzyfinder.Find(
		tracks,
		func(i int) string {
			return tracks[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			mins, secs := utils.MsToMinutesAndSeconds(tracks[i].DurationMs)

			return fmt.Sprintf("Track: %s\nDuration: %sm:%ss",
				tracks[i].Name,
				mins,
				secs,
			)
		}))

	HideCursor()

	return idx, err
}

// Hides the user's cursor after fzf.
func HideCursor() {
	fmt.Print("\033[?25l")
}
