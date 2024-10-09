package tui

import (
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/player"
)

// Handles updates associate with the current selected view.
func (m *Program) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Terminal.UpdateSize()

	if !m.Terminal.IsValid() {
		m.CurrentView = TERMINAL_WARNING_VIEW
	}

	if m.CurrentView == TERMINAL_WARNING_VIEW && m.Terminal.IsValid() {
		m.CurrentView = PLAYER_VIEW
	}

	switch msg := msg.(type) {
	case tickMsg:
		// If state is unaccessible, likely due to user closing
		// their playerback device, and attempt reconnect to closed device.
		if m.PlayerState() == nil {
			m.Views.Player.StatusBar.Update(m.PlayerState())

			m.player.Resume(m.session, false)

			return m, tea.Tick(4*UPDATE_RATE_SEC, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		m.Views.Player.EnsureProgressSynced()

		return m, tea.Tick(UPDATE_RATE_SEC, func(time.Time) tea.Msg {
			return tickMsg{}
		})

	case tea.KeyMsg:
		// Prevents search query from activating any commands, enless esc or enter.
		// key := msg.String()
		// if m.CurrentView == SEARCH_QUERY_VIEW &&
		// key != "enter" &&
		// key != "esc" &&
		// key != "f1" &&
		// key != "f2" &&
		// key != "f4" &&
		// key != "f5" {
		// var cmd tea.Cmd
		// m.Views.Squery.textInput, cmd = m.Views.Squery.textInput.Update(msg)
		// return m, cmd
		// }

		switch msg.String() {
		case "esc":
			if m.CurrentView == SEARCH_QUERY_VIEW {
				m.CurrentView = PLAYER_VIEW
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		case " ":
			m.Views.Player.PlayPause()

		case "[":
			vol := m.PlayerState().Device.VolumePercent
			m.Views.Player.Player.SetVolume(m.session, vol-VOLUME_INCREMENT_PERCENT)

		case "]":
			vol := m.PlayerState().Device.VolumePercent
			m.Views.Player.Player.SetVolume(m.session, vol+VOLUME_INCREMENT_PERCENT)

		case "f1", "1":
			m.CurrentView = PLAYER_VIEW

		case "f2", "2":
			m.CurrentView = PLAYLIST_VIEW

		case "f3", "3":
			m.CurrentView = SEARCH_VIEW
			// m.Views.Squery.textInput.SetValue("") // Resets the search value.

		case "f4", "4":
			// m.Views.Device = NewDeviceView(m.session) // Updates the list of available devices.
			m.Views.Device.UpdateDevices()
			m.CurrentView = DEVICE_VIEW

		case "f5", "5":
			m.CurrentView = HELP_VIEW

		case "enter":

			// The enter key has different actions it needs to perform depending on the
			// current view.
			switch m.CurrentView {
			case PLAYLIST_VIEW:
				// if i, ok := m.Views.Playlist.PlaylistList.list.SelectedItem().(views.Item); ok {
				// 	m.Views.Playlist.PlaylistList.choice = string(i)
				// 	err := m.player.Play(m.Views.Playlist.GetSelectedPlaylist().URI, "", m.session)
				// 	if err != nil {
				// 		log.Fatal(string(i), m.Views.Playlist.GetSelectedPlaylist().URI)
				// 	}
				// }

			case SEARCH_TYPE_VIEW:
				m.CurrentView = SEARCH_RESULT_TRACK

			case DEVICE_VIEW:
				device := m.Views.Device.GetSelectedDevice()

				m.player.SetDevice(device, m.Config)

				// Transfers playback to the newly select device.
				m.player.Resume(m.session, false)

				// case SEARCH_QUERY_VIEW:
				// if m.Views.Squery.Query() != "" {
				// 	m.CurrentView = SEARCH_TYPE_VIEW
				// }
			}

		case "ctrl+r":
			// Refreshes the terminal fixing any visual glitches. This doesn't yet force any
			// updates to, for example, listed playlist devices.
			go func() {
				view := m.CurrentView

				m.CurrentView = REFRESH_VIEW
				time.Sleep(UPDATE_RATE_SEC)
				m.CurrentView = view

				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()

		case "t":
			if m.CurrentView == PLAYLIST_VIEW {
				m.CurrentView = PLAYLIST_TRACK_VIEW
			}

		case "a":
			m.CurrentView = ALBUM_TRACK_VIEW

		case "s":
			// Enables or disables shuffling on current album or playlist.
			state := m.PlayerState().ShuffleState

			m.PlayerState().ShuffleState = !state

			m.player.Shuffle(!state, m.session)

		case "r":
			switch m.PlayerState().RepeatState {
			case "off":
				m.player.Repeat(true, m.session)
				m.PlayerState().RepeatState = "context"
			default:
				m.player.Repeat(false, m.session)
				m.PlayerState().RepeatState = "off"
			}
		}

		var cmd tea.Cmd

		// Handles updates from the playlist list.
		if m.CurrentView == PLAYLIST_VIEW {
			m.Views.Playlist.PlaylistList, cmd = m.Views.Playlist.PlaylistList.Update(msg)
			return m, cmd
		}

		// Handles updates from the device list.
		// if m.CurrentView == DEVICE_VIEW {
		// 	m.Views.Device.ListModel.list, cmd = m.Views.Device.ListModel.list.Update(msg)
		// 	return m, cmd
		// }

		// if m.CurrentView == SEARCH_TYPE_VIEW {
		// 	m.Views.SearchType.ListModel, cmd = m.Views.SearchType.ListModel.Update(msg)
		// 	return m, cmd
		// }
	}

	return m, nil
}

// Returns the player state from the model's player view.
func (m *Program) PlayerState() *player.State {
	return m.Views.Player.State
}
