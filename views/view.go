package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
		return HelpString()

	case PLAYLIST_TRACK_VIEW:
		state := m.Views.Player.State.CurrentPlayingType

		switch state {
		case "track":
			playlist := m.Views.Playlist.GetSelectedPlaylist()

			if playlist == nil {
				m.CurrentView = PLAYLIST_VIEW
				return ""
			}

			t, err := spotify.PlaylistTracks(m.Session, playlist.ID)
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
				image := Image{FilePath: fmt.Sprint(imagePath, i, ".jpeg")}
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
				m.Player.Play(contextUri, tracks[idx].Uri, m.Session)
			}

			fmt.Print("\033[?25l") // Hide cursor after fzf showing fuzzyfinder

			if err == nil {
				m.CurrentView = PLAYER_VIEW
			} else {
				m.CurrentView = PLAYLIST_VIEW
			}

			return ""

		case "episode":
			return "TODO"

		default:
			return ""
		}

	case ALBUM_TRACK_VIEW:
		if m.Views.Player.State == nil || m.Views.Player.State.Track == nil {
			m.CurrentView = PLAYER_VIEW
			return ""
		}

		album := &m.Views.Player.State.Track.Album

		t, _ := spotify.AlbumTracks(m.Session, album.ID)
		tracks := *t

		// Fzf tracks from the album currently playing
		// and plays the selected track.
		idx, err := FzfAlbumTracks(t)
		if err == nil {
			m.Player.Play(album.Uri, tracks[idx].Uri, m.Session)

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

	case SEARCH_TYPE_VIEW:
		return m.Views.SearchType.View(m.Views.Player, m.Terminal)

	case SEARCH_QUERY_VIEW:
		return m.Views.Squery.View()

	case SEARCH_ALBUM_VIEW:
		return "TODO album"

	case SEARCH_TRACK_VIEW:
		return "TODO track"

	case SEARCH_PLAYLIST_VIEW:
		return "TODO playlist"

	case DEVICE_VIEW:
		return m.Views.Device.View(m.Terminal, m.Views.Player.State.Device)

	default:
		return "TODO"
	}
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

	fmt.Print("\033[?25l") // Hide cursor after fzf showing fuzzyfinder

	return idx, err
}
