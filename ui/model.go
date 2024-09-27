package ui

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"golang.org/x/term"
)

const (
	POLLING_RATE_MS          = 500 * time.Millisecond
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

type Terminal struct {
	Height int
	Width  int
}

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

	m.Views.Player = NewPlayerView(auth, player, config)
	m.Views.Playlist = NewPlaylistView(auth, config)
	m.Views.SearchType = NewSearchTypeView(auth)
	m.Views.Device = NewDeviceView(m.Session)
	// m.Views.SearchQuery = NewSearchQueryView(m.Session)

	m.Views.Squery = NewSearchQuery()

	return m
}

// Asyncronously updates the terminal dimensions.
func updateTerminalSize(terminal *Terminal) {
	// Channel to receive terminal size change signals (SIGWINCH)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)

	// If their is a change in terminal dimensions, updates terminal.
	go func() {
		for range sigCh {
			w, h := getTerminalSize()
			if w != terminal.Width || h != terminal.Height {
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()

			}

			terminal.Width, terminal.Height = w, h
		}
	}()
}

// Gets the current dimensions of the user's terminal.
func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return -1, -1
	}

	return width, height
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
