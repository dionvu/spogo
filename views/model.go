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

const (
	PAUSED      = "Paused"
	NO_PLAYER   = "Player Inactive"
	NOW_PLAYING = "Now Playing"
)

type Model struct {
	Session     *auth.Session
	Player      *player.Player
	Config      *config.Config
	CurrentView string
	Views       struct {
		// Tracks player state and current progress,
		// displaying information in a media player.
		Player     *PlayerView
		Playlist   *PlaylistView
		SearchType *SearchTypeView
		// SearchQuery *SearchQueryView
		Device *DeviceView

		Squery SearchQuery
	}

	Terminal Terminal

	CurrentWarning string
}

type tickMsg struct{}

func New(
	auth *auth.Session, player *player.Player,
	config *config.Config,
) *Model {
	m := &Model{
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

func (m *Model) Init() tea.Cmd {
	return tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
