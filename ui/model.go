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
		Device     *DeviceView
	}

	Terminal Terminal

	CurrentWarning string
}

type tickMsg struct{}

type Terminal struct {
	Height int
	Width  int
}

var TERMINALSIZE = struct {
	Small  int
	Normal int
}{
	Small:  30,
	Normal: 40,
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

	// A nil state due to invalid device will be handled
	// after view is updated.
	m.Views.Player = NewPlayerView(auth, player, config)
	m.Views.Playlist = NewPlaylistView(auth, config)
	m.Views.SearchType = NewSearchTypeView(auth)
	m.Views.Device = NewDeviceView(m.Session)

	m.Terminal.Width, m.Terminal.Height = getTerminalSize()
	return m
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

			}
			*width, *height = w, h
		}
	}()
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(POLLING_RATE_MS, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
