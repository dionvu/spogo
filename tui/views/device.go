package views

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify/auth"
	comp "github.com/dionvu/spogo/tui/views/components"
)

type Device struct {
	Session    *auth.Session
	NumDevices int
}

func (dv *Device) UpdateNumberDevices() {
	devices, _ := player.GetDevices(dv.Session)
	if devices != nil {
		dv.NumDevices = len(*devices)
	}
}

func (dv *Device) View(term comp.Terminal, device *player.Device, config *config.Config) string {
	var currDeviceInfo string

	if device == nil {
		currDeviceInfo = comp.Content("Avaliable Devices: " + fmt.Sprint(dv.NumDevices) +
			"\n\nCtrl+D to select a device\n\nCurrent Device: " + "none").String()
	} else {
		currDeviceInfo = comp.Content("Avaliable Devices: " + fmt.Sprint(dv.NumDevices) +
			"\n\nCtrl+D to select a device\n\nCurrent Device: " + device.Name + " " + "(" +
			device.Type + ")").String()
	}

	return comp.Join([]string{
		comp.InvisibleBarV(10).String(),
		currDeviceInfo,
		comp.InvisibleBarV(7).String(),
		ViewStatus{CurrentView: DEVICE_VIEW}.Content().String(),
	}).CenterVertical(term).CenterHorizontal(term).String()
}

type ViewStatus struct {
	CurrentView string
}

// Renders the ViewStatus as a content string based on the
// it's current view.
func (vs ViewStatus) Content() comp.Content {
	style := struct {
		Selected lg.Style
		Normal   lg.Style
	}{
		Normal:   lg.NewStyle().Faint(true),
		Selected: lg.NewStyle(),
	}

	switch vs.CurrentView {
	case PLAYER_VIEW:
		return comp.Join([]string{
			style.Selected.Render("[ "),
			style.Selected.Render("F1 Player"),
			style.Normal.Render(" | F2 Playlists | F3 Search | F4 Help ]"),
		}, "")

	case PLAYLIST_VIEW:
		return comp.Join([]string{
			style.Normal.Render("[ F1 Player | "),
			style.Selected.Render("F2 Playlists"),
			style.Normal.Render(" | F3 Search | F4 Help ]"),
		}, "")

	case HELP_VIEW:
		return comp.Join([]string{
			style.Normal.Render("[ F1 Player | F2 Playlists | F3 Search "),
			style.Selected.Render("| F4 Help ]"),
		}, "")

	case SEARCH_VIEW_QUERY, SEARCH_VIEW_TYPE, SEARCH_VIEW_RESULTS:
		return comp.Join([]string{
			style.Normal.Render("[ F1 Player | F2 Playlists | "),
			style.Selected.Render("F3 Search"),
			style.Normal.Render(" | F4 Help ]"),
		}, "")

	// case DEVICE_VIEW:
	// 	return comp.Join([]string{
	// 		style.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | "),
	// 		style.Selected.Render("F4 Device"),
	// 		style.Normal.Render(" | F5 Help ]"),
	// 	}, "")

	default:
		return "Unknown View"
	}
}
