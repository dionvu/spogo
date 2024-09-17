package ui

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
	"golang.org/x/term"
)

const (
	POLLING_RATE_MS          = 500 * time.Millisecond
	VOLUME_INCREMENT_PERCENT = 5
)

var TERMINALSIZE = struct {
	Small  int
	Normal int
}{
	Small:  30,
	Normal: 40,
}

const (
	PLAYER_VIEW   = "PLAYER_VIEW"
	PLAYLIST_VIEW = "PLAYLIST_VIEW"
	REFRESH_VIEW  = "REFRESH_VIEW"
	HELP_VIEW     = "HELP_VIEW"
	PAUSED        = "Paused"
	NO_PLAYER     = "Player Inactive"
	NOW_PLAYING   = "Now Playing"
)

type Model struct {
	Session     *session.Session
	Player      *player.Player
	Config      *config.Config
	CurrentView string
	Views       struct {
		// Tracks player state and current progress,
		// displaying information in a media player.
		Player   *PlayerView
		Playlist *PlaylistView
	}

	Terminal struct {
		Height int
		Width  int
	}
}

type tickMsg struct{}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateTerminalSize(&m.Terminal.Width, &m.Terminal.Height)

	switch msg := msg.(type) {
	case tickMsg:
		m.Views.Player.UpdateStateAsync()

		// If state is unaccessible, likely due to user closing
		// their playerback device, and attempt reconnect to closed device.
		if m.Views.Player.State == nil {
			m.Views.Player.PlayingStatusStyle = &PlayerViewStyle.StatusBar.NoPlayer
			m.Views.Player.PlayingStatus = NO_PLAYER

			m.Player.Resume(m.Session, false)

			return m, tea.Tick(4*POLLING_RATE_MS, func(time.Time) tea.Msg {
				return tickMsg{}
			})
		}

		m.Views.Player.EnsureSynced()

		return m, tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
			return tickMsg{}
		})

	case tea.KeyMsg:
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

		case "f1":
			m.CurrentView = PLAYER_VIEW

		case "f2":
			m.CurrentView = PLAYLIST_VIEW

		case "f5":
			m.CurrentView = HELP_VIEW

		case "enter":
			if m.CurrentView == PLAYLIST_VIEW {
				i, ok := m.Views.Playlist.ListModel.list.SelectedItem().(Item)
				if ok {
					m.Views.Playlist.ListModel.choice = string(i)
				}
			}

		case "r":
			go func() {
				view := m.CurrentView

				m.CurrentView = REFRESH_VIEW
				time.Sleep(POLLING_RATE_MS)
				m.CurrentView = view

				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()

		}

		if m.CurrentView == PLAYLIST_VIEW {
			var cmd tea.Cmd
			m.Views.Playlist.ListModel.list, cmd = m.Views.Playlist.ListModel.list.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func New(
	session *session.Session, player *player.Player,
	config *config.Config,
) *Model {
	m := &Model{
		Session:     session,
		Player:      player,
		Config:      config,
		CurrentView: PLAYER_VIEW,
	}

	// A nil state here could be due to an inactive device.
	// Transfers playback to inactive player.
	if initialState, _ := player.State(session); initialState == nil {
		player.Resume(session, false)
	}

	// A nil state due to invalid device will be handled
	// after view is updated.
	m.Views.Player = NewPlayerView(session, player, config)
	m.Views.Playlist = NewPlaylistView(session, config)

	m.Terminal.Width, m.Terminal.Height = getTerminalSize()
	return m
}

func (m *Model) View() string {
	switch m.CurrentView {
	case PLAYER_VIEW:
		return m.Views.Player.View(m.Terminal.Height)
	case PLAYLIST_VIEW:
		return m.Views.Playlist.View(m.Views.Player, m.Terminal.Width)
	case REFRESH_VIEW:
		return "Refreshing..."
	case HELP_VIEW:
		return MainControlsView(HELP_VIEW) + "\n\n" + padLines(m.Config.HelpString(), TAB_WIDTH)
	default:
		return "TODO"
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return -1, -1
	}
	return width, height
}

func updateTerminalSize(width *int, height *int) {
	// Channel to receive terminal size change signals (SIGWINCH)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)

	// Listen for size change events
	go func() {
		for range sigCh {
			w, h := getTerminalSize()
			if w != *width || h != *height {
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()

				// Avoids bubbletea glitching out the ascii.
				// time.Sleep(POLLING_RATE_MS * 2)
			}
			*width, *height = w, h
		}
	}()
}
