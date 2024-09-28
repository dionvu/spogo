package ui

import (
	"fmt"
	"strconv"

	lg "github.com/charmbracelet/lipgloss"
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

func (pd *PlayerDetails) Update(track *spotify.Track, progressMs int, volume int) {
	if track == nil {
		return
	}

	pd.Track = track.Name

	pd.Artists = ""

	for i := 0; i < len(track.Artists); i++ {
		pd.Artists += track.Artists[i].Name

		if i != len(track.Artists)-1 {
			pd.Artists += ", "
		}
	}

	pd.ProgressSec = strconv.Itoa(((progressMs / 1000) % 60))
	pd.ProgressMin = strconv.Itoa((progressMs / 1000) / 60)

	pd.DurationSec = strconv.Itoa((track.DurationMs / 1000) % 60)
	pd.DurationMin = strconv.Itoa((track.DurationMs / 1000) / 60)

	for _, time := range []*string{&pd.ProgressSec, &pd.DurationSec} {
		if len(*time) == 1 {
			*time = "0" + *time
		}
	}

	pd.VolumePercent = strconv.Itoa(volume)
}

// Renders the player details as a string.
func (pd *PlayerDetails) Render(track *spotify.Track, volume int, shuffleState bool, progressMs int) string {
	pd.Update(track, progressMs, volume)

	var shuffle string

	if shuffleState {
		shuffle = "on"
	} else {
		shuffle = "off"
	}

	title := Content(fmt.Sprintf("%s - %s", pd.Track, pd.Artists))

	timer := Content(fmt.Sprintf("%sm:%ss / %sm:%ss", pd.ProgressMin, pd.ProgressSec, pd.DurationMin, pd.DurationSec))

	options := Content(fmt.Sprintf("vol: %s%% sfl: %v", pd.VolumePercent, shuffle))

	return Join([]Content{title, timer, options}, "\n\n").String()
}

// Renders the player details as a content string.
func (pd *PlayerDetails) Content(track *spotify.Track, volume int, shuffleState bool, progressMs int) Content {
	pd.Update(track, progressMs, volume)

	var shuffle string

	if shuffleState {
		shuffle = "on"
	} else {
		shuffle = "off"
	}

	title := Content(fmt.Sprintf("%s - %s", pd.Track, pd.Artists))

	timer := Content(fmt.Sprintf("%sm:%ss / %sm:%ss", pd.ProgressMin, pd.ProgressSec, pd.DurationMin, pd.DurationSec))

	options := Content(fmt.Sprintf("vol: %s%% sfl: %v", pd.VolumePercent, shuffle))

	return Join([]Content{title, timer, options}, "\n\n")
}

type StatusBar struct {
	// The title status bar indicating, playing, paused or invalid device.
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

func (sb *StatusBar) Update(state *player.PlayerState) {
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
