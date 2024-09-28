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
	"github.com/jedib0t/go-pretty/v6/table"
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
		return pv.viewSmall(terminal)
	}

	if pv.State == nil {
		return fmt.Sprintf("\n\n%s\n\n%s",
			CenterString(MainControlsRender(PLAYER_VIEW), terminal),
			CenterString(PlayerStatusView(pv), terminal))
	}

	AsciiNewUrl := pv.State.Track.Album.Images[0].Url

	if AsciiNewUrl != pv.AsciiCurrentUrl {
		cacheImage(AsciiNewUrl, pv.CachedImagePath())

		pv.AsciiCurrentUrl = AsciiNewUrl
	}

	ascii, _ := AsciiRender(pv.CachedImagePath(), AsciiFlagsNormal())

	return fmt.Sprintf("\n\n%s\n\n%s\n\n%s\n\n%s",
		CenterString(ascii, terminal),
		CenterString(PlayerStatusView(pv), terminal),
		PlayerInfoRender(pv, terminal),
		CenterString(MainControlsRender(PLAYER_VIEW), terminal),
	)
}

func (pv *PlayerView) viewSmall(terminal Terminal) string {
	if pv.State == nil {
		return fmt.Sprintf("\n\n%s",
			CenterString(PlayerStatusView(pv), terminal))
	}

	err := cacheImage(pv.State.Track.Album.Images[0].Url, pv.CachedImagePath())
	if err != nil {
		return ""
	}
	pv.AsciiCurrentUrl = pv.State.Track.Album.Images[0].Url

	ascii, err := AsciiRender(pv.CachedImagePath(), AsciiFlagsSmall())
	if err != nil {
		ascii = "Ascii image unavailable"
	}

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	line1 := CenterString(ascii,
		terminal)
	line2 := CenterString(PlayerStatusView(pv),
		terminal, -1)
	line3 := PlayerInfoRender(pv, terminal)

	t.AppendRows([]table.Row{
		{"\n\n" + line1},
		{"\n" + line2 + "\n"},
		{line3},
	})

	return t.Render()
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
	return pv.PlayingStatusStyle.Render(pv.PlayingStatus)
}

func PlayerInfoRender(pv *PlayerView, terminal Terminal) string {
	track, artist,
		progressMin, progressSec,
		durationMin, durationSec := pv.State.Track.InfoString(pv.Config, pv.ProgressMs)

	var shuffle string

	if pv.State.ShuffleState {
		shuffle = "on"
	} else {
		shuffle = "off"
	}

	line1 := CenterString(
		fmt.Sprintf("%s - %s", track, artist),
		terminal, -1)

	line2 := CenterString(fmt.Sprintf("%sm:%ss / %sm:%ss",
		progressMin, progressSec, durationMin, durationSec),
		terminal, -1)

	line3 := CenterString(fmt.Sprintf("vol: %v%% sfl: %v", pv.State.Device.VolumePercent, shuffle),
		terminal, -1)

	return line1 + "\n\n" + line2 + "\n\n" + line3
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
