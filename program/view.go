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
	"github.com/ktr0731/go-fuzzyfinder"
)

func (m *Program) View() string {
	switch m.CurrentView {
	case PLAYER_VIEW:
		return m.Views.Player.View(m.Terminal)

	case PLAYLIST_VIEW:
		return m.Views.Playlist.View(m.Views.Player, m.Terminal)

	case HELP_VIEW:
		// return HelpString()
		return ""

	case PLAYLIST_TRACK_VIEW:
		playlist := m.Views.Playlist.GetSelectedPlaylist()

		if playlist == nil {
			m.CurrentView = PLAYLIST_VIEW
			return ""
		}

		t, err := spotify.PlaylistTracks(m.session, playlist.ID)
		if t == nil || err != nil || len(*t) < 1 {
			m.CurrentView = PLAYLIST_VIEW
			return ""
		}

		tracks := *t
		asciis := make([]string, len(tracks))

		m.CurrentView = PLAYER_VIEW

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

		contextUri := m.Views.Playlist.GetSelectedPlaylist().URI

		// Prevents user pressing Esc from playing the first track.
		if err == nil {
			m.player.Play(contextUri, tracks[idx].Uri, m.session)
		}

		HideCursor()

		if err == nil {
			m.CurrentView = PLAYER_VIEW
		} else {
			m.CurrentView = PLAYLIST_VIEW
		}

		return ""

	case ALBUM_TRACK_VIEW:
		if m.Views.Player.State == nil || m.Views.Player.State.Track == nil {
			m.CurrentView = PLAYER_VIEW
			return ""
		}

		album := &m.Views.Player.State.Track.Album

		t, _ := spotify.AlbumTracks(m.session, album.ID)
		tracks := *t

		// Fzf tracks from the album currently playing
		// and plays the selected track.
		idx, err := FzfAlbumTracks(t)
		if err == nil {
			m.player.Play(album.Uri, tracks[idx].Uri, m.session)

			go func() {
				// After playing spotify takes a moment to update the state
				// to match the newly played song.
				time.Sleep(time.Second)

				// Syncs state to new song.
				m.Views.Player.UpdateStateSync()
			}()
		}

		m.CurrentView = PLAYER_VIEW

		return ""

	case REFRESH_VIEW:
		return "Refreshing..."

	case TERMINAL_WARNING_VIEW:
		return m.Terminal.WarningString()

		// case SEARCH_TYPE_VIEW:
		// 	return m.Views.SearchType.View(m.Views.Player, m.Terminal)

		// case SEARCH_QUERY_VIEW:
		// 	return m.Views.Squery.View()

		// case SEARCH_RESULT_TRACK:
		// 	fmt.Println(m.Views.SearchType.ListModel.list.SelectedItem())
		// 	fmt.Println(m.Views.SearchType.itemsMap[m.Views.SearchType.ListModel.list.SelectedItem()])
		// 	m.Views.SearchResult = NewSearchResultView(m.Views.Squery.Query(), m.Views.SearchType.itemsMap[m.Views.SearchType.ListModel.list.SelectedItem()], m.session)
		// 	return m.Views.SearchResult.items.Tracks[0].Name

	case SEARCH_VIEW:

	case DEVICE_VIEW:
		return m.Views.Device.View(m.Terminal, m.Views.Player.State.Device)
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
