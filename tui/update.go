package tui

import (
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/tui/views"
)

const (
	KEY_ESC      = "esc"
	KEY_ENTER    = "enter"
	KEY_QUIT     = "q"
	KEY_QUIT_ALT = "ctrl+c"

	KEY_VISUAL_REFRESH = "ctrl+r"

	KEY_PLAYER_VIEW      = "f1"
	KEY_PLAYER_VIEW_ALT  = "ctrl+v"
	KEY_PLAY_PAUSE       = " "
	KEY_TOGGLE_SHUFFLING = "s"
	KEY_TOGGLE_REPEAT    = "r"

	KEY_PLAYLIST_VIEW       = "f2"
	KEY_PLAYLIST_VIEW_ALT   = "ctrl+p"
	KEY_FZF_PLAYLIST_TRACKS = "t"

	KEY_SEARCH_VIEW     = "f3"
	KEY_SEARCH_VIEW_ALT = "/"

	KEY_DEVICE_VIEW = "UNDEFINED"

	KEY_HELP_VIEW     = "f4"
	KEY_HELP_VIEW_ALT = "ctrl+h"

	KEY_VOLUME_DOWN_BIG   = "["
	KEY_VOLUME_DOWN_SMALL = "{"
	KEY_VOLUME_UP_BIG     = "]"
	KEY_VOLUME_UP_SMALL   = "}"

	KEY_FZF_DEVICES      = "ctrl+d"
	KEY_FZF_ALBUM_TRACKS = "ctrl+a"

	VOLUME_INCREMENT_PERCENT = 5
)

