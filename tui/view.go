package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/tui/views"
	comp "github.com/dionvu/spogo/tui/views/components"
	"github.com/ktr0731/go-fuzzyfinder"
)

func (p *Program) View() string {
	switch p.CurrentView {
	case views.PLAYER_VIEW:
		if p.Player.State != nil && p.Player.State.CurrentPlayingType == views.EPISODE {
			return comp.Content("Does not support podcasts, but support is coming soon!").
				CenterVertical(p.Terminal).CenterHorizontal(p.Terminal).String()
		}

		return p.Player.View(p.Terminal)

	case views.PLAYLIST_VIEW:
		return p.Playlist.View(p.Player, p.Terminal)

	case views.HELP_VIEW:
		return p.Help.View()

	case views.REAUTH_VIEW:
		err := p.session.Reauth(p.Config)
		if err != nil {
			log.Fatal("ERR: Failed to reauthenticate: ", err)
			errors.Log(err)
		}
		p.CurrentView = views.PLAYER_VIEW

		return "reauthenticating..."

	case views.DEVICE_FZF_VIEW:
		devices, err := player.GetDevices(p.session)
		if errors.IsReauthenticationErr(err) {
			p.CurrentView = views.REAUTH_VIEW
		}

		idx, err := FzfDevices(devices)
		if err == nil {
			p.player.SetDevice(&(*devices)[idx], p.Config)

			if p.Player.State != nil && p.Player.State.IsPlaying {
				p.player.Resume(p.session, true)
			} else {
				p.player.Resume(p.session, false)
			}
		}

		p.CurrentView = views.PLAYER_VIEW

	case views.PLAYLIST_TRACK_VIEW:
		playlist := p.Playlist.GetSelectedPlaylist()

		if playlist == nil {
			p.CurrentView = views.PLAYLIST_VIEW
			return EMPTY
		}

		tracks, err := spotify.PlaylistTracks(p.session, playlist.ID)
		if tracks == nil || err != nil || len(*tracks) < 1 {
			p.CurrentView = views.PLAYLIST_VIEW
			return EMPTY
		}

		p.CurrentView = views.PLAYER_VIEW

		idx, err := FzfPlaylistTracks(tracks)

		if err == nil {
			p.CurrentView = views.PLAYER_VIEW
			p.player.Play(p.Playlist.GetSelectedPlaylist().Uri, (*tracks)[idx].Uri, p.session)
		} else {
			p.CurrentView = views.PLAYLIST_VIEW
		}

		return EMPTY

	case views.ALBUM_TRACK_VIEW:
		if p.Player.State == nil || p.Player.State.Track == nil {
			p.CurrentView = views.PLAYER_VIEW
			return EMPTY
		}

		album := &p.Player.State.Track.Album

		tracks, _ := spotify.AlbumTracks(p.session, album.ID)

		// Fzf tracks from the album currently playing
		// and plays the selected track.
		idx, err := FzfAlbumTracks(tracks)
		if err == nil {
			err := p.player.Play(album.Uri, (*tracks)[idx].Uri, p.session)
			if errors.IsReauthenticationErr(err) {
				p.CurrentView = views.REAUTH_VIEW
			}

			go func() {
				// After playing spotify takes a moment to update the state
				// to match the newly played song.
				time.Sleep(time.Second)

				// Syncs state to new song.
				p.Player.UpdateStateSync()
			}()
		}

		p.CurrentView = views.PLAYER_VIEW

		return EMPTY

	case views.REFRESH_VIEW:
		return "Refreshing..."

	case views.TERMINAL_WARNING_VIEW:
		return p.Terminal.WarningString()

	case views.SEARCH_VIEW_QUERY, views.SEARCH_VIEW_TYPE, views.SEARCH_VIEW_RESULTS:
		return p.Search.View(p.Terminal, p.CurrentView)

	case views.DEVICE_VIEW:
		if p.Player.State == nil {
			return p.Device.View(p.Terminal, nil, p.Config)
		}

		if p.player.Device() == nil {
			return p.Device.View(p.Terminal, nil, p.Config)
		}

		return p.Device.View(p.Terminal, p.player.Device(), p.Config)
	}

	return EMPTY
}

func FzfAlbumTracks(albumTracks *[]spotify.AlbumTrack) (int, error) {
	tracks := *albumTracks

	idx, err := fuzzyfinder.Find(
		tracks,
		func(i int) string {
			if tracks[i].Name == "" {
				return "Unavailable"
			}

			return tracks[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return EMPTY
			}

			if tracks[i].Name == "" {
				return "Content is unavailable :("
			}

			mins, secs := views.MsToMinutesAndSeconds(tracks[i].DurationMs)

			return fmt.Sprintf("Track: %s\nDuration: %sm:%ss",
				tracks[i].Name,
				mins,
				secs,
			)
		}))

	HideCursor()

	return idx, err
}

func FzfDevices(devices *[]player.Device) (int, error) {
	tracks := *devices

	idx, err := fuzzyfinder.Find(
		tracks,
		func(i int) string {
			return tracks[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return EMPTY
			}

			if tracks[i].Name == "" {
				return "Unavailable"
			}

			return fmt.Sprintf("Name: %s\nType: %s\nVol: %v%%",
				tracks[i].Name,
				tracks[i].Type,
				tracks[i].VolumePercent,
			)
		}))

	HideCursor()

	return idx, err
}

func FzfPlaylistTracks(t *[]spotify.Track) (int, error) {
	if t == nil {
		log.Fatal("Unreachable")
	}

	tracks := *t

	idx, err := fuzzyfinder.Find(
		tracks,
		func(i int) string {
			if tracks[i].Name == "" {
				return "Unavailable"
			}

			return tracks[i].Name + " - " + tracks[i].Artists[0].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return EMPTY
			}

			mins, secs := views.MsToMinutesAndSeconds(tracks[i].DurationMs)

			if tracks[i].Name == "" {
				return "Content is unavailable :("
			}

			return fmt.Sprintf("Track: %s \nArtist: %s\nAlbum: %s\nDuration: %sm:%ss", //\n\n%s",
				tracks[i].Name,
				tracks[i].Artists[0].Name,
				tracks[i].Album.Name,
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
