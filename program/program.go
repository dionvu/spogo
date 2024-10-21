package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/components"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/views"
)

const (
	UPDATE_RATE_SEC          = time.Second
	POLLING_RATE_STATE_SEC   = time.Second * 5
	VOLUME_INCREMENT_PERCENT = 5
)

// The struct that integrates every view
// into a single cohesive program. Handles
// updates for views and controls which views
// are to be displayed.
type Program struct {
	CurrentView string

	// Tracks player state and current progress,
	// displaying information in a media player.
	Player *views.Player

	// Displays the user's playlists in a list
	// format, allowing the user to select one
	// to transfer playback to.
	Playlist *views.Playlist

	// Allows the user to search for tracks,
	// albums, etc., depending on the selection.
	// SearchType *SearchTypeView

	// SearchQuery *SearchQueryView
	Device *views.Device

	Search views.Search

	// The programs's current terminal size, this
	// is updated consistantly.
	Terminal components.Terminal

	// Stuff necessary to access the spotify api.
	session *auth.Session
	player  *player.Player

	// Configuration options.
	Config *config.Config
}

type tickMsg struct{}

func New(
	auth *auth.Session, player *player.Player,
	config *config.Config,
) *Program {
	p := &Program{
		session:     auth,
		player:      player,
		Config:      config,
		CurrentView: views.PLAYER_VIEW,
	}

	// A nil state here could be due to an inactive device.
	// Transfers playback to inactive player.
	if initialState, _ := player.State(auth); initialState == nil {
		player.Resume(auth, false)
	}

	p.Terminal.Width, p.Terminal.Height = components.GetTerminalSize()

	p.Player = views.NewPlayerView(auth, player)
	p.Playlist = views.NewPlaylistView(auth, p.Terminal)
	p.Device = views.NewDeviceView(p.session)
	p.Search = views.NewSearch(p.session)

	return p
}

func (program *Program) Run() error {
	tp := tea.NewProgram(program, tea.WithAltScreen())
	if _, err := tp.Run(); err != nil {
		return err
	}

	return nil
}

func (p *Program) Init() tea.Cmd {
	p.Player.UpdateStateLoop()
	return tea.Tick(UPDATE_RATE_SEC, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
