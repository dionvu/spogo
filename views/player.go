package ui

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"

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

	// Tracks time independent of state progress
	// to improve performance, periodically will
	// be checked for error.
	ProgressMs int

	StatusBar *StatusBar

	PlayerDetails *PlayerDetails

	Ascii *Ascii

	// Kept to track if progressMs is in sync with the song.
	TrackID string

	// // Tracks the current ascii uri
	// AsciiCurrentUrl string
}

func NewPlayerView(
	auth *auth.Session, player *player.Player,
	// c *config.Config,
) *PlayerView {
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "player_ascii.jpeg")

	pv := &PlayerView{
		Session: auth,
		Player:  player,
		// Config:  c,

		PlayerDetails: &PlayerDetails{},
		StatusBar:     &StatusBar{},
		Ascii:         &Ascii{FilePath: path},
	}

	pv.UpdateStateSync()

	pv.StatusBar.Update(pv.State)

	if pv.State != nil {
		pv.ProgressMs = pv.State.ProgressMs
		pv.TrackID = pv.State.Track.ID
	}

	return pv
}

// Ensures that player time progress is within 2 * polling rate.
func (pv *PlayerView) EnsureProgressSynced() {
	// Checks pv state for external pausing or playing not captured by
	// the update method.
	if pv.State != nil {
		if !pv.State.IsPlaying && pv.StatusBar.Status != PAUSED {
			pv.StatusBar.Style = &PlayerViewStyle.StatusBar.Paused
			pv.StatusBar.Status = PAUSED
		}

		if pv.State.IsPlaying && pv.StatusBar.Status != NOW_PLAYING {
			pv.StatusBar.Style = &PlayerViewStyle.StatusBar.NowPlaying
			pv.StatusBar.Status = NOW_PLAYING
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

	pv.StatusBar.Update(pv.State)
}

// Returns a string containing the entire player view, centered with
// the size dynamic to the terminal size.
func (pv *PlayerView) View(terminal Terminal) string {
	if terminal.IsSizeSmall() {
		return pv.viewSmall(terminal)
	}

	if pv.State == nil {
		return pv.viewNoState(terminal)
	}

	pv.Ascii.UpdateImage(pv.State.Track.Album.Images[0].Url)

	return fmt.Sprintf("\n\n%s\n\n%s\n\n%s\n\n%s",
		pv.Ascii.Center(AsciiFlagsNormal(), terminal),
		pv.StatusBar.Render(terminal),
		pv.PlayerDetails.Render(pv.State.Track, pv.State.Device.VolumePercent, pv.State.ShuffleState, pv.ProgressMs, terminal),
		CenterString(MainControlsRender(PLAYER_VIEW), terminal),
	)
}

func (pv *PlayerView) viewSmall(terminal Terminal) string {
	pv.Ascii.UpdateImage(pv.State.Track.Album.Images[0].Url)

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	// ascii = CenterString(ascii, terminal)
	ascii := pv.Ascii.Center(AsciiFlagsSmall(), terminal)
	statusBar := pv.StatusBar.Render(terminal)
	playerInfo := pv.PlayerDetails.Render(pv.State.Track, pv.State.Device.VolumePercent, pv.State.ShuffleState, pv.ProgressMs, terminal)

	t.AppendRows([]table.Row{
		{"\n\n" + ascii}, {"\n" + statusBar + "\n"}, {playerInfo},
	})

	return t.Render()
}

// The player view when state is nil (player and device is not active).
func (pv *PlayerView) viewNoState(terminal Terminal) string {
	if terminal.IsSizeSmall() {
		return fmt.Sprintf("\n\n%s",
			pv.StatusBar.Render(terminal))
	}

	return fmt.Sprintf("\n\n%s\n\n%s",
		CenterString(MainControlsRender(PLAYER_VIEW), terminal),
		pv.StatusBar.Render(terminal))
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