// Handles updates associate with the current selected view.
func (p *Program) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p.Terminal.UpdateSize()

	if !p.Terminal.IsValid() {
		p.CurrentView = views.TERMINAL_WARNING_VIEW
	}

	if p.CurrentView == views.TERMINAL_WARNING_VIEW && p.Terminal.IsValid() {
		p.CurrentView = views.PLAYER_VIEW
	}

	switch msg := msg.(type) {
	case tickMsg:
		// If state is unaccessible, likely due to user closing
		// their playerback device, and attempt reconnect to closed device.
		if p.PlayerState() == nil {
			p.Player.StatusBar.Update(p.PlayerState())

			if p.Player.State != nil && p.Player.State.IsPlaying {
				p.player.Resume(p.session, false)
			} else {
				p.player.Resume(p.session, true)
			}

			return p, tea.Tick(4*UPDATE_RATE_SEC, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		p.Player.EnsureProgressSynced()

		return p, tea.Tick(UPDATE_RATE_SEC, func(time.Time) tea.Msg {
			return tickMsg{}
		})

	case tea.KeyMsg:
		// Prevents search query from activating any commands, enless esc or enter.
		key := msg.String()

		if p.CurrentView != views.SEARCH_VIEW_QUERY &&
			(key == KEY_SEARCH_VIEW || key == KEY_SEARCH_VIEW_ALT) {
			p.Search.Input.Text.Focus()
			p.CurrentView = views.SEARCH_VIEW_QUERY
			return p, nil
		}

		if p.CurrentView == views.SEARCH_VIEW_QUERY && !IsImportantKey(key) {
			var cmd tea.Cmd
			p.Search.Input, cmd = p.Search.Input.Update(msg)
			return p, cmd
		}

		if p.CurrentView == views.SEARCH_VIEW_TYPE && !IsImportantKey(key) {
			var cmd tea.Cmd
			p.Search.TypeList, cmd = p.Search.TypeList.Update(msg)
			return p, cmd
		}

		if (key == KEY_SEARCH_VIEW || key == KEY_SEARCH_VIEW_ALT) &&
			p.CurrentView != views.SEARCH_VIEW_QUERY {
			p.CurrentView = views.SEARCH_VIEW_QUERY
			return p, nil
		}

		switch msg.String() {
		case KEY_ESC:
			switch p.CurrentView {
			case views.SEARCH_VIEW_QUERY:
				p.CurrentView = views.PLAYER_VIEW
			default:
			}

		case KEY_QUIT, KEY_QUIT_ALT:
			return p, tea.Quit

		case KEY_PLAY_PAUSE:
			err := p.Player.PlayPause()
			if errors.IsReauthenticationErr(err) {
				p.CurrentView = views.REAUTH_VIEW
			}

		case KEY_FZF_DEVICES:
			p.CurrentView = views.DEVICE_FZF_VIEW

		case KEY_VOLUME_DOWN_BIG:
			// Spotify doesn't have a volume control for mobile devices.
			if p.player.Device() != nil && !p.player.Device().IsMobile() {
				vol := p.PlayerState().Device.VolumePercent
				newVol := vol - VOLUME_INCREMENT_PERCENT

				if 0 < vol && vol <= 5 {
					newVol = 0
				}

				if !player.IsValidVolume(newVol) {
					break
				}

				err := p.Player.Player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.CurrentView = views.REAUTH_VIEW
				}
				p.Player.State.Device.VolumePercent = newVol
			}

		case KEY_VOLUME_UP_BIG:
			if p.player.Device() != nil && !p.player.Device().IsMobile() {
				vol := p.PlayerState().Device.VolumePercent
				newVol := vol + VOLUME_INCREMENT_PERCENT

				if 95 <= vol && vol < 100 {
					newVol = 100
				}

				if !player.IsValidVolume(newVol) {
					break
				}

				err := p.Player.Player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.CurrentView = views.REAUTH_VIEW
				}
				p.Player.State.Device.VolumePercent = newVol
			}

		case KEY_VOLUME_DOWN_SMALL:
			if p.player.Device() != nil && !p.player.Device().IsMobile() {
				vol := p.PlayerState().Device.VolumePercent
				newVol := vol - 1

				if 0 < vol && vol <= 5 {
					newVol = 0
				}

				if !player.IsValidVolume(newVol) {
					break
				}

				err := p.Player.Player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.CurrentView = views.REAUTH_VIEW
				}
				p.Player.State.Device.VolumePercent = newVol
			}

		case KEY_VOLUME_UP_SMALL:
			if p.player.Device() != nil && !p.player.Device().IsMobile() {
				vol := p.PlayerState().Device.VolumePercent
				newVol := vol + 1

				if 95 <= vol && vol < 100 {
					newVol = 100
				}

				if !player.IsValidVolume(newVol) {
					break
				}

				err := p.Player.Player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.CurrentView = views.REAUTH_VIEW
				}
				p.Player.State.Device.VolumePercent = newVol
			}

		case KEY_PLAYER_VIEW, KEY_PLAYER_VIEW_ALT:
			p.CurrentView = views.PLAYER_VIEW

		case KEY_PLAYLIST_VIEW, KEY_PLAYLIST_VIEW_ALT:
			p.CurrentView = views.PLAYLIST_VIEW

		case KEY_SEARCH_VIEW, KEY_SEARCH_VIEW_ALT:
			// Requires handling priority, logic is at the top.

		// case KEY_DEVICE_VIEW:
		// 	p.Device.UpdateNumberDevices()
		// 	p.CurrentView = views.DEVICE_VIEW

		case KEY_HELP_VIEW, KEY_HELP_VIEW_ALT:
			p.CurrentView = views.HELP_VIEW

		case KEY_ENTER:

			// The enter key has different actions it needs to perform depending on the
			// current view.
			switch p.CurrentView {
			case views.PLAYLIST_VIEW:
				pl := p.Playlist.GetSelectedPlaylist()
				p.player.Play(pl.Uri, "", p.session)

				p.Player.UpdateStateSync()

			case views.SEARCH_VIEW_QUERY:
				p.Search.Input = p.Search.Input.HideCursor()
				p.CurrentView = views.SEARCH_VIEW_TYPE

			case views.SEARCH_VIEW_TYPE:
				p.Search.Results = p.Search.Results.Refresh(p.Search.Input.Query(), p.Search.SelectedType(), p.session)
				p.CurrentView = views.SEARCH_VIEW_RESULTS

			case views.SEARCH_VIEW_RESULTS:
				switch p.Search.SelectedType() {
				case views.TRACK:
					err := p.player.Play(p.Search.Results.SelectedTrack().Album.Uri, p.Search.Results.SelectedTrack().Uri, p.session)
					if errors.IsReauthenticationErr(err) {
						p.CurrentView = views.REAUTH_VIEW
					}
					p.Player.UpdateStateSync()
				case views.ALBUM:
					err := p.player.Play(p.Search.Results.SelectedAlbum().Uri, "", p.session)
					if errors.IsReauthenticationErr(err) {
						p.CurrentView = views.REAUTH_VIEW
					}
					p.Player.UpdateStateSync()

				case views.PLAYLIST:
					err := p.player.Play(p.Search.Results.SelectedPlaylist().Uri, "", p.session)
					if errors.IsReauthenticationErr(err) {
						p.CurrentView = views.REAUTH_VIEW
					}
					p.Player.UpdateStateSync()
				}
			}

		case KEY_VISUAL_REFRESH:
			// Refreshes the terminal fixing any visual glitches. This doesn't yet force any
			// updates to, for example, listed playlist devices.
			go func() {
				view := p.CurrentView

				p.CurrentView = views.REFRESH_VIEW
				time.Sleep(UPDATE_RATE_SEC)
				p.CurrentView = view

				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()

		case KEY_FZF_PLAYLIST_TRACKS:
			if p.CurrentView == views.PLAYLIST_VIEW {
				p.CurrentView = views.PLAYLIST_TRACK_VIEW
			}

		case KEY_FZF_ALBUM_TRACKS:
			p.CurrentView = views.ALBUM_TRACK_VIEW

		case KEY_TOGGLE_SHUFFLING:
			// Enables or disables shuffling on current album or playlist.
			state := p.PlayerState().ShuffleState

			p.PlayerState().ShuffleState = !state

			p.player.Shuffle(!state, p.session)

		case KEY_TOGGLE_REPEAT:
			switch p.PlayerState().RepeatState {
			case "off":
				err := p.player.Repeat(true, p.session)
				if errors.IsReauthenticationErr(err) {
					p.CurrentView = views.REAUTH_VIEW
				}

				p.PlayerState().RepeatState = "context"
			default:
				err := p.player.Repeat(false, p.session)
				if errors.IsReauthenticationErr(err) {
					p.CurrentView = views.REAUTH_VIEW
				}

				p.PlayerState().RepeatState = "off"
			}
		}

		var cmd tea.Cmd

		// Handles updates from the playlist list.
		if p.CurrentView == views.PLAYLIST_VIEW {
			p.Playlist.PlaylistList, cmd = p.Playlist.PlaylistList.Update(msg)
			return p, cmd
		}

		if p.CurrentView == views.SEARCH_VIEW_QUERY {
			p.Search.Input, cmd = p.Search.Input.Update(msg)
			return p, cmd
		}

		if p.CurrentView == views.SEARCH_VIEW_RESULTS {
			p.Search.Results, cmd = p.Search.Results.Update(msg)
			return p, cmd
		}
	}

	return p, nil
}

// Returns the player state from the model's player view.
func (p *Program) PlayerState() *player.State {
	return p.Player.State
}

// Returns is a keys that should always go through,
// for exmaple keys to change the view, quit keys
// are excluded since search query uses q.
func IsImportantKey(key string) bool {
	keys := []string{
		KEY_ENTER,
		KEY_PLAYER_VIEW, KEY_PLAYER_VIEW_ALT,
		KEY_PLAYLIST_VIEW, KEY_PLAYER_VIEW_ALT,
		KEY_SEARCH_VIEW, KEY_SEARCH_VIEW_ALT,
		KEY_DEVICE_VIEW,
		KEY_HELP_VIEW, KEY_HELP_VIEW_ALT,
		KEY_FZF_DEVICES,
		KEY_FZF_ALBUM_TRACKS,
	}

	for _, k := range keys {
		if key == k {
			return true
		}
	}

	return false
}
