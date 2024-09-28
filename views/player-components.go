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

func (pd *PlayerDetails) Render(track *spotify.Track, volumePercent int, shuffleState bool, progressMs int, terminal Terminal) string {
	pd.Update(track, progressMs, volumePercent)

	var shuffle string

	if shuffleState {
		shuffle = "on"
	} else {
		shuffle = "off"
	}

	line1 := CenterString(
		fmt.Sprintf("%s - %s", pd.Track, pd.Artists),
		terminal, -1)

	line2 := CenterString(fmt.Sprintf("%sm:%ss / %sm:%ss",
		pd.ProgressMin, pd.ProgressSec, pd.DurationMin, pd.DurationSec),
		terminal, -1)

	line3 := CenterString(fmt.Sprintf("vol: %s%% sfl: %v", pd.VolumePercent, shuffle),
		terminal, -1)

	return line1 + "\n\n" + line2 + "\n\n" + line3
}

type StatusBar struct {
	// The title status bar indicating, playing, paused or invalid device.
	Status string
	Style  *lg.Style
}

func (sb *StatusBar) Render(terminal Terminal) string {
	return CenterString(sb.Style.Render(sb.Status), terminal)
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
