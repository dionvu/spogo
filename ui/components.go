package ui

import (
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

func MainControlsRender(view string) string {
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

	if view == SEARCH_TYPE_VIEW {
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | ")+
			CommonStyle.MainControls.Selected.Render("F3 Search")+
			CommonStyle.MainControls.Normal.Render(" | F4 Devices | F5 Help ]"), TAB_WIDTH)
	}

	if view == DEVICE_VIEW {
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | ")+
			CommonStyle.MainControls.Selected.Render("F4 Devices")+CommonStyle.MainControls.Normal.Render(" | F5 Help ]"), TAB_WIDTH)
	}

	return "Unknown View"
}

func AsciiRender(filepath string, flags aic_package.Flags) string {
	ascii, _ := aic_package.Convert(filepath, flags)

	ascii = ascii

	return ascii
}
