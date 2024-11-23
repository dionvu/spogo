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
	KEY_ESC                  = "esc"
	KEY_ENTER                = "enter"
	KEY_QUIT                 = "q"
	KEY_QUIT_ALT             = "ctrl+c"
	KEY_VISUAL_REFRESH       = "ctrl+r"
	KEY_PLAYER_VIEW          = "f1"
	KEY_PLAYER_VIEW_ALT      = "ctrl+o"
	KEY_PLAY_PAUSE           = " "
	KEY_TOGGLE_SHUFFLING     = "s"
	KEY_TOGGLE_REPEAT        = "r"
	KEY_PLAYLIST_VIEW        = "f2"
	KEY_PLAYLIST_VIEW_ALT    = "ctrl+p"
	KEY_FZF_PLAYLIST_TRACKS  = "t"
	KEY_SEARCH_VIEW          = "f3"
	KEY_SEARCH_VIEW_ALT      = "/"
	KEY_DEVICE_VIEW          = "UNDEFINED"
	KEY_HELP_VIEW            = "f4"
	KEY_HELP_VIEW_ALT        = "ctrl+h"
	KEY_VOLUME_DOWN_BIG      = "["
	KEY_VOLUME_DOWN_SMALL    = "{"
	KEY_VOLUME_UP_BIG        = "]"
	KEY_VOLUME_UP_SMALL      = "}"
	KEY_FZF_DEVICES          = "ctrl+d"
	KEY_FZF_ALBUM_TRACKS     = "ctrl+a"
	KEY_NEXT_TRACK           = ">"
	KEY_PREV_TRACK           = "<"
	KEY_FORWARD              = "."
	KEY_BACKWARD             = ","
	VOLUME_INCREMENT_PERCENT = 5
	EMPTY                    = ""

	ENABLED  = "on"
	DISABLED = "off"
)

