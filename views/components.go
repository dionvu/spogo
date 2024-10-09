package views

import (
	"github.com/charmbracelet/lipgloss"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/components"
)

const (
	PLAYER_VIEW           = "player_view"
	PLAYLIST_VIEW         = "playlist_view"
	PLAYLIST_TRACK_VIEW   = "playlist_track_view"
	ALBUM_TRACK_VIEW      = "album_track_view"
	REFRESH_VIEW          = "refresh_view"
	HELP_VIEW             = "help_view"
	TERMINAL_WARNING_VIEW = "terminal_warning_view"

	SEARCH_VIEW       = "search_view"
	SEARCH_TYPE_VIEW  = "search_type_view"
	SEARCH_QUERY_VIEW = "search_query_view"

	SEARCH_PLAYLIST_VIEW = "search_playlist_view"
	SEARCH_TRACK_VIEW    = "search_track_view"
	SEARCH_ALBUM_VIEW    = "search_album_view"

	SEARCH_RESULT_TRACK = "search_result_track"

	DEVICE_VIEW = "device_view"
)

type ViewStatus struct {
	CurrentView string
}

// Updates the current view.
func (vs *ViewStatus) Update(view string) {
	vs.CurrentView = view
}

// Renders the ViewStatus as a content string based on the
// it's current view.
func (vs *ViewStatus) Content() components.Content {
	switch vs.CurrentView {
	case PLAYER_VIEW:
		return components.Join([]string{
			CommonStyle.MainControls.Selected.Render("[ "),
			CommonStyle.MainControls.Selected.Render("F1 Player"),
			CommonStyle.MainControls.Normal.Render(" | F2 Playlists | F3 Search | F4 Devices | F5 Help ]"),
		}, "")

	case PLAYLIST_VIEW:
		return components.Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | "),
			CommonStyle.MainControls.Selected.Render("F2 Playlists"),
			CommonStyle.MainControls.Normal.Render(" | F3 Search | F4 Devices | F5 Help ]"),
		}, "")

	case HELP_VIEW:
		return components.Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | F4 Devices "),
			CommonStyle.MainControls.Selected.Render("| F5 Help ]"),
		}, "")

	case SEARCH_TYPE_VIEW:
		return components.Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | "),
			CommonStyle.MainControls.Selected.Render("F3 Search"),
			CommonStyle.MainControls.Normal.Render(" | F4 Devices | F5 Help ]"),
		}, "")

	case DEVICE_VIEW:
		return components.Join([]string{
			CommonStyle.MainControls.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | "),
			CommonStyle.MainControls.Selected.Render("F4 Devices"),
			CommonStyle.MainControls.Normal.Render(" | F5 Help ]"),
		}, "")

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
