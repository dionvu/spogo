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

	// Tracks the current ascii uri
	AsciiCurrentUrl string
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

func (pv *PlayerView) View(terminal Terminal) string {
	if terminal.IsSizeSmall() {
		return pv.viewSmall()
	}

	if pv.State != nil {
		if pv.State.Track.Album.Images[0].Url != pv.AsciiCurrentUrl {
			err := cacheImage(pv.State.Track.Album.Images[0].Url, pv.CachedImagePath())
			if err != nil {
				return ""
			}
			pv.AsciiCurrentUrl = pv.State.Track.Album.Images[0].Url
		}

		return fmt.Sprintf("\n\n%s\n\n%s\n\n%s\n\n%s",
			MainControlsRender(PLAYER_VIEW),
			padLines(AsciiRender(pv.CachedImagePath(), AsciiFlagsNormal()), TAB_WIDTH),
			PlayerStatusView(pv),
			PlayerInfoView(pv))
	}

	return fmt.Sprintf("\n\n%s\n\n%s",
		MainControlsRender(PLAYER_VIEW),
		PlayerStatusView(pv))
}

func (pv *PlayerView) viewSmall() string {
	if pv.State.Track.Album.Images[0].Url != pv.AsciiCurrentUrl {
		err := cacheImage(pv.State.Track.Album.Images[0].Url, pv.CachedImagePath())
		if err != nil {
			return ""
		}
		pv.AsciiCurrentUrl = pv.State.Track.Album.Images[0].Url
	}

	if pv.State != nil {
		return fmt.Sprintf("\n\n%s\n\n%s\n\n%s",
			padLines(AsciiRender(pv.CachedImagePath(), AsciiFlagsSmall()), TAB_WIDTH),
			PlayerStatusView(pv),
			PlayerInfoView(pv))
	}

	return fmt.Sprintf("\n\n%s",
		PlayerStatusView(pv))
}

// Ensures that player time progress is within 2 * polling rate.
func (pv *PlayerView) EnsureProgressSynced() {
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
	if pv.State == nil {
		return
	}

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

var PlayerStatusView = func(pv *PlayerView) string {
	return padLines(pv.PlayingStatusStyle.Render(pv.PlayingStatus), 4)
}

var PlayerInfoView = func(pv *PlayerView) string {
	if pv.State == nil {
		return "invalid player state"
	}

	track, artist,
		progressMin, progressSec,
		durationMin, durationSec := pv.State.Track.InfoString(pv.Config, pv.ProgressMs)

	var shuffle string

	if pv.State.ShuffleState {
		shuffle = "on"
	} else {
		shuffle = "off"
	}

	return padLines(fmt.Sprintf(
		"%s - %s\n\n%sm:%ss / %sm:%ss\n\nvol: %v%% sfl: %v",
		track,
		artist,
		progressMin,
		progressSec,
		durationMin,
		durationSec,
		pv.State.Device.VolumePercent,
		shuffle,
	), TAB_WIDTH)
}

func cacheImage(url string, path string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (pv *PlayerView) CachedImagePath() string {
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "player_ascii.jpeg")
	return path
}