// Handles updates associate with the current selected view.
func (p *Program) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p.terminal.UpdateSize()

	if !p.terminal.IsValid() {
		p.currentView = views.TERMINAL_WARNING_VIEW
	}

	if p.currentView == views.TERMINAL_WARNING_VIEW && p.terminal.IsValid() {
		p.currentView = views.PLAYER_VIEW
	}

	var cmd tea.Cmd
	if p.currentView == views.HELP_VIEW {
		p.help, cmd = p.help.Update(msg)
		return p, cmd
	}

	switch msg := msg.(type) {
	case tickMsg:
		// If state is unaccessible, likely due to user closing
		// their playerback device, and attempt reconnect to closed device.
		if p.PlayerState() == nil {
			p.playerView.UpdateStatusBar(p.PlayerState())

			if p.playerView.State != nil && p.playerView.State.IsPlaying {
				p.player.Resume(p.session, false)
			} else {
				p.player.Resume(p.session, true)
			}

			return p, tea.Tick(4*UPDATE_RATE_SEC, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		p.playerView.EnsureProgressSynced()

		return p, tea.Tick(UPDATE_RATE_SEC, func(time.Time) tea.Msg {
			return tickMsg{}
		})

	case tea.KeyMsg:
		// Prevents search query from activating any commands, enless esc or enter.
		key := msg.String()

		if p.currentView != views.SEARCH_VIEW_QUERY &&
			(key == KEY_SEARCH_VIEW || key == KEY_SEARCH_VIEW_ALT) {
			p.search.Input.Text.Focus()
			p.currentView = views.SEARCH_VIEW_QUERY
			return p, nil
		}

		if p.currentView == views.SEARCH_VIEW_QUERY && !IsImportantKey(key) {
			var cmd tea.Cmd
			p.search.Input, cmd = p.search.Input.Update(msg)
			return p, cmd
		}

		if p.currentView == views.SEARCH_VIEW_TYPE && !IsImportantKey(key) {
			var cmd tea.Cmd
			p.search.TypeList, cmd = p.search.TypeList.Update(msg)
			return p, cmd
		}

		if (key == KEY_SEARCH_VIEW || key == KEY_SEARCH_VIEW_ALT) &&
			p.currentView != views.SEARCH_VIEW_QUERY {
			p.currentView = views.SEARCH_VIEW_QUERY
			return p, nil
		}

		switch msg.String() {
		case KEY_ESC:
			switch p.currentView {
			case views.SEARCH_VIEW_QUERY:
				p.currentView = views.PLAYER_VIEW
			default:
			}

		case KEY_QUIT, KEY_QUIT_ALT:
			return p, tea.Quit

		case KEY_PLAY_PAUSE:
			err := p.playerView.PlayPause()
			if errors.IsReauthenticationErr(err) {
				p.currentView = views.REAUTH_VIEW
			}

		case KEY_FZF_DEVICES:
			p.currentView = views.DEVICE_FZF_VIEW

		case KEY_PREV_TRACK:
			if p.currentView == views.PLAYER_VIEW {
				p.player.SkipPrev(p.session)
			}

			const STATE_DELAY_INTERVAL = time.Second / 100

			time.Sleep(STATE_DELAY_INTERVAL)

			p.playerView.UpdateStateSync()

		case KEY_NEXT_TRACK:
			if p.currentView == views.PLAYER_VIEW {
				p.player.SkipNext(p.session)
			}

			time.Sleep(time.Second / 100)

			p.playerView.UpdateStateSync()

		case KEY_FORWARD:
			if p.currentView == views.PLAYER_VIEW && p.playerView.State != nil &&
				p.playerView.State.Track != nil {
				pos := p.playerView.State.ProgressMs + 10000
				if pos > p.playerView.State.Track.DurationMs {
					pos = p.playerView.State.Track.DurationMs
				} else if pos < 0 {
					pos = 0
				}

				p.player.Seek(pos, p.session)
			}

			time.Sleep(time.Second / 100)

			p.playerView.UpdateStateSync()

		case KEY_BACKWARD:

			if p.currentView == views.PLAYER_VIEW && p.playerView.State != nil &&
				p.playerView.State.Track != nil {
				pos := p.playerView.State.ProgressMs - 10000
				if pos > p.playerView.State.Track.DurationMs {
					pos = p.playerView.State.Track.DurationMs
				} else if pos < 0 {
					pos = 0
				}

				p.player.Seek(pos, p.session)
			}

			time.Sleep(time.Second / 100)

			p.playerView.UpdateStateSync()

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

				err := p.player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.currentView = views.REAUTH_VIEW
				}
				p.playerView.State.Device.VolumePercent = newVol
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

				err := p.player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.currentView = views.REAUTH_VIEW
				}
				p.playerView.State.Device.VolumePercent = newVol
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

				err := p.player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.currentView = views.REAUTH_VIEW
				}
				p.playerView.State.Device.VolumePercent = newVol
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

				err := p.player.SetVolume(p.session, newVol)
				if errors.IsReauthenticationErr(err) {
					p.currentView = views.REAUTH_VIEW
				}
				p.playerView.State.Device.VolumePercent = newVol
			}

		case KEY_PLAYER_VIEW, KEY_PLAYER_VIEW_ALT:
			p.currentView = views.PLAYER_VIEW

		case KEY_PLAYLIST_VIEW, KEY_PLAYLIST_VIEW_ALT:
			p.currentView = views.PLAYLIST_VIEW

		case KEY_SEARCH_VIEW, KEY_SEARCH_VIEW_ALT:
			// Requires handling priority, logic is at the top.

		// case KEY_HELP_VIEW, KEY_HELP_VIEW_ALT:
		// 	p.currentView = views.HELP_VIEW

		case KEY_ENTER:
			switch p.currentView {
			case views.PLAYLIST_VIEW:
				pl := p.playlistView.GetSelectedPlaylist()
				p.player.Play(pl.Uri, "", p.session)

				p.playerView.UpdateStateSync()

			case views.SEARCH_VIEW_QUERY:
				if p.search.Input.Text.Value() != EMPTY {
					p.search.Input = p.search.Input.HideCursor()
					p.currentView = views.SEARCH_VIEW_TYPE
				}

			case views.SEARCH_VIEW_TYPE:
				p.search.Results = p.search.Results.Refresh(p.search.Input.Query(), p.search.SelectedType(), p.session)
				p.currentView = views.SEARCH_VIEW_RESULTS

			case views.SEARCH_VIEW_RESULTS:

				switch p.search.SelectedType() {
				case views.TRACK:
					if p.search.Results.SelectedTrack() == nil {
						return p, nil
					}

					err := p.player.Play(p.search.Results.SelectedTrack().Album.Uri, p.search.Results.SelectedTrack().Uri, p.session)
					if errors.IsReauthenticationErr(err) {
						p.currentView = views.REAUTH_VIEW
					}

					p.playerView.UpdateStateSync()

				case views.ALBUM:
					if p.search.Results.SelectedAlbum() == nil {
						return p, nil
					}

					err := p.player.Play(p.search.Results.SelectedAlbum().Uri, EMPTY, p.session)
					if errors.IsReauthenticationErr(err) {
						p.currentView = views.REAUTH_VIEW
					}

					p.playerView.UpdateStateSync()

				case views.PLAYLIST:
					if p.search.Results.SelectedPlaylist() == nil {
						return p, nil
					}

					err := p.player.Play(p.search.Results.SelectedPlaylist().Uri, EMPTY, p.session)
					if errors.IsReauthenticationErr(err) {
						p.currentView = views.REAUTH_VIEW
					}

					p.playerView.UpdateStateSync()

				default:
					return p, nil
				}

			}

		case KEY_VISUAL_REFRESH:
			// Refreshes the terminal fixing any visual glitches. This doesn't yet force any
			// updates to, for example, listed playlist devices.
			go func() {
				view := p.currentView

				p.currentView = views.REFRESH_VIEW
				time.Sleep(UPDATE_RATE_SEC)
				p.currentView = view

				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()

		case KEY_FZF_PLAYLIST_TRACKS:
			if p.currentView == views.PLAYLIST_VIEW {
				p.currentView = views.PLAYLIST_TRACK_VIEW
			}

		case KEY_FZF_ALBUM_TRACKS:
			p.currentView = views.ALBUM_TRACK_VIEW

		case KEY_TOGGLE_SHUFFLING:
			// Enables or disables shuffling on current album or playlist.
			state := p.PlayerState().ShuffleState

			p.PlayerState().ShuffleState = !state

			p.player.Shuffle(!state, p.session)

		case KEY_TOGGLE_REPEAT:
			switch p.PlayerState().RepeatState {
			case DISABLED:
				err := p.player.Repeat(true, p.session)
				if errors.IsReauthenticationErr(err) {
					p.currentView = views.REAUTH_VIEW
				}

				p.PlayerState().RepeatState = "context"
			default:
				err := p.player.Repeat(false, p.session)
				if errors.IsReauthenticationErr(err) {
					p.currentView = views.REAUTH_VIEW
				}

				p.PlayerState().RepeatState = DISABLED
			}
		}

		var cmd tea.Cmd

		// Handles updates from the playlist list.
		if p.currentView == views.PLAYLIST_VIEW {
			p.playlistView.PlaylistList, cmd = p.playlistView.PlaylistList.Update(msg)
			return p, cmd
		}

		if p.currentView == views.SEARCH_VIEW_QUERY {
			p.search.Input, cmd = p.search.Input.Update(msg)
			return p, cmd
		}

		if p.currentView == views.SEARCH_VIEW_RESULTS {
			p.search.Results, cmd = p.search.Results.Update(msg)
			return p, cmd
		}
	}

	return p, nil
}

// Returns the player state from the model's player view.
func (p *Program) PlayerState() *player.State {
	return p.playerView.State
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
