package ui

import (
	"fmt"
	"strconv"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify"
)

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
		errors.LogError(errors.PlayerViewInvalidState.New("invalid state passed, cannot update player details"))
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

// Renders the player details as a string.
func (pd *PlayerDetails) Render(track *spotify.Track, progressMs int, state *player.State) string {
	pd.Update(progressMs, state)

	title := Content(fmt.Sprintf("%s - %s", pd.Track, pd.Artists))

	timer := Content(fmt.Sprintf("%sm:%ss / %sm:%ss", pd.ProgressMin, pd.ProgressSec, pd.DurationMin, pd.DurationSec))

	options := Content(fmt.Sprintf("vol: %s%% sfl: %v", pd.VolumePercent, state.ShuffleState))

	return Join([]Content{title, timer, options}, "\n\n").String()
}

// Renders the player details as a content string.
func (pd *PlayerDetails) Content(track *spotify.Track, progressMs int, state *player.State) Content {
	pd.Update(progressMs, state)

	title := Content(fmt.Sprintf("%s - %s", pd.Track, pd.Artists))

	timer := Content(fmt.Sprintf("%sm:%ss / %sm:%ss", pd.ProgressMin, pd.ProgressSec, pd.DurationMin, pd.DurationSec))

	options := Content(fmt.Sprintf("vol: %s%% sfl: %v", pd.VolumePercent, state.ShuffleState))

	return Join([]Content{title, timer, options}, "\n\n")
}

// The title status bar indicating whether the player is
// playing, paused or an invalid device is selected.
type StatusBar struct {
	Status string
	Style  *lg.Style
}

// Renders the status bar as a string.
func (sb *StatusBar) Render() string {
	return sb.Style.Render(sb.Status)
}

// Renders the status bar as a content string.
func (sb *StatusBar) Content() Content {
	return Content(sb.Style.Render(sb.Status))
}

// Updates the status bar given the player's state.
func (sb *StatusBar) Update(state *player.State) {
	const (
		PAUSED      = "Paused"
		NO_PLAYER   = "Player Inactive"
		NOW_PLAYING = "Now Playing"
	)

	if state != nil && state.IsPlaying {
		sb.Style = &PlayerViewStyle.StatusBar.NowPlaying
		sb.Status = NOW_PLAYING
	} else if state != nil && !state.IsPlaying {
		sb.Style = &PlayerViewStyle.StatusBar.Paused
		sb.Status = PAUSED
	} else {
		sb.Style = &PlayerViewStyle.StatusBar.NoPlayer
		sb.Status = NO_PLAYER
	}
}

type ViewStatus struct {
	CurrentView string
}

// Updates the current view.
func (vs *ViewStatus) Update(view string) {
	vs.CurrentView = view
}

// Renders the ViewStatus as a content string based on the
// it's current view.
func (vs *ViewStatus) Content() Content {
	switch vs.CurrentView {
	case PLAYER_VIEW:
		return Join([]string{
			CommonStyle.MainControls.Selected.Render("[ "),
			CommonStyle.MainControls.Selected.Render("F1 Player"),
			CommonStyle.MainControls.Normal.Render(" | F2 Playlists | F3 Search | F4 Devices | F5 Help ]"),
		}, "")

	case PLAYLIST_VIEW:
		return Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | "),
			CommonStyle.MainControls.Selected.Render("F2 Playlists"),
			CommonStyle.MainControls.Normal.Render(" | F3 Search | F4 Devices | F5 Help ]"),
		}, "")

	case HELP_VIEW:
		return Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | F4 Devices "),
			CommonStyle.MainControls.Selected.Render("| F5 Help ]"),
		}, "")

	case SEARCH_TYPE_VIEW:
		return Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | "),
			CommonStyle.MainControls.Selected.Render("F3 Search"),
			CommonStyle.MainControls.Normal.Render(" | F4 Devices | F5 Help ]"),
		}, "")

	case DEVICE_VIEW:
		return Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | "),
			CommonStyle.MainControls.Selected.Render("F4 Devices"),
			CommonStyle.MainControls.Normal.Render(" | F5 Help ]"),
		}, "")

	default:
		return "Unknown View"
	}
}
