package ui

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
)

type PlayerView struct {
	Session *auth.Session
	Player  *player.Player
	Config  *config.Config
	State   *player.PlayerState

	// The title status bar indicating, playing, paused or invalid device.
	PlayingStatus      string
	PlayingStatusStyle *lipgloss.Style

	// Tracks time independent of state progress
	// to improve performance, periodically will
	// be checked for error.
	ProgressMs int

	// Kept to track if progressMs is in sync with the song.
	TrackID string
}

func NewPlayerView(
	auth *auth.Session, player *player.Player,
	config *config.Config,
) *PlayerView {
	pv := &PlayerView{
		Session: auth,
		Player:  player,
		Config:  config,
	}

	pv.UpdateStateSync()

	if pv.State != nil && pv.State.IsPlaying {
		pv.PlayingStatusStyle = &PlayerViewStyle.StatusBar.NowPlaying
		pv.PlayingStatus = NOW_PLAYING
	} else if pv.State != nil && !pv.State.IsPlaying {
		pv.PlayingStatusStyle = &PlayerViewStyle.StatusBar.Paused
		pv.PlayingStatus = PAUSED
	} else {
		pv.PlayingStatusStyle = &PlayerViewStyle.StatusBar.NoPlayer
		pv.PlayingStatus = NO_PLAYER
	}

	if pv.State != nil {
		pv.ProgressMs = pv.State.ProgressMs
		pv.TrackID = pv.State.Track.ID
	}

	return pv
}

func (pv *PlayerView) View(terminalSize int) string {
	if pv.State != nil {
		res, _ := http.Get(pv.State.Track.Album.Images[0].Url)

		cd, _ := os.UserCacheDir()
		filepath := filepath.Join(cd, config.APPNAME, "image.jpeg")

		file, _ := os.Create(filepath)

		io.Copy(file, res.Body)

		if terminalSize <= TERMINALSIZE.Small {
			return fmt.Sprintf("\n\n%s\n\n%s\n\n%s",
				AsciiView(filepath, ASCII_FLAGS_SMALL),
				PlayerStatusView(pv),
				PlayerInfoView(pv))
		}

		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
			MainControlsView(PLAYER_VIEW),
			AsciiView(filepath, ASCII_FLAGS_NORMAL),
			PlayerStatusView(pv),
			PlayerInfoView(pv))
	}

	if terminalSize <= TERMINALSIZE.Small {
		return fmt.Sprintf("\n\n%s",
			PlayerStatusView(pv))
	}

	return fmt.Sprintf("%s\n\n%s",
		MainControlsView(PLAYER_VIEW),
		PlayerStatusView(pv))
}

// Ensures that player time progress is within 2 * polling rate.
func (pv *PlayerView) EnsureSynced() {
	// Checks pv state for external pausing or playing not captured by
	// the update method.
	if pv.State != nil {
		if !pv.State.IsPlaying && pv.PlayingStatus != PAUSED {
			pv.PlayingStatusStyle = &PlayerViewStyle.StatusBar.Paused
			pv.PlayingStatus = PAUSED
		}

		if pv.State.IsPlaying && pv.PlayingStatus != NOW_PLAYING {
			pv.PlayingStatusStyle = &PlayerViewStyle.StatusBar.NowPlaying
			pv.PlayingStatus = NOW_PLAYING
		}
	}

	if pv.State.IsPlaying {
		pv.ProgressMs += int(POLLING_RATE_MS.Milliseconds())
	}

	// Syncs progress time if it differs too much (2 * Polling rate).
	if math.Abs(float64(pv.State.ProgressMs-pv.ProgressMs)) >
		float64(2*POLLING_RATE_MS.Milliseconds()) ||
		pv.State.Track.ID != pv.TrackID {

		pv.ProgressMs = pv.State.ProgressMs

		pv.TrackID = pv.State.Track.ID
	}
}

// Updates state asyncchronously to improve progress timer smoothness.
func (pv *PlayerView) UpdateStateAsync() {
	go func() {
		pv.State, _ = pv.Player.State(pv.Session)
	}()
}

// Update state synchronously for percision.
func (pv *PlayerView) UpdateStateSync() {
	pv.State, _ = pv.Player.State(pv.Session)
}

func (pv *PlayerView) PlayPause() {
	if pv.State.IsPlaying {
		pv.Player.Pause(pv.Session)
	} else {
		pv.Player.Resume(pv.Session, true)
	}

	pv.UpdateStateSync()

	if pv.State.IsPlaying {
		pv.PlayingStatusStyle = &PlayerViewStyle.StatusBar.NowPlaying
		pv.PlayingStatus = NOW_PLAYING
	} else {
		pv.PlayingStatusStyle = &PlayerViewStyle.StatusBar.Paused
		pv.PlayingStatus = PAUSED
	}
}
