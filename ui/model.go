package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
)

const POLLING_RATE_MS = 500 * time.Millisecond

const (
	PLAYER_VIEW   = "PLAYER_VIEW"
	PLAYLIST_VIEW = "PLAYLIST_VIEW"
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
}

type tickMsg struct{}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.Views.Player.UpdateStateAsync()

		// If state is unaccessible, likely due to user closing
		// their playerback device, and attempt reconnect to closed device.
		if m.Views.Player.State == nil {
			m.Views.Player.PlayingStatusStyle = &NO_PLAYER_STYLE
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

		case "p", " ":
			m.Views.Player.UpdateStateSync()
			m.Views.Player.PlayPause()

		case "f1":
			m.CurrentView = PLAYER_VIEW

		case "f2":
			m.CurrentView = PLAYLIST_VIEW

		case "enter":
			if m.CurrentView == PLAYLIST_VIEW {
				i, ok := m.Views.Playlist.List.list.SelectedItem().(Item)
				if ok {
					m.Views.Playlist.List.choice = string(i)
				}
			}

		}

		if m.CurrentView == PLAYLIST_VIEW {
			var cmd tea.Cmd
			m.Views.Playlist.List.list, cmd = m.Views.Playlist.List.list.Update(msg)
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
	m.Views.Playlist = NewPlaylistView(session)

	return m
}

func (m *Model) View() string {
	switch m.CurrentView {
	case PLAYER_VIEW:
		return m.Views.Player.View()
	case PLAYLIST_VIEW:
		return m.Views.Playlist.View(m.Views.Player)
	default:
		return "TODO"
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
