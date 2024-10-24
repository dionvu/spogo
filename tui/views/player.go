package views

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/spotify/auth"
	comp "github.com/dionvu/spogo/tui/views/components"
)

const (
	PLAYER_VIEW              = "player_view"
	PLAYLIST_VIEW            = "playlist_view"
	PLAYLIST_TRACK_VIEW      = "playlist_track_view"
	ALBUM_TRACK_VIEW         = "album_track_view"
	REFRESH_VIEW             = "refresh_view"
	HELP_VIEW                = "help_view"
	TERMINAL_WARNING_VIEW    = "terminal_warning_view"
	SEARCH_VIEW_QUERY        = "search_view_query"
	SEARCH_VIEW_TYPE         = "search_view_type"
	SEARCH_VIEW_RESULTS      = "search_view_results"
	DEVICE_VIEW              = "device_view"
	REAUTH_VIEW              = "reauthentication_view"
	DEVICE_FZF_VIEW          = "device_fzf_view"
	ENABLED                  = "on"
	DISABLED                 = "off"
	PLAYER_IMAGE_FILE        = "player" + comp.FILE_EXTENSION
	UPDATE_RATE_SEC          = time.Second
	POLLING_RATE_STATE_SEC   = time.Second * 5
	VOLUME_INCREMENT_PERCENT = 5
)

// The view struct that displays player state
// details, the current track's album art,
// and other relevant information to the user.
type Player struct {
	// The string of content to be displayed when the
	// player is viewed.
	Content comp.Content

	// Indicates the playing status of the the track.
	StatusBar *StatusBar

	// Holds the track & artist names, the progress of the track
	// and various other options relevant to the user.
	PlayerDetails *PlayerDetails

	// Displays the alternative main views, with the
	// current view (player view) highlighted.
	ViewStatus *ViewStatus

	// Album art image of the track currently playing.
	Image *comp.Image

	Session *auth.Session
	Config  *config.Config
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
	auth *auth.Session, player *player.Player, cfg *config.Config,
) *Player {
	pv := &Player{
		Session: auth,
		Player:  player,
		Config:  cfg,

		PlayerDetails: &PlayerDetails{},
		StatusBar:     &StatusBar{},
		ViewStatus:    &ViewStatus{CurrentView: PLAYER_VIEW},
		Image:         &comp.Image{FilePath: filepath.Join(cfg.CachePath(), IMAGES_FOLDER_NAME, PLAYER_IMAGE_FILE)},
	}

	os.MkdirAll(filepath.Join(cfg.CachePath(), IMAGES_FOLDER_NAME), os.ModePerm)

	pv.UpdateStateSync()

	if pv.State != nil {
		pv.ProgressMs = pv.State.ProgressMs
		pv.TrackID = pv.State.Track.ID
	}

	pv.StatusBar.Style = struct {
		NowPlaying lg.Style
		Paused     lg.Style
		NoPlayer   lg.Style
	}{
		NowPlaying: lg.NewStyle().
			Bold(cfg.Player.StatusBar.Style.NowPlaying.Bold).
			Foreground(lg.Color(cfg.Player.StatusBar.Style.NowPlaying.Foreground)).
			Background(lg.Color(cfg.Player.StatusBar.Style.NowPlaying.Background)).
			PaddingLeft(1).
			PaddingRight(1),

		Paused: lg.NewStyle().
			Bold(cfg.Player.StatusBar.Style.Paused.Bold).
			Foreground(lg.Color(cfg.Player.StatusBar.Style.Paused.Foreground)).
			Background(lg.Color(cfg.Player.StatusBar.Style.Paused.Background)).
			PaddingLeft(1).
			PaddingRight(1),

		NoPlayer: lg.NewStyle().
			Bold(cfg.Player.StatusBar.Style.NoPlayer.Bold).
			Foreground(lg.Color(cfg.Player.StatusBar.Style.NoPlayer.Foreground)).
			Background(lg.Color(cfg.Player.StatusBar.Style.NoPlayer.Background)).
			PaddingLeft(1).
			PaddingRight(1),
	}

	pv.StatusBar.Update(pv.State)

	return pv
}

