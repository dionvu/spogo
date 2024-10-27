package views

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Delta456/box-cli-maker/v2"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/spotify/auth"
	comp "github.com/dionvu/spogo/tui/views/components"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	GLOBAL_VIEW_WIDTH     = 80
	PLAYER_VIEW           = "player_view"
	PLAYLIST_VIEW         = "playlist_view"
	PLAYLIST_TRACK_VIEW   = "playlist_track_view"
	ALBUM_TRACK_VIEW      = "album_track_view"
	REFRESH_VIEW          = "refresh_view"
	HELP_VIEW             = "help_view"
	TERMINAL_WARNING_VIEW = "terminal_warning_view"
	SEARCH_VIEW_QUERY     = "search_view_query"
	SEARCH_VIEW_TYPE      = "search_view_type"
	SEARCH_VIEW_RESULTS   = "search_view_results"
	DEVICE_VIEW           = "device_view"
	REAUTH_VIEW           = "reauthentication_view"
	DEVICE_FZF_VIEW       = "device_fzf_view"

	UPDATE_RATE_SEC          = time.Second
	POLLING_RATE_STATE_SEC   = time.Second * 5
	PLAYER_MAX_CHAR_NORMAL   = 50
	VOLUME_INCREMENT_PERCENT = 5
	PLAYER_IMAGE_FILE        = "player" + comp.FILE_EXTENSION
	ENABLED                  = "on"
	DISABLED                 = "off"
)

var Box = box.New(box.Config{Px: 3, Py: 1, Type: "Round", Color: "HiGreen", TitlePos: "Bottom"})

// The view struct that displays player state
// details, the current track's album art.
type Player struct {
	// Indicates the playing status of the the track.
	statusBar *statusBar

	// Holds the track & artist names, the progress of the track
	// and various other options relevant to the user.
	playerDetails *PlayerDetails

	// Album art image of the track currently playing.
	image *comp.Image

	State   *player.State
	session *auth.Session
	config  *config.Config
	player  *player.Player

	// Tracks time independent of state progress
	// to improve performance, periodically will
	// be checked for error.
	progressMs int

	// Kept to track if progressMs is in sync with the song.
	trackID string
}

func (pv *Player) UpdateStatusBar(s *player.State) {
	pv.statusBar.Update(s)
}

