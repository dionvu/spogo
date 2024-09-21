package ui

import (
	"fmt"

	"github.com/dionvu/spogo/spotify"
	"github.com/fatih/color"
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

	SEARCH_TYPE_VIEW     = "SEARCH_TYPE_VIEW"
	SEARCH_PLAYLIST_VIEW = "SEARCH_PLAYLIST_VIEW"
	SEARCH_TRACK_VIEW    = "SEARCH_TRACK_VIEW"
	SEARCH_ALBUM_VIEW    = "SEARCH_ALBUM_VIEW"

	DEVICE_VIEW = "DEVICE_VIEW"
)

func (m *Model) View() string {
	switch m.CurrentView {
	case PLAYER_VIEW:
		return m.Views.Player.View(m.Terminal.Height)

	case PLAYLIST_VIEW:
		return m.Views.Playlist.View(m.Views.Player, m.Terminal.Height)

	case HELP_VIEW:
		return MainControlsView(HELP_VIEW) + "\n\n" + padLines(m.Config.HelpString(), TAB_WIDTH)

	case PLAYLIST_TRACK_VIEW:
		state := m.Views.Player.State.CurrentPlayingType

		switch state {
		case "track":
			selectedItem := m.Views.Playlist.PlaylistListModel.list.SelectedItem()

			playlistName := m.Views.Playlist.ItemsMap[selectedItem]

			playlist := m.Views.Playlist.playlistsMap[playlistName]

			t, err := spotify.PlaylistTracks(m.Session, playlist.ID)
			if t == nil || err != nil || len(*t) < 1 {
				m.CurrentView = PLAYLIST_VIEW
				return ""
			}

			tracks := *t

			idx, _ := fuzzyfinder.FindMulti(
				tracks,
				func(i int) string {
					return tracks[i].Name + " - " + tracks[i].Artists[0].Name
				},
				fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
					if i == -1 {
						return ""
					}

					mins, secs := msToMinutesAndSeconds(tracks[i].DurationMs)

					return fmt.Sprintf("Track: %s \nArtist: %s\nAlbum: %s\nDuration: %sm:%ss",
						tracks[i].Name,
						tracks[i].Artists[0].Name,
						tracks[i].Album.Name,
						mins,
						secs,
					)
				}))

			if len(idx) > 0 && idx[0] < len(tracks) {
				m.Player.Play("", tracks[idx[0]].Uri, m.Session)
			}

			m.CurrentView = PLAYLIST_VIEW

			// Hide cursor after fzf showing fuzzyfinder
			fmt.Print("\033[?25l")

			return "changing views..."

		case "episode":
			return "TODO"

		default:
			return ""
		}

	case ALBUM_TRACK_VIEW:
		albumID := m.Views.Player.State.Track.Album.ID

		t, _ := spotify.AlbumTracks(m.Session, albumID)

		tracks := *t

		idx, _ := fuzzyfinder.FindMulti(
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

		if len(idx) > 0 && idx[0] < len(tracks) {
			m.Player.Play("", tracks[idx[0]].Uri, m.Session)
		}

		fmt.Print("\033[?25l") // Hide cursor after fzf showing fuzzyfinder

		m.CurrentView = PLAYER_VIEW

		return ""

	case REFRESH_VIEW:
		return "Refreshing..."

	case TERMINAL_WARNING_VIEW:
		return color.RedString(fmt.Sprint("Terminal of size ", m.Terminal.Height, "x", m.Terminal.Width, " is prone to visual glitches.\nMinimum required height is ", MIN_TERMINAL_HEIGHT, "."))

	case SEARCH_TYPE_VIEW:
		return m.Views.SearchType.View(m.Views.Player, m.Terminal)

	case SEARCH_ALBUM_VIEW:
		return "TODO album"

	case SEARCH_TRACK_VIEW:
		return "TODO track"

	case SEARCH_PLAYLIST_VIEW:
		return "TODO playlist"

	case DEVICE_VIEW:
		return "todo"

	default:
		return "TODO"
	}
}

func msToMinutesAndSeconds(ms int) (string, string) {
	minutes := ms / 60000
	seconds := (ms % 60000) / 1000

	return fmt.Sprint(minutes), fmt.Sprint(seconds)
}
