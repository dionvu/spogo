package ui

import (
	"log"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Handles updates associate with the current selected view.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Terminal.UpdateSize()

	if !m.Terminal.IsValid() {
		m.CurrentView = TERMINAL_WARNING_VIEW
	}

	if m.CurrentView == TERMINAL_WARNING_VIEW && m.Terminal.IsValid() {
		m.CurrentView = PLAYER_VIEW
	}

	switch msg := msg.(type) {
	case tickMsg:
		m.Views.Player.UpdateStateAsync()

		// If state is unaccessible, likely due to user closing
		// their playerback device, and attempt reconnect to closed device.
		if m.Views.Player.State == nil {
			// m.Views.Player.PlayingStatusStyle = &PlayerViewStyle.StatusBar.NoPlayer
			// m.Views.Player.PlayingStatus = NO_PLAYER
			m.Views.Player.StatusBar.Update(m.Views.Player.State)

			m.Player.Resume(m.Session, false)

			return m, tea.Tick(4*POLLING_RATE_MS, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		m.Views.Player.EnsureProgressSynced()

		return m, tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
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
			m.Views.Player.UpdateStateSync()
			m.Views.Player.PlayPause()

		case "[":
			vol := m.Views.Player.State.Device.VolumePercent
			m.Views.Player.Player.SetVolume(m.Session, vol-VOLUME_INCREMENT_PERCENT)

		case "]":
			vol := m.Views.Player.State.Device.VolumePercent
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
				if i, ok := m.Views.Playlist.PlaylistListModel.list.SelectedItem().(Item); ok {
					m.Views.Playlist.PlaylistListModel.choice = string(i)
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

		case "r":
			// Refreshes the terminal fixing any visual glitches. This doesn't yet force any
			// updates to, for example, listed playlist devices.
			go func() {
				view := m.CurrentView

				m.CurrentView = REFRESH_VIEW
				time.Sleep(POLLING_RATE_MS)
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
			state := m.Views.Player.State.ShuffleState
			m.Player.Shuffle(!state, m.Session)

		}

		// Handles updates from the playlist list.
		if m.CurrentView == PLAYLIST_VIEW && m.Views.Playlist.PlaylistListModel.choice != "" {
			var cmd tea.Cmd
			m.Views.Playlist.PlaylistListModel.list, cmd = m.Views.Playlist.PlaylistListModel.list.Update(msg)
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