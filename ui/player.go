package ui

import (
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/session"
)

type tickMsg struct{}

const (
	POLLINTERVALMS = 500 * time.Millisecond

	PLAYERVIEW   = "playerview"
	PLAYLISTVIEW = "playlistview"
	PAUSED       = "Paused"
	NOWPLAYING   = "Now Playing"
)

type Model struct {
	Session     *session.Session
	Player      *player.Player
	Config      *config.Config
	State       *player.PlayerState
	CurrentView string
	Views       struct {
		Player struct {
			PlayingStatus      string
			PlayingStatusStyle *lipgloss.Style
			ProgressMs         int
			TrackID            string
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(POLLINTERVALMS, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *Model) View() string {
	if m.CurrentView == PLAYERVIEW {

		if !m.State.IsPlaying && m.Views.Player.PlayingStatus != PAUSED {
			m.Views.Player.PlayingStatusStyle = &PAUSEDSTYLE
			m.Views.Player.PlayingStatus = PAUSED
		}

		if m.State.IsPlaying && m.Views.Player.PlayingStatus != NOWPLAYING {
			m.Views.Player.PlayingStatusStyle = &NOWPLAYINGSTYLE
			m.Views.Player.PlayingStatus = NOWPLAYING
		}

		mainControls := "[F1 Player | F2 Playlists | F3 Search | F4 Devices]"

		track, artist,
			progressMin, progressSec,
			durationMin, durationSec := m.State.Track.InfoString(m.Config, m.Views.Player.ProgressMs)

		playerInfo := fmt.Sprintf(
			"%s\n\n%s\n\n[%sm:%ss / %sm:%ss]",
			track,
			artist,
			progressMin,
			progressSec,
			durationMin,
			durationSec,
		)

		playerStatus := m.Views.Player.PlayingStatusStyle.Render(m.Views.Player.PlayingStatus)

		return mainControls + "\n\n" + playerStatus + "\n\n" + playerInfo + "\n"
	}

	return "TODO"
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:

		go func() {
			state, err := m.Player.State(m.Session)
			if err == nil {
				m.State = state
			}
		}()

		if m.State.IsPlaying {
			m.Views.Player.ProgressMs += int(POLLINTERVALMS.Milliseconds())
		}

		progressOutOfSync := math.Abs(float64(m.State.ProgressMs-m.Views.Player.ProgressMs)) >
			float64(2*POLLINTERVALMS.Milliseconds())

		songOutOfSync := m.State.Track.ID != m.Views.Player.TrackID

		if progressOutOfSync || songOutOfSync {
			m.Views.Player.ProgressMs = m.State.ProgressMs
		}

		return m, tea.Tick(POLLINTERVALMS, func(time.Time) tea.Msg {
			return tickMsg{}
		})

	case tea.KeyMsg:
		state, err := m.Player.State(m.Session)
		if err != nil {
			return m, tea.Quit
		}
		m.State = state

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p", " ":
			if m.State.IsPlaying {
				m.Player.Pause(m.Session)
				m.Views.Player.PlayingStatusStyle = &PAUSEDSTYLE
				m.Views.Player.PlayingStatus = PAUSED
			} else {
				m.Player.Resume(m.Session, true)
				m.Views.Player.PlayingStatusStyle = &NOWPLAYINGSTYLE
				m.Views.Player.PlayingStatus = NOWPLAYING
			}

		}
	}

	return m, nil
}

func New(
	session *session.Session, player *player.Player,
	config *config.Config, inititalState *player.PlayerState,
) *Model {
	m := &Model{
		Session:     session,
		Player:      player,
		Config:      config,
		State:       inititalState,
		CurrentView: PLAYERVIEW,
	}

	if m.State.IsPlaying {
		m.Views.Player.PlayingStatusStyle = &NOWPLAYINGSTYLE
		m.Views.Player.PlayingStatus = NOWPLAYING
	} else {
		m.Views.Player.PlayingStatusStyle = &PAUSEDSTYLE
		m.Views.Player.PlayingStatus = PAUSED
	}

	m.Views.Player.ProgressMs = inititalState.ProgressMs
	m.Views.Player.TrackID = inititalState.Track.ID

	return m
}
