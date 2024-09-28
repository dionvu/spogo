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
)

type PlayerView struct {
	Session *auth.Session
	Player  *player.Player
	State   *player.PlayerState

	// Tracks time independent of state progress
	// to improve performance, periodically will
	// be checked for error.
	ProgressMs int

	// The string of content to be displayed when player viewed.
	Content Content

	StatusBar *StatusBar

	PlayerDetails *PlayerDetails

	Image *Image

	// Kept to track if progressMs is in sync with the song.
	TrackID string
}

func NewPlayerView(
	auth *auth.Session, player *player.Player,
) *PlayerView {
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "player_ascii.jpeg")

	pv := &PlayerView{
		Session: auth,
		Player:  player,

		PlayerDetails: &PlayerDetails{},
		StatusBar:     &StatusBar{},
		Image:         &Image{FilePath: path},
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
func (pv *PlayerView) View(term Terminal) string {
	pv.UpdateContent(term)

	return pv.Content.CenterHorizontal(term).CenterVertical(term).String()
}

// Updates the view content based on the state of the player,
// and the current size of the terminal.
func (pv *PlayerView) UpdateContent(term Terminal) {
	switch term.IsSizeSmall() {
	case true:
		switch pv.State {
		case nil:
			pv.Content = Content(fmt.Sprintf("\n\n%s", pv.StatusBar.Render()))
		default:
			pv.Image.UpdateImage(pv.State.Track.Album.Images[0].Url)

			pv.Image.UpdateImage(pv.State.Track.Album.Images[0].Url)

			ascii := pv.Image.AsciiSmall().Content()

			statusBar := pv.StatusBar.Content()

			playerDetails := pv.PlayerDetails.Content(pv.State.Track,
				pv.State.Device.VolumePercent,
				pv.State.ShuffleState,
				pv.ProgressMs)

			pv.Content = Join([]Content{
				ascii, statusBar, playerDetails,
			}, "\n\n")

			// pv.Content = Content(container.Render())
		}

	case false:
		switch pv.State {
		case nil:
			pv.Content = Content(fmt.Sprintf("\n\n%s\n\n%s",
				MainControlsRender(PLAYER_VIEW),
				pv.StatusBar.Render()))

		default:
			pv.Image.UpdateImage(pv.State.Track.Album.Images[0].Url)

			ascii := pv.Image.AsciiNormal().Content()

			statusBar := pv.StatusBar.Content()

			playerDetails := pv.PlayerDetails.Content(pv.State.Track,
				pv.State.Device.VolumePercent,
				pv.State.ShuffleState,
				pv.ProgressMs)

			mainControls := MainControlsRender(PLAYER_VIEW)

			pv.Content = Join([]Content{
				ascii, statusBar, playerDetails, Content(mainControls),
			}, "\n\n")
		}
	}
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
