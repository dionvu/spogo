package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
	"github.com/dionvu/spogo/spotify/auth"
	comp "github.com/dionvu/spogo/tui/views/components"
)

type Device struct {
	Session    *auth.Session
	Cfg        *config.Config
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
		ViewStatus{CurrentView: DEVICE_VIEW}.Content(dv.Cfg).String(),
	}).CenterVertical(term).CenterHorizontal(term).String()
}

type ViewStatus struct {
	CurrentView string
}

// Renders the ViewStatus as a content string based on the
// it's current view.
func (vs ViewStatus) Content(cfg *config.Config) comp.Content {
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
			// "[ Spogo 󰝚 ] ",
			style.Selected.Render("[ "),
			style.Selected.Render("Player"),
			style.Normal.Render(" - Playlists - Search - F4 Help ]"),
		}, "")

	case PLAYLIST_VIEW:
		return comp.Join([]string{
			// "[ Spogo 󰝚 ] ",
			style.Normal.Render("[ Player - "),
			style.Selected.Render("Playlists"),
			style.Normal.Render(" - Search - F4 Help ]"),
		}, "")

	case HELP_VIEW:
		return comp.Join([]string{
			// "[ Spogo 󰝚 ] ",
			style.Normal.Render("[ Player - Playlists - Search "),
			style.Selected.Render("- F4 Help ]"),
		}, "")

	case SEARCH_VIEW_QUERY, SEARCH_VIEW_TYPE, SEARCH_VIEW_RESULTS:
		return comp.Join([]string{
			// "[ Spogo 󰝚 ] ",
			style.Normal.Render("[ Player - Playlists - "),
			style.Selected.Render("Search"),
			style.Normal.Render(" - F4 Help ]"),
		}, "")

	default:
		return "Unknown View"
	}
}

var (
	titleStyle = func() lg.Style {
		b := lg.RoundedBorder()
		b.Right = "├"
		return lg.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lg.Style {
		b := lg.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type Help struct {
	content  string
	viewport viewport.Model
}

func NewHelpView() Help {
	content := ""

	x, y := comp.GetTerminalSize()

	vp := viewport.New(int(float64(x)*0.5), int(float64(y)*0.5))

	return Help{content: string(content), viewport: vp}
}

func (m Help) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m Help) headerView() string {
	title := titleStyle.Render("Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lg.Width(title)))
	return lg.JoinHorizontal(lg.Center, title, line)
}

func (m Help) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lg.Width(info)))
	return lg.JoinHorizontal(lg.Center, line, info)
}

func (m Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// x, y := comp.GetTerminalSize()
	// m.viewport.Height = int(float64(x) * 0.5)
	// m.viewport.Width = int(float64(y) * 0.5)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

		// case tea.WindowSizeMsg:
		// 	headerHeight := lg.Height(m.headerView())
		// 	footerHeight := lg.Height(m.footerView())
		// 	verticalMarginHeight := headerHeight + footerHeight
		//
		// 	m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		// 	m.viewport.YPosition = headerHeight
		// 	m.viewport.SetContent(m.content)
		//
		// 	m.viewport.YPosition = headerHeight + 1
		// }
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Help) Init() tea.Cmd {
	return nil
}