// Ensures that player time progress is within 2 * polling rate.
func (pv *Player) EnsureProgressSynced() {
	if pv.State == nil {
		return
	}

	// Checks pv state for external pausing or playing not captured by
	// the update method.
	pv.StatusBar.Update(pv.State)

	if pv.State.IsPlaying && pv.ProgressMs < pv.State.Track.DurationMs {
		pv.ProgressMs += int(UPDATE_RATE_SEC.Milliseconds())
	}

	// Syncs progress time if it differs too much (5 * Polling rate).
	if math.Abs(float64(pv.State.ProgressMs-pv.ProgressMs)) >
		float64(5*UPDATE_RATE_SEC.Milliseconds()) ||
		pv.State.Track.ID != pv.TrackID {

		pv.ProgressMs = pv.State.ProgressMs
		pv.TrackID = pv.State.Track.ID
	}

	// Updates the progress percisly when player is paused.
	if !pv.State.IsPlaying {
		pv.ProgressMs = pv.State.ProgressMs
		pv.TrackID = pv.State.Track.ID
	}
}

// Updates the view content based on the state of the player,
// and the current size of the terminal.
func (pv *Player) UpdateContent(term comp.Terminal) {
	pv.Content = func() comp.Content {
		if term.HeightIsSmall() {
			switch pv.State {
			case nil:
				return "Ctrl+D to select a device\n\n" + pv.StatusBar.Content()
			default:
				pv.Image.Update(pv.State.Track.Album.Images[0].Url)

				return comp.Join([]comp.Content{
					pv.Image.AsciiSmall(pv.Config).Content(),
					pv.StatusBar.Content(),
					pv.PlayerDetails.Content(pv.State.Track, pv.ProgressMs, pv.State),
				}, "\n\n")
			}
		} else if term.WidthIsSmall() {
			switch pv.State {
			case nil:
				return "Ctrl+D to select a device\n\n" + pv.StatusBar.Content()
			default:
				pv.Image.Update(pv.State.Track.Album.Images[0].Url)

				return comp.Join([]comp.Content{
					pv.Image.AsciiNormal(pv.Config).Content(),
					pv.StatusBar.Content(),
					pv.PlayerDetails.Content(pv.State.Track, pv.ProgressMs, pv.State),
				}, "\n\n")
			}
		}

		switch pv.State {
		case nil:
			return comp.Join([]comp.Content{
				"",
				pv.StatusBar.Content().Append('\n', 10).Prepend('\n', 12),
				pv.ViewStatus.Content(pv.Config),
			}, "\n\n")

		default:
			if len(pv.State.Track.Album.Images) > 0 {
				pv.Image.Update(pv.State.Track.Album.Images[0].Url)
			}

			return comp.Join([]comp.Content{
				"",
				pv.Image.AsciiNormal(pv.Config).Content(),
				pv.StatusBar.Content(),
				pv.PlayerDetails.Content(pv.State.Track, pv.ProgressMs, pv.State),
				pv.ViewStatus.Content(pv.Config),
			}, "\n\n")
		}
	}()
}

// PlayPause toggles playback and updates the
// StatusBar accordingly.
func (pv *Player) PlayPause() error {
	if pv.State == nil {
		return nil
	}

	switch pv.State.IsPlaying {
	case true:
		// Updates to ensure player updates immediately
		// since state only updates every POLLING_RATE seconds.
		pv.State.IsPlaying = !pv.State.IsPlaying

		err := pv.Player.Pause(pv.Session)
		if err != nil {
			return err
		}

	default:
		pv.State.IsPlaying = !pv.State.IsPlaying

		err := pv.Player.Resume(pv.Session, true)
		if err != nil {
			return err
		}
	}

	pv.StatusBar.Update(pv.State)

	return nil
}

// Returns a string containing the entire player view, centered with
// the size dynamic to the terminal size.
func (pv *Player) View(term comp.Terminal) string {
	pv.UpdateContent(term)
	pv.PlayerDetails.Update(pv.ProgressMs, pv.State)

	return pv.Content.CenterHorizontal(term).CenterVertical(term).String()
}

// Update state synchronously for percision.
func (pv *Player) UpdateStateSync() {
	pv.State, _ = pv.Player.State(pv.Session)
}