func NewPlayerView(
	auth *auth.Session, player *player.Player, cfg *config.Config,
) Player {
	pv := Player{
		session: auth,
		player:  player,
		config:  cfg,

		playerDetails: &PlayerDetails{},
		statusBar:     &statusBar{},
		image:         &comp.Image{FilePath: filepath.Join(cfg.CachePath(), IMAGES_FOLDER_NAME, PLAYER_IMAGE_FILE)},
	}

	os.MkdirAll(filepath.Join(cfg.CachePath(), IMAGES_FOLDER_NAME), os.ModePerm)

	pv.UpdateStateSync()

	if pv.State != nil {
		pv.progressMs = pv.State.ProgressMs
		pv.trackID = pv.State.Track.ID
	}

	pv.statusBar.Style = struct {
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

	pv.statusBar.Update(pv.State)

	return pv
}

// Ensures that player time progress is within 2 * polling rate.
func (pv *Player) EnsureProgressSynced() {
	if pv.State == nil {
		return
	}

	// Checks pv state for external pausing or playing not captured by
	// the update method.
	pv.statusBar.Update(pv.State)

	if pv.State.IsPlaying && pv.progressMs < pv.State.Track.DurationMs {
		pv.progressMs += int(UPDATE_RATE_SEC.Milliseconds())
	}

	// Syncs progress time if it differs too much (5 * Polling rate).
	if math.Abs(float64(pv.State.ProgressMs-pv.progressMs)) >
		float64(5*UPDATE_RATE_SEC.Milliseconds()) ||
		pv.State.Track.ID != pv.trackID {

		pv.progressMs = pv.State.ProgressMs
		pv.trackID = pv.State.Track.ID
	}

	// Updates the progress percisly when player is paused.
	if !pv.State.IsPlaying {
		pv.progressMs = pv.State.ProgressMs
		pv.trackID = pv.State.Track.ID
	}
}

// PlayPause toggles playback and updates the
// statusBar  accordingly.
func (pv *Player) PlayPause() error {
	if pv.State == nil {
		return nil
	}

	switch pv.State.IsPlaying {
	case true:
		// Updates to ensure player updates immediately
		// since state only updates every POLLING_RATE seconds.
		pv.State.IsPlaying = !pv.State.IsPlaying

		err := pv.player.Pause(pv.session)
		if err != nil {
			return err
		}

	default:
		pv.State.IsPlaying = !pv.State.IsPlaying

		err := pv.player.Resume(pv.session, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns a string containing the entire player view, centered with
// the size dynamic to the terminal size.
func (pv *Player) View(term comp.Terminal) string {
	pv.playerDetails.Update(pv.progressMs, pv.State)

	if term.HeightIsSmall() {
		content := func() comp.Content {
			switch pv.State {
			case nil:
				return "Ctrl+D to select a device\n\n" + pv.statusBar.Content()
			default:
				pv.image.Update(pv.State.Track.Album.Images[0].Url)

				return comp.Join([]comp.Content{
					pv.image.AsciiSmall(pv.config).Content(),
					pv.statusBar.Content(),
					pv.playerDetails.Content(pv.State.Track, pv.progressMs, pv.State),
				}, "\n\n")
			}
		}()
		return content.CenterVertical(term).PadLinesLeft(3).String()
	}

	if term.WidthIsSmall() {
		content := func() comp.Content {
			switch pv.State {
			case nil:
				return "Ctrl+D to select a device\n\n" + pv.statusBar.Content()
			default:
				pv.image.Update(pv.State.Track.Album.Images[0].Url)

				return comp.Join([]comp.Content{
					pv.image.AsciiNormal(pv.config).Content(),
					pv.statusBar.Content(),
					pv.playerDetails.Content(pv.State.Track, pv.progressMs, pv.State),
				}, "\n\n")
			}
		}()

		return content.CenterVertical(term).PadLinesLeft(3).String()
	}

	content := func() comp.Content {
		switch pv.State {
		case nil:
			c := comp.Join([]comp.Content{
				"'Ctrl+D' to select a playback device",
				pv.statusBar.Content().Prepend(NL, 1),
				comp.InvisibleBarV(12),
			})

			mainContainer := func() comp.Content {
				t := comp.NewDefaultTable()

				img := comp.Image{FilePath: filepath.Join(pv.config.CachePath(), IMAGES_FOLDER_NAME, "temp"+comp.FILE_EXTENSION)}
				img.Update("https://i.pinimg.com/736x/ad/7a/16/ad7a164adabc065fae659a5b9dce9f69.jpg")
				t.AppendRow(table.Row{
					img.AsciiNormal(pv.config),
					c.PadLinesLeft(3),
				})

				return comp.Content(t.Render())
			}()

			return comp.Content(Box.String(
				ViewStatus{CurrentView: PLAYER_VIEW}.Content(pv.config).String(),
				comp.InvisibleBar(GLOBAL_VIEW_WIDTH).Append('\n', 1).String()+mainContainer.Append(NL, 1).String(),
			))

		default:
			if len(pv.State.Track.Album.Images) > 0 {
				pv.image.Update(pv.State.Track.Album.Images[0].Url)
			}

			c := comp.Join([]comp.Content{
				pv.playerDetails.Content(pv.State.Track, pv.progressMs, pv.State).Prepend(NL, 3).AdjustFit(PLAYER_MAX_CHAR_NORMAL),
				pv.statusBar.Content(),
			}, "\n\n")

			mainContainer := func() comp.Content {
				t := comp.NewDefaultTable()

				t.AppendRow(table.Row{
					pv.image.AsciiNormal(pv.config).Content(),
					c.PadLinesLeft(3),
				})

				return comp.Content(t.Render())
			}()

			return comp.Content(Box.String(
				ViewStatus{CurrentView: PLAYER_VIEW}.Content(pv.config).String(),
				comp.InvisibleBar(GLOBAL_VIEW_WIDTH).Append('\n', 1).String()+mainContainer.Append(NL, 1).String(),
			))
		}
	}()

	return content.CenterHorizontal(term).CenterVertical(term).String()
}

// Update state synchronously for percision.
func (pv *Player) UpdateStateSync() {
	pv.State, _ = pv.player.State(pv.session)
}

// Updates state continuously and asyncchronously, runs reauthentication if requried.
func (pv *Player) UpdateStateLoop(session *auth.Session, config *config.Config) {
	go func() {
		var err error

		pv.State, err = pv.player.State(pv.session)
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
	Album         string
	ShuffleState  bool
	RepeatState   string
}

// Renders the player details as a content string.
func (pd *PlayerDetails) Content(track *spotify.Track, progressMs int, state *player.State) comp.Content {
	pd.Update(progressMs, state)

	repeat := func() string {
		if pd.RepeatState == DISABLED {
			return DISABLED
		}
		return ENABLED
	}()

	shuffle := func() string {
		if pd.ShuffleState {
			return ENABLED
		}
		return DISABLED
	}()

	timer := fmt.Sprintf("%sm:%ss / %sm:%ss", pd.ProgressMin, pd.ProgressSec, pd.DurationMin, pd.DurationSec)

	options := fmt.Sprintf("Vol: %s%%  Sfl: %v  Rep: %v", pd.VolumePercent, shuffle, repeat)

	return comp.Join([]string{
		color.HiGreenString("Track:     ") + pd.Track,
		color.HiGreenString("Artist:    ") + pd.Artists,
		color.HiGreenString("Album:     ") + pd.Album,
		color.HiGreenString("Progress:  ") + timer,
		color.HiGreenString("Options:   ") + options,
	}, "\n\n")
}

// Updates PlayerDetails with information in given state and track.
func (pd *PlayerDetails) Update(progressMs int, state *player.State) {
	if state == nil || state.Device == nil || state.Track == nil {
		return
	}

	pd.Track = state.Track.Name
	pd.Album = state.Track.Album.Name
	pd.ProgressSec = strconv.Itoa(((progressMs / 1000) % 60))
	pd.ProgressMin = strconv.Itoa((progressMs / 1000) / 60)
	pd.DurationSec = strconv.Itoa((state.Track.DurationMs / 1000) % 60)
	pd.DurationMin = strconv.Itoa((state.Track.DurationMs / 1000) / 60)
	pd.VolumePercent = strconv.Itoa(state.Device.VolumePercent)
	pd.RepeatState = state.RepeatState
	pd.ShuffleState = state.ShuffleState
	pd.Artists = state.Track.ArtistsString()

	for _, time := range []*string{&pd.ProgressSec, &pd.DurationSec} {
		if len(*time) == 1 {
			*time = "0" + *time
		}
	}
}

// The title status bar indicating whether the player is
// playing, paused or an invalid device is selected.
type statusBar struct {
	Style struct {
		NowPlaying lg.Style
		Paused     lg.Style
		NoPlayer   lg.Style
	}
	Status       string
	CurrentStyle lg.Style
}

// Renders the status bar as a string.
func (sb *statusBar) Render() string {
	return sb.CurrentStyle.Render(sb.Status)
}

// Renders the status bar as a content string.
func (sb *statusBar) Content() comp.Content {
	return comp.Content(sb.CurrentStyle.Render(sb.Status))
}

// Updates the status bar given the player's state.
func (sb *statusBar) Update(state *player.State) {
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
