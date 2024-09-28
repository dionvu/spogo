package ui

import (
	"github.com/charmbracelet/lipgloss"
	lg "github.com/charmbracelet/lipgloss"
)

func MainControlsRender(view string) string {
	switch view {
	case PLAYER_VIEW:
		return padLines(CommonStyle.MainControls.Selected.Render("[ ")+
			CommonStyle.MainControls.Selected.Render("F1 Player")+
			CommonStyle.MainControls.Normal.Render(" | F2 Playlists | F3 Search | F4 Devices | F5 Help ]"), 4)

	case PLAYLIST_VIEW:
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | ")+
			CommonStyle.MainControls.Selected.Render("F2 Playlists")+
			CommonStyle.MainControls.Normal.Render(" | F3 Search | F4 Devices | F5 Help ]"), TAB_WIDTH)

	case HELP_VIEW:
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | F4 Devices ")+
			CommonStyle.MainControls.Selected.Render("| F5 Help ]"), TAB_WIDTH)

	case SEARCH_TYPE_VIEW:
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | ")+
			CommonStyle.MainControls.Selected.Render("F3 Search")+
			CommonStyle.MainControls.Normal.Render(" | F4 Devices | F5 Help ]"), TAB_WIDTH)

	case DEVICE_VIEW:
		return padLines(CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | ")+
			CommonStyle.MainControls.Selected.Render("F4 Devices")+
			CommonStyle.MainControls.Normal.Render(" | F5 Help ]"), TAB_WIDTH)

	default:
		return "Unknown View"
	}
}

func HelpString() string {
	h1 := lg.NewStyle().Underline(true).Foreground(lipgloss.Color("#458588"))
	h2 := lg.NewStyle().Underline(true)

	return h1.Render("CONTROLS") +
		"\n\n" + h2.Render("Global") +
		"\n\n" +
		"r - refresh, fixes visual issues" +
		"\n" +
		"a - select track from album of current playing track" +
		"\n" +
		"s - toggle shuffling of album/playlist" +
		"\n\n" + h2.Render("Player") +
		"\n\n" +
		"space - play / pause" +
		"\n" +
		"n - next track" +
		"\n" +
		"p - previous track" +
		"\n" +
		"[ - volume down" +
		"\n" +
		"] - volume up"
}
