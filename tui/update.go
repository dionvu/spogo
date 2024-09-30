package ui

import (
	"log"
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

			m.Player.Resume(m.Session, false)

			return m, tea.Tick(4*UPDATE_RATE_SEC, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		m.Views.Player.EnsureProgressSynced()

		return m, tea.Tick(UPDATE_RATE_SEC, func(time.Time) tea.Msg {
			return tickMsg{}
		})

	case tea.KeyMsg:

		// Prevents search query from activating any commands.
		if m.CurrentView == SEARCH_QUERY_VIEW && msg.String() != "enter" {
			var cmd tea.Cmd
			m.Views.Squery.textInput, cmd = m.Views.Squery.textInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case " ":
			// m.Views.Player.UpdateStateSync()
			m.Views.Player.PlayPause()

		case "[":
			vol := m.PlayerState().Device.VolumePercent
			m.Views.Player.Player.SetVolume(m.Session, vol-VOLUME_INCREMENT_PERCENT)

		case "]":
			vol := m.PlayerState().Device.VolumePercent
			m.Views.Player.Player.SetVolume(m.Session, vol+VOLUME_INCREMENT_PERCENT)

		case "f1", "1":
			m.CurrentView = PLAYER_VIEW

		case "f2", "2":
			m.CurrentView = PLAYLIST_VIEW

		case "f3", "3":
			m.CurrentView = SEARCH_TYPE_VIEW

		case "f4", "4":
			m.Views.Device = NewDeviceView(m.Session) // Updates the list of available devices.
			m.Views.Device.UpdateDevices()
			m.CurrentView = DEVICE_VIEW

		case "f5", "5":
			m.CurrentView = HELP_VIEW

		case "enter":

			// The enter key has different actions it needs to perform depending on the
			// current view.
			switch m.CurrentView {

			case PLAYLIST_VIEW:
				if i, ok := m.Views.Playlist.PlaylistList.list.SelectedItem().(Item); ok {
					m.Views.Playlist.PlaylistList.choice = string(i)
					err := m.Player.Play(m.Views.Playlist.playlistsMap[string(i)].URI, "", m.Session)
					if err != nil {
						log.Fatal(string(i), m.Views.Playlist.playlistsMap[string(i)].URI)
					}
				}

			case SEARCH_TYPE_VIEW:
				if i, ok := m.Views.SearchType.ListModel.list.SelectedItem().(Item); ok {
					m.Views.SearchType.ListModel.choice = string(i)

					m.CurrentView = SEARCH_QUERY_VIEW
				}

			case DEVICE_VIEW:
				device := m.Views.Device.GetSelectedDevice()

				m.Player.SetDevice(device, m.Config)

				// Transfers playback to the newly select device.
				m.Player.Resume(m.Session, false)

			case SEARCH_QUERY_VIEW:
				if m.Views.Squery.Query() != "" {
					m.CurrentView = PLAYER_VIEW
				}

				m.Views.Squery = NewSearchQuery()
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

			m.Player.Shuffle(!state, m.Session)

		case "r":
			switch m.PlayerState().RepeatState {
			case "off":
				m.Player.Repeat(true, m.Session)
				m.PlayerState().RepeatState = "context"
			default:
				m.Player.Repeat(false, m.Session)
				m.PlayerState().RepeatState = "off"
			}
		}

		// Handles updates from the playlist list.
		if m.CurrentView == PLAYLIST_VIEW && m.Views.Playlist.PlaylistList.choice != "" {
			var cmd tea.Cmd
			model, cmd := m.Views.Playlist.PlaylistList.Update(msg)
			m.Views.Playlist.PlaylistList = model.(PlaylistList)
			return m, cmd
		}

		// Handles updates from the search list.
		if m.CurrentView == SEARCH_TYPE_VIEW && m.Views.SearchType.ListModel.choice != "" {
			var cmd tea.Cmd
			m.Views.SearchType.ListModel.list, cmd = m.Views.SearchType.ListModel.list.Update(msg)
			return m, cmd
		}

		// Handles updates from the device list.
		if m.CurrentView == DEVICE_VIEW && m.Views.Device.ListModel.choice != "" {
			var cmd tea.Cmd
			m.Views.Device.ListModel.list, cmd = m.Views.Device.ListModel.list.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// Returns the player state from the model's player view.
func (m *Program) PlayerState() *player.State {
	return m.Views.Player.State
}
