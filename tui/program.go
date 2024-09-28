package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
)

const (
	POLLING_RATE_MS          = time.Second
	VOLUME_INCREMENT_PERCENT = 5
)

// The struct that integrates every view
// into a single cohesive program. Handles
// updates for views and controls which views
// are to be displayed.
type Program struct {
	CurrentView string
	Views       struct {
		// Tracks player state and current progress,
		// displaying information in a media player.
		Player *PlayerView

		// Displays the user's playlists in a list
		// format, allowing the user to select one
		// to transfer playback to.
		Playlist *PlaylistView

		// Allows the user to search for tracks,
		// albums, etc., depending on the selection.
		SearchType *SearchTypeView

		// SearchQuery *SearchQueryView
		Device *DeviceView

		Squery SearchQuery
	}

	Terminal Terminal

	CurrentWarning string

	Session *auth.Session
	Player  *player.Player
	Config  *config.Config
}

type tickMsg struct{}

func New(
	auth *auth.Session, player *player.Player,
	config *config.Config,
) *Program {
	m := &Program{
		Session:     auth,
		Player:      player,
		Config:      config,
		CurrentView: PLAYER_VIEW,
	}

	// A nil state here could be due to an inactive device.
	// Transfers playback to inactive player.
	if initialState, _ := player.State(auth); initialState == nil {
		player.Resume(auth, false)
	}

	m.Terminal.Width, m.Terminal.Height = getTerminalSize()

	m.Views.Player = NewPlayerView(auth, player)
	m.Views.Playlist = NewPlaylistView(auth, config)
	m.Views.SearchType = NewSearchTypeView(auth)
	m.Views.Device = NewDeviceView(m.Session)

	m.Views.Squery = NewSearchQuery()

	return m
}

func (m *Program) Init() tea.Cmd {
	return tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
