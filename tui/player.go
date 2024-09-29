package ui

import (
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
)

// The view struct that displays player state
// details, the current track's album art,
// and other relevant information to the user.
type PlayerView struct {
	// The string of content to be displayed when the
	// player is viewed.
	Content Content

	// Indicates the playing status of the the track.
	StatusBar *StatusBar

	// Holds the track & artist names, the progress of the track
	// and various other options relevant to the user.
	PlayerDetails *PlayerDetails

	// Displays the alternative main views, with the
	// current view (player view) highlighted.
	ViewStatus *ViewStatus

	// Album art image of the track currently playing.
	Image *Image

	Session *auth.Session
	Player  *player.Player
	State   *player.State

	// Tracks time independent of state progress
	// to improve performance, periodically will
	// be checked for error.
	ProgressMs int

	// Kept to track if progressMs is in sync with the song.
	TrackID string
}

func NewPlayerView(
	auth *auth.Session, player *player.Player,
) *PlayerView {
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "player.jpeg")

	pv := &PlayerView{
		Session: auth,
		Player:  player,

		PlayerDetails: &PlayerDetails{},
		StatusBar:     &StatusBar{},
		ViewStatus:    &ViewStatus{},
		Image:         &Image{FilePath: path},
	}

	pv.UpdateStateSync()

	pv.StatusBar.Update(pv.State)
	pv.ViewStatus.Update(PLAYER_VIEW)

	if pv.State != nil {
		pv.ProgressMs = pv.State.ProgressMs
		pv.TrackID = pv.State.Track.ID
	}

	return pv
}

// Ensures that player time progress is within 2 * polling rate.
func (pv *PlayerView) EnsureProgressSynced() {
	if pv.State == nil {
		return
	}

	// Checks pv state for external pausing or playing not captured by
	// the update method.
	pv.StatusBar.Update(pv.State)

	if pv.State.IsPlaying {
		pv.ProgressMs += int(POLLING_RATE_MS.Milliseconds())
	}

	// Syncs progress time if it differs too much (2 * Polling rate).
	if math.Abs(float64(pv.State.ProgressMs-pv.ProgressMs)) >
		float64(10*POLLING_RATE_MS.Milliseconds()) ||
		pv.State.Track.ID != pv.TrackID {

		pv.ProgressMs = pv.State.ProgressMs
		pv.TrackID = pv.State.Track.ID
	}
}

// Updates the view content based on the state of the player,
// and the current size of the terminal.
func (pv *PlayerView) UpdateContent(term Terminal) {
	pv.Content = func() Content {
		if term.IsSizeSmall() {
			switch pv.State {
			case nil:
				return pv.StatusBar.Content()
			default:
				pv.Image.Update(pv.State.Track.Album.Images[0].Url)

				return Join([]Content{
					pv.Image.AsciiSmall().Content(),
					pv.StatusBar.Content(),
					pv.PlayerDetails.Content(pv.State.Track, pv.ProgressMs, pv.State),
				}, "\n\n")
			}
		}

		switch pv.State {
		case nil:
			return Join([]Content{
				pv.StatusBar.Content(),
				pv.ViewStatus.Content(),
			}, "\n\n")

		default:
			pv.Image.Update(pv.State.Track.Album.Images[0].Url)

			return Join([]Content{
				pv.Image.AsciiNormal().Content(),
				pv.StatusBar.Content(),
				pv.PlayerDetails.Content(pv.State.Track, pv.ProgressMs, pv.State),
				pv.ViewStatus.Content(),
			}, "\n\n")
		}
	}()
}

// PlayPause toggles playback and updates the
// StatusBar accordingly.
func (pv *PlayerView) PlayPause() {
	if pv.State == nil {
		return
	}

	switch pv.State.IsPlaying {
	case true:
		pv.Player.Pause(pv.Session)
	default:
		pv.Player.Resume(pv.Session, true)
	}

	// pv.UpdateStateSync()

	pv.StatusBar.Update(pv.State)
}

// Returns a string containing the entire player view, centered with
// the size dynamic to the terminal size.
func (pv *PlayerView) View(term Terminal) string {
	pv.UpdateContent(term)
	pv.PlayerDetails.Update(pv.ProgressMs, pv.State)

	return pv.Content.CenterHorizontal(term).CenterVertical(term).String()
}

// Update state synchronously for percision.
func (pv *PlayerView) UpdateStateSync() {
	pv.State, _ = pv.Player.State(pv.Session)
}

// Updates state continuously and asyncchronously.
func (pv *PlayerView) UpdateStateLoop() {
	go func() {
		pv.State, _ = pv.Player.State(pv.Session)
		time.Sleep(time.Second)
		pv.UpdateStateLoop()
	}()
}

// func (pv *PlayerView) UpdateStateAsync() {
// 	go func() {
// 		pv.State, _ = pv.Player.State(pv.Session)
// 	}()
// }
