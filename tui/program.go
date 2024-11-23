package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify/auth"
	"github.com/dionvu/spogo/tui/views"
	comp "github.com/dionvu/spogo/tui/views/components"
)

const (
	UPDATE_RATE_SEC        = time.Second
	POLLING_RATE_STATE_SEC = time.Second * 5
)

// Handles updates for views and controls
// which views are to be displayed.
type Program struct {
	currentView string

	// Tracks player state and current progress,
	// displaying information in a media player.
	playerView views.Player

	// Displays the user's playlists in a list
	// format, allowing the user to select one
	// to transfer playback to.
	playlistView views.Playlist
	search       views.Search
	help         views.Help

	terminal comp.Terminal

	session *auth.Session

	player *player.Player

	config *config.Config
}

type tickMsg struct{}

func New(
	auth *auth.Session, player *player.Player,
	config *config.Config,
) *Program {
	p := &Program{
		session:     auth,
		player:      player,
		config:      config,
		currentView: views.PLAYER_VIEW,
		help:        views.NewHelpView(),
	}

	if initialState, _ := player.State(auth); initialState == nil {
		player.Resume(auth, false)
	}

	p.terminal.Width, p.terminal.Height = comp.GetTerminalSize()

	p.playerView = views.NewPlayerView(auth, player, config)
	p.playlistView = views.NewPlaylistView(auth, p.terminal, config)
	p.search = views.NewSearch(p.session, p.config)

	return p
}

func (program *Program) Run() error {
	tp := tea.NewProgram(program, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := tp.Run(); err != nil {
		return err
	}

	return nil
}

func (p *Program) Init() tea.Cmd {
	p.playerView.UpdateStateLoop(p.session, p.config)

	return tea.Tick(UPDATE_RATE_SEC, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
