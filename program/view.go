package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/views"
	"github.com/ktr0731/go-fuzzyfinder"
)

func (p *Program) View() string {
	switch p.CurrentView {
	case views.PLAYER_VIEW:
		return p.Player.View(p.Terminal)

	case views.PLAYLIST_VIEW:
		return p.Playlist.View(p.Player, p.Terminal)

	case views.HELP_VIEW:
		return "idk man just press them keys"

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
			p.player.Resume(p.session, false)
		}

		p.CurrentView = views.PLAYER_VIEW

	case views.PLAYLIST_TRACK_VIEW:
		playlist := p.Playlist.GetSelectedPlaylist()

		if playlist == nil {
			p.CurrentView = views.PLAYLIST_VIEW
			return ""
		}

		t, err := spotify.PlaylistTracks(p.session, playlist.ID)
		if t == nil || err != nil || len(*t) < 1 {
			p.CurrentView = views.PLAYLIST_VIEW
			return ""
		}

		tracks := *t
		// asciis := make([]string, len(tracks))

		p.CurrentView = views.PLAYER_VIEW

		// cd, _ := os.UserCacheDir()
		// imagePath := filepath.Join(cd, config.APPNAME, playlist.ID, "temp")
		// os.MkdirAll(imagePath, os.ModePerm)
		// for i, track := range tracks {
		// 	image := components.Image{FilePath: fmt.Sprint(imagePath, i, ".jpeg")}
		// 	image.Update(track.Album.Images[0].Url)
		// 	asciis[i] = image.AsciiSmall().String()
		// }

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

				mins, secs := views.MsToMinutesAndSeconds(tracks[i].DurationMs)

				return fmt.Sprintf("Track: %s \nArtist: %s\nAlbum: %s\nDuration: %sm:%ss", //\n\n%s",
					tracks[i].Name,
					tracks[i].Artists[0].Name,
					tracks[i].Album.Name,
					mins,
					secs,
					// asciis[i],
				)
			}))

		contextUri := p.Playlist.GetSelectedPlaylist().Uri

		// Prevents user pressing Esc from playing the first track.
		if err == nil {
			p.player.Play(contextUri, tracks[idx].Uri, p.session)
		}

		HideCursor()

		if err == nil {
			p.CurrentView = views.PLAYER_VIEW
		} else {
			p.CurrentView = views.PLAYLIST_VIEW
		}

		return ""

	case views.ALBUM_TRACK_VIEW:
		if p.Player.State == nil || p.Player.State.Track == nil {
			p.CurrentView = views.PLAYER_VIEW
			return ""
		}

		album := &p.Player.State.Track.Album

		t, _ := spotify.AlbumTracks(p.session, album.ID)
		tracks := *t

		// Fzf tracks from the album currently playing
		// and plays the selected track.
		idx, err := FzfAlbumTracks(t)
		if err == nil {
			err := p.player.Play(album.Uri, tracks[idx].Uri, p.session)
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

		return ""

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
				return ""
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

// Hides the user's cursor after fzf.
func HideCursor() {
	fmt.Print("\033[?25l")
}
