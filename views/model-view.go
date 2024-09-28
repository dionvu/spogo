package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
	"github.com/ktr0731/go-fuzzyfinder"
)

const (
	PLAYER_VIEW           = "PLAYER_VIEW"
	PLAYLIST_VIEW         = "PLAYLIST_VIEW"
	PLAYLIST_TRACK_VIEW   = "PLAYLIST_TRACK_VIEW"
	ALBUM_TRACK_VIEW      = "ALBUM_TRACK_VIEW"
	REFRESH_VIEW          = "REFRESH_VIEW"
	HELP_VIEW             = "HELP_VIEW"
	TERMINAL_WARNING_VIEW = "TERMINAL_WARNING_VIEW"

	SEARCH_TYPE_VIEW  = "SEARCH_TYPE_VIEW"
	SEARCH_QUERY_VIEW = "SEARCH_QUERY_VIEW"

	SEARCH_PLAYLIST_VIEW = "SEARCH_PLAYLIST_VIEW"
	SEARCH_TRACK_VIEW    = "SEARCH_TRACK_VIEW"
	SEARCH_ALBUM_VIEW    = "SEARCH_ALBUM_VIEW"

	DEVICE_VIEW = "DEVICE_VIEW"
)

func (m *Model) View() string {
	switch m.CurrentView {
	case PLAYER_VIEW:
		return m.Views.Player.View(m.Terminal)

	case PLAYLIST_VIEW:
		return m.Views.Playlist.View(m.Views.Player, m.Terminal)

	case HELP_VIEW:
		return "\n\n" + MainControlsRender(HELP_VIEW) + "\n\n" + padLines(HelpString(), TAB_WIDTH)

	case PLAYLIST_TRACK_VIEW:
		state := m.Views.Player.State.CurrentPlayingType

		switch state {
		case "track":
			selectedItem := m.Views.Playlist.PlaylistListModel.list.SelectedItem()

			playlistName := m.Views.Playlist.ItemsMap[selectedItem]

			playlist := m.Views.Playlist.playlistsMap[playlistName]

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

					mins, secs := msToMinutesAndSeconds(tracks[i].DurationMs)

					cd, _ := os.UserCacheDir()
					imagePath := filepath.Join(cd, config.APPNAME, "temp.jpeg")

					_ = cacheImage(tracks[i].Album.Images[0].Url, imagePath)

					var ascii string
					// if m.Terminal.IsSizeSmall() {
					// 	ascii, _ = AsciiRender(imagePath, AsciiFlagsSmall())
					// } else {
					// 	ascii, _ = AsciiRender(imagePath, AsciiFlagsNormal())
					// }

					return fmt.Sprintf("Track: %s \nArtist: %s\nAlbum: %s\nDuration: %sm:%ss\n\n%s",
						tracks[i].Name,
						tracks[i].Artists[0].Name,
						tracks[i].Album.Name,
						mins,
						secs,
						ascii,
					)
				}))

			contextUri := m.Views.Playlist.GetSelectedPlaylist().URI

			// Prevents user pressing Esc from playing the first track.
			if err == nil {
				m.Player.Play(contextUri, tracks[idx].Uri, m.Session)
			}

			m.CurrentView = PLAYLIST_VIEW

			fmt.Print("\033[?25l") // Hide cursor after fzf showing fuzzyfinder

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
		idx, err := fuzzyfinder.Find(
			tracks,
			func(i int) string {
				return tracks[i].Name
			},
			fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
				if i == -1 {
					return ""
				}

				mins, secs := msToMinutesAndSeconds(tracks[i].DurationMs)

				return fmt.Sprintf("Track: %s\nDuration: %sm:%ss",
					tracks[i].Name,
					mins,
					secs,
				)
			}))

		// Prevents user pressing Esc from playing the first track.
		if err == nil {
			m.Player.Play(album.Uri, tracks[idx].Uri, m.Session)
		}

		fmt.Print("\033[?25l") // Hide cursor after fzf showing fuzzyfinder

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

// Converts the number of milliseconds into two string values
// of minutes and addittional seconds.
func msToMinutesAndSeconds(ms int) (minutes string, seconds string) {
	m := ms / 60000
	s := (ms % 60000) / 1000

	minutes = fmt.Sprint(m)
	seconds = fmt.Sprint(s)

	return minutes, seconds
}