// Updates state continuously and asyncchronously, runs reauthentication if requried.
func (pv *Player) UpdateStateLoop(session *auth.Session, config *config.Config) {
	go func() {
		var err error

		pv.State, err = pv.Player.State(pv.Session)
		if err != nil {
			err = session.Reauth(config)
			if err != nil {
				log.Fatal("ERR: Failed to reauthenticate: ", err)
				errors.Log(err)
			}
		}

		time.Sleep(POLLING_RATE_STATE_SEC)

		pv.UpdateStateLoop(session, config)
	}()
}

type PlayerDetails struct {
	Track         string
	Artists       string
	ProgressMin   string
	ProgressSec   string
	DurationMin   string
	DurationSec   string
	VolumePercent string
}

// Updates PlayerDetails with information in given state and track.
func (pd *PlayerDetails) Update(progressMs int, state *player.State) {
	if state == nil || state.Device == nil || state.Track == nil {
		return
	}

	pd.Track = state.Track.Name
	pd.Artists = ""
	pd.ProgressSec = strconv.Itoa(((progressMs / 1000) % 60))
	pd.ProgressMin = strconv.Itoa((progressMs / 1000) / 60)
	pd.DurationSec = strconv.Itoa((state.Track.DurationMs / 1000) % 60)
	pd.DurationMin = strconv.Itoa((state.Track.DurationMs / 1000) / 60)
	pd.VolumePercent = strconv.Itoa(state.Device.VolumePercent)

	for i := 0; i < len(state.Track.Artists); i++ {
		pd.Artists += state.Track.Artists[i].Name
		if i != len(state.Track.Artists)-1 {
			pd.Artists += ", "
		}
	}

	for _, time := range []*string{&pd.ProgressSec, &pd.DurationSec} {
		if len(*time) == 1 {
			*time = "0" + *time
		}
	}
}

// Renders the player details as a content string.
func (pd *PlayerDetails) Content(track *spotify.Track, progressMs int, state *player.State) comp.Content {
	pd.Update(progressMs, state)

	title := comp.Content(fmt.Sprintf("%s - %s", pd.Track, pd.Artists))

	timer := comp.Content(fmt.Sprintf("%sm:%ss / %sm:%ss", pd.ProgressMin, pd.ProgressSec, pd.DurationMin, pd.DurationSec))

	var repeat string
	switch state.RepeatState {
	case DISABLED:
		repeat = DISABLED
	default:
		repeat = ENABLED
	}

	var shuffle string
	switch state.ShuffleState {
	case true:
		shuffle = ENABLED
	default:
		shuffle = DISABLED
	}

	options := comp.Content(fmt.Sprintf("Vol: %s%%  Sfl: %v  Rep: %v", pd.VolumePercent, shuffle, repeat))

	return comp.Join([]comp.Content{title, timer, options}, "\n\n")
}

// The title status bar indicating whether the player is
// playing, paused or an invalid device is selected.
type StatusBar struct {
	Style struct {
		NowPlaying lg.Style
		Paused     lg.Style
		NoPlayer   lg.Style
	}
	Status       string
	CurrentStyle lg.Style
}

// Renders the status bar as a string.
func (sb *StatusBar) Render() string {
	return sb.CurrentStyle.Render(sb.Status)
}

// Renders the status bar as a content string.
func (sb *StatusBar) Content() comp.Content {
	return comp.Content(sb.CurrentStyle.Render(sb.Status))
}

// Updates the status bar given the player's state.
func (sb *StatusBar) Update(state *player.State) {
	const (
		PAUSED      = "Paused"
		NO_PLAYER   = "Player Inactive"
		NOW_PLAYING = "Now Playing"
	)

	if state != nil && state.IsPlaying {
		sb.CurrentStyle = sb.Style.NowPlaying
		sb.Status = NOW_PLAYING
	} else if state != nil && !state.IsPlaying {
		sb.CurrentStyle = sb.Style.Paused
		sb.Status = PAUSED
	} else {
		sb.CurrentStyle = sb.Style.NoPlayer
		sb.Status = NO_PLAYER
	}
}
