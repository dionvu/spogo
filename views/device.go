package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/auth"
	comp "github.com/dionvu/spogo/components"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/player"
)

const MAX_DEVICE_ITEM_WIDTH = comp.DEFAULT_WIDTH - 4

type Device struct {
	Session   *auth.Session
	ListModel DeviceListModel
	deviceMap map[string]*player.Device
	itemMap   map[list.Item]string
}

// Creates a new device view with a list model for the
// user to select available playback devices.
func NewDeviceView(s *auth.Session) *Device {
	return &Device{
		Session: s,
	}
}

func (dv *Device) UpdateDevices() {
	items := []list.Item{}
	dv.deviceMap = map[string]*player.Device{}
	dv.itemMap = map[list.Item]string{}

	devices, _ := player.GetDevices(dv.Session)

	if devices != nil {
		for _, device := range *devices {
			item := comp.ListItem(comp.Content(device.Name).AdjustFit(MAX_DEVICE_ITEM_WIDTH))
			items = append(items, item)
			dv.deviceMap[device.Name] = &device
			dv.itemMap[item] = device.Name
		}
	}

	dv.ListModel = DeviceListModel{list: comp.NewDefaultList(items, "Devices")}
}

func (dv *Device) View(term comp.Terminal, device *player.Device, config *config.Config) string {
	var currDeviceInfo string

	if device == nil {
		currDeviceInfo = comp.Content("\nCurrent Selected Device: " + "none" + "\n\n" + "Type: " + "none").String()
	} else {
		currDeviceInfo = comp.Content("\nCurrent Selected Device: " + device.Name + "\n\n" + "Type: " + device.Type).String()
	}

	return comp.Join([]string{
		comp.InvisibleBarV(5).String(),
		dv.ListModel.View(),
		currDeviceInfo,
		"\n\n",
		ViewStatus{CurrentView: DEVICE_VIEW}.Content().Prepend('\n', 1).String(),
	}).CenterVertical(term).CenterHorizontal(term).String()
}

func (dv *Device) GetDeviceFromChoice(choice string) *player.Device {
	return dv.deviceMap[choice]
}

func (dv *Device) GetSelectedDevice() *player.Device {
	return dv.deviceMap[dv.itemMap[dv.ListModel.list.SelectedItem()]]
}

type DeviceListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m DeviceListModel) Update(msg tea.Msg) (DeviceListModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			return m, nil
		}
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (dlm DeviceListModel) View() string {
	return dlm.list.View()
}

func (dlm DeviceListModel) Content() comp.Content {
	return comp.Content(dlm.list.View())
}

func (_ DeviceListModel) Init() tea.Cmd {
	return nil
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
			style.Normal.Render(" | F2 Playlists | F3 Search | F4 Devices | F5 Help ]"),
		}, "")

	case PLAYLIST_VIEW:
		return comp.Join([]string{
			style.Normal.Render("[ F1 Player | "),
			style.Selected.Render("F2 Playlists"),
			style.Normal.Render(" | F3 Search | F4 Devices | F5 Help ]"),
		}, "")

	case HELP_VIEW:
		return comp.Join([]string{
			style.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | F4 Devices "),
			style.Selected.Render("| F5 Help ]"),
		}, "")

	case SEARCH_VIEW_QUERY, SEARCH_VIEW_TYPE, SEARCH_VIEW_RESULTS:
		return comp.Join([]string{
			style.Normal.Render("[ F1 Player | F2 Playlists | "),
			style.Selected.Render("F3 Search"),
			style.Normal.Render(" | F4 Devices | F5 Help ]"),
		}, "")

	case DEVICE_VIEW:
		return comp.Join([]string{
			style.Normal.Render("[ F1 Player | F2 Playlists | F3 Search | "),
			style.Selected.Render("F4 Devices"),
			style.Normal.Render(" | F5 Help ]"),
		}, "")

	default:
		return "Unknown View"
	}
}
