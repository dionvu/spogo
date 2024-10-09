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
		p.CurrentView = TERMINAL_WARNING_VIEW
	}

	if p.CurrentView == TERMINAL_WARNING_VIEW && p.Terminal.IsValid() {
		p.CurrentView = PLAYER_VIEW
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

		switch msg.String() {
		case "esc":
			switch p.CurrentView {
			case views.SEARCH_VIEW_QUERY:
				p.CurrentView = PLAYER_VIEW
			default:
			}

		case "ctrl+c", "q":
			return p, tea.Quit

		case " ":
			p.Player.PlayPause()

		case "[":
			vol := p.PlayerState().Device.VolumePercent
			p.Player.Player.SetVolume(p.session, vol-VOLUME_INCREMENT_PERCENT)

		case "]":
			vol := p.PlayerState().Device.VolumePercent
			p.Player.Player.SetVolume(p.session, vol+VOLUME_INCREMENT_PERCENT)

		case "f1", "1":
			p.CurrentView = PLAYER_VIEW

		case "f2", "2":
			p.CurrentView = PLAYLIST_VIEW

		case "f3", "3":
			p.Search.Input.Text.Focus()
			p.CurrentView = views.SEARCH_VIEW_QUERY
			// p.Views.Squery.textInput.SetValue("") // Resets the search value.

		case "f4", "4":
			// p.Views.Device = NewDeviceView(p.session) // Updates the list of available devices.
			p.Device.UpdateDevices()
			p.CurrentView = DEVICE_VIEW

		case "f5", "5":
			p.CurrentView = HELP_VIEW

		case "enter":

			// The enter key has different actions it needs to perform depending on the
			// current view.
			switch p.CurrentView {
			case PLAYLIST_VIEW:
				// if i, ok := p..Playlist.PlaylistList.list.SelectedItem().(views.Item); ok {
				// 	p.Views.Playlist.PlaylistList.choice = string(i)
				// 	err := p.player.Play(p.Views.Playlist.GetSelectedPlaylist().URI, "", p.session)
				// 	if err != nil {
				// 		log.Fatal(string(i), p.Views.Playlist.GetSelectedPlaylist().URI)
				// 	}
				// }

			case views.SEARCH_VIEW_QUERY:
				p.Search.Input = p.Search.Input.HideCursor()
				p.CurrentView = views.SEARCH_VIEW_TYPE

			case DEVICE_VIEW:
				device := p.Device.GetSelectedDevice()

				p.player.SetDevice(device, p.Config)

				// Transfers playback to the newly select device.
				p.player.Resume(p.session, false)

				// case SEARCH_QUERY_VIEW:
				// if p.Views.Squery.Query() != "" {
				// 	p.CurrentView = SEARCH_TYPE_VIEW
				// }

			case views.SEARCH_VIEW_TYPE:
				p.Search.Results = p.Search.Results.Refresh(p.Search.Input.Query(), p.session)
				p.CurrentView = views.SEARCH_VIEW_RESULTS
			}

		case "ctrl+r":
			// Refreshes the terminal fixing any visual glitches. This doesn't yet force any
			// updates to, for example, listed playlist devices.
			go func() {
				view := p.CurrentView

				p.CurrentView = REFRESH_VIEW
				time.Sleep(UPDATE_RATE_SEC)
				p.CurrentView = view

				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()

		case "t":
			if p.CurrentView == PLAYLIST_VIEW {
				p.CurrentView = PLAYLIST_TRACK_VIEW
			}

		case "a":
			p.CurrentView = ALBUM_TRACK_VIEW

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
		if p.CurrentView == PLAYLIST_VIEW {
			p.Playlist.PlaylistList, cmd = p.Playlist.PlaylistList.Update(msg)
			return p, cmd
		}

		// Handles updates from the device list.
		// if p.CurrentView == DEVICE_VIEW {
		// 	p.Views.Device.ListModel.list, cmd = p.Views.Device.ListModel.list.Update(msg)
		// 	return p, cmd
		// }

		if p.CurrentView == views.SEARCH_VIEW_QUERY {
			p.Search.Input, cmd = p.Search.Input.Update(msg)
			return p, cmd
		}

		if p.CurrentView == views.SEARCH_VIEW_RESULTS {
			p.Search.Results, cmd = p.Search.Results.Update(msg)
			return p, cmd
		}

		// if p.CurrentView == SEARCH_TYPE_VIEW {
		// 	p.Views.SearchType.ListModel, cmd = p.Views.SearchType.ListModel.Update(msg)
		// 	return p, cmd
		// }
	}

	return p, nil
}

// Returns the player state from the model's player view.
func (p *Program) PlayerState() *player.State {
	return p.Player.State
}
