package tui

import (
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/views"
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

			p.player.Resume(p.session, false)

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

		if p.CurrentView == views.SEARCH_VIEW_RESULTS && key == "/" {
			p.Search.Input.Text.Focus()
			p.CurrentView = views.SEARCH_VIEW_QUERY
			return p, nil
		}

		if p.CurrentView == views.SEARCH_VIEW_QUERY &&
			key != "enter" && key != "esc" &&
			key != "f1" && key != "f2" && key != "f4" && key != "f5" {

			var cmd tea.Cmd
			p.Search.Input, cmd = p.Search.Input.Update(msg)
			return p, cmd

		}

		if p.CurrentView == views.SEARCH_VIEW_TYPE &&
			key != "enter" && key != "esc" &&
			key != "f1" && key != "f2" && key != "f4" && key != "f5" {

			var cmd tea.Cmd
			p.Search.TypeList, cmd = p.Search.TypeList.Update(msg)
			return p, cmd
		}

		if (key == "/" || key == "3") && p.CurrentView != views.SEARCH_VIEW_QUERY {
			p.CurrentView = views.SEARCH_VIEW_QUERY
			return p, nil
		}

		switch msg.String() {
		case "esc":
			switch p.CurrentView {
			case views.SEARCH_VIEW_QUERY:
				p.CurrentView = views.PLAYER_VIEW
			default:
			}

		case "ctrl+c", "q":
			return p, tea.Quit

		case " ":
			p.Player.PlayPause()

		case "[":
			t := p.player.Device().Type
			// Spotify doesn't have a volume control for mobile devices.
			if p.player.Device() != nil && (t != "smartphone" || t != "tablet") {
				vol := p.PlayerState().Device.VolumePercent
				p.Player.Player.SetVolume(p.session, vol-VOLUME_INCREMENT_PERCENT)
				p.Player.State.Device.VolumePercent = vol - VOLUME_INCREMENT_PERCENT
			}

		case "]":
			t := p.player.Device().Type
			// Spotify doesn't have a volume control for mobile devices.
			if p.player.Device() != nil && (t != "smartphone" || t != "tablet") {
				vol := p.PlayerState().Device.VolumePercent
				p.Player.Player.SetVolume(p.session, vol+VOLUME_INCREMENT_PERCENT)
				p.Player.State.Device.VolumePercent = vol + VOLUME_INCREMENT_PERCENT
			}

		case "{":
			t := p.player.Device().Type
			// Spotify doesn't have a volume control for mobile devices.
			if p.player.Device() != nil && (t != "smartphone" || t != "tablet") {
				vol := p.PlayerState().Device.VolumePercent
				p.Player.Player.SetVolume(p.session, vol-1)
				p.Player.State.Device.VolumePercent = vol - 1
			}

		case "}":
			t := p.player.Device().Type
			// Spotify doesn't have a volume control for mobile devices.
			if p.player.Device() != nil && (t != "smartphone" || t != "tablet") {
				vol := p.PlayerState().Device.VolumePercent
				p.Player.Player.SetVolume(p.session, vol+1)
				p.Player.State.Device.VolumePercent = vol + 1
			}

		case "f1", "1":
			p.CurrentView = views.PLAYER_VIEW

		case "f2", "2":
			p.CurrentView = views.PLAYLIST_VIEW

		case "f3", "3", "/":
			p.Search.Input.Text.Focus()
			p.CurrentView = views.SEARCH_VIEW_QUERY

		case "f4", "4":
			p.Device.UpdateDevices()
			p.CurrentView = views.DEVICE_VIEW

		case "f5", "5":
			p.CurrentView = views.HELP_VIEW

		case "enter":

			// The enter key has different actions it needs to perform depending on the
			// current view.
			switch p.CurrentView {
			case views.PLAYLIST_VIEW:
				pl := p.Playlist.GetSelectedPlaylist()
				p.player.Play(pl.URI, "", p.session)

				p.Player.UpdateStateSync()

			case views.SEARCH_VIEW_QUERY:
				p.Search.Input = p.Search.Input.HideCursor()
				p.CurrentView = views.SEARCH_VIEW_TYPE

			case views.DEVICE_VIEW:
				device := p.Device.GetSelectedDevice()

				p.player.SetDevice(device, p.Config)

				// Transfers playback to the newly select device.
				p.player.Resume(p.session, false)

			case views.SEARCH_VIEW_TYPE:
				p.Search.Results = p.Search.Results.Refresh(p.Search.Input.Query(), p.Search.SelectedType(), p.session)
				p.CurrentView = views.SEARCH_VIEW_RESULTS

			case views.SEARCH_VIEW_RESULTS:
				switch p.Search.SelectedType() {
				case "track":
					p.player.Play(p.Search.Results.SelectedTrack().Album.Uri, p.Search.Results.SelectedTrack().Uri, p.session)
					p.Player.UpdateStateSync()
				case "album":
					p.player.Play(p.Search.Results.SelectedAlbum().Uri, "", p.session)
					p.Player.UpdateStateSync()
				}
			}

		case "ctrl+r":
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

		case "t":
			if p.CurrentView == views.PLAYLIST_VIEW {
				p.CurrentView = views.PLAYLIST_TRACK_VIEW
			}

		case "a":
			p.CurrentView = views.ALBUM_TRACK_VIEW

		case "s":
			// Enables or disables shuffling on current album or playlist.
			state := p.PlayerState().ShuffleState

			p.PlayerState().ShuffleState = !state

			p.player.Shuffle(!state, p.session)

		case "r":
			switch p.PlayerState().RepeatState {
			case "off":
				p.player.Repeat(true, p.session)
				p.PlayerState().RepeatState = "context"
			default:
				p.player.Repeat(false, p.session)
				p.PlayerState().RepeatState = "off"
			}
		}

		var cmd tea.Cmd

		// Handles updates from the playlist list.
		if p.CurrentView == views.PLAYLIST_VIEW {
			p.Playlist.PlaylistList, cmd = p.Playlist.PlaylistList.Update(msg)
			return p, cmd
		}

		// Handles updates from the device list.
		if p.CurrentView == views.DEVICE_VIEW {
			p.Device.ListModel, cmd = p.Device.ListModel.Update(msg)
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
