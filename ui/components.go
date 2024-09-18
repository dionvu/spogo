package ui

import (
	"fmt"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

var ASCII_FLAGS_NORMAL aic_package.Flags = func() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Colored = true
	flags.Dimensions = []int{40, 20}
	flags.Braille = true
	flags.Threshold = 20

	return flags
}()

var ASCII_FLAGS_SMALL aic_package.Flags = func() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Colored = true
	flags.Dimensions = []int{20, 10}
	flags.Braille = true
	flags.Threshold = 20

	return flags
}()

var ASCII_FLAGS_TINY aic_package.Flags = func() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Colored = true
	flags.Dimensions = []int{16, 8}
	flags.Braille = true

	return flags
}()

var MainControlsView = func(view string) string {
	if view == PLAYER_VIEW {
		return padLines(CommonStyle.MainControls.Selected.Render("[ ")+
			CommonStyle.MainControls.Selected.Render("F1 Player")+
			CommonStyle.MainControls.Normal.Render(" | F2 Playlists | F3 Search | F4 Devices | F5 Help ]"), 4)
	}
	if view == PLAYLIST_VIEW {
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | ")+
			CommonStyle.MainControls.Selected.Render("F2 Playlists")+
			CommonStyle.MainControls.Normal.Render(" | F3 Search | F4 Devices | F5 Help ]"), TAB_WIDTH)
	}

	if view == HELP_VIEW {
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | F4 Devices ")+
			CommonStyle.MainControls.Selected.Render("| F5 Help ]"), TAB_WIDTH)
	}

	return "Unknown View"
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

	return padLines(fmt.Sprintf(
		"%s\n\n%s\n\n[%sm:%ss / %sm:%ss]\n\nvol: %v%%",
		track,
		artist,
		progressMin,
		progressSec,
		durationMin,
		durationSec,
		pv.State.Device.VolumePercent,
	), TAB_WIDTH)
}

var AsciiView = func(filepath string, flags aic_package.Flags) string {
	ascii, err := aic_package.Convert(filepath, flags)
	if err != nil {
		return "invalid ascii filepath or flags"
	}

	ascii = padLines(ascii, TAB_WIDTH)

	return ascii
}