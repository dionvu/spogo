package views

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/components"
	"github.com/dionvu/spogo/player"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Device struct {
	Session   *auth.Session
	ListModel DeviceListModel
	deviceMap map[string]*player.Device
	itemMap   map[list.Item]string
}

type DeviceListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

// Creates a new device view with a list model for the
// user to select available playback devices.
func NewDeviceView(s *auth.Session) *Device {
	items := []list.Item{}
	deviceMap := map[string]*player.Device{}
	itemMap := map[list.Item]string{}

	devices, _ := player.GetDevices(s)

	for _, device := range *devices {
		item := components.ListItem(device.Name)
		items = append(items, item)
		deviceMap[device.Name] = &device
		itemMap[item] = device.Name
	}

	dv := Device{
		Session:   s,
		deviceMap: deviceMap,
		itemMap:   itemMap,
	}

	if len((*devices)) > 0 {
		dv.ListModel.choice = (*devices)[0].Name
	}

	return &dv
}

func (dv *Device) UpdateDevices() {
	items := []list.Item{}
	dv.deviceMap = map[string]*player.Device{}
	dv.itemMap = map[list.Item]string{}

	devices, _ := player.GetDevices(dv.Session)

	for _, device := range *devices {
		item := components.ListItem(components.Content(device.Name).AdjustFit(components.DEFAULT_WIDTH - 4))
		items = append(items, item)
		dv.deviceMap[device.Name] = &device
		dv.itemMap[item] = device.Name
	}

	dv.ListModel = DeviceListModel{list: components.NewCustomList(items,
		"Devices", components.DEFAULT_WIDTH, components.LIST_HEIGHT_NORMAL)}
}

func (dv *Device) View(term components.Terminal, device *player.Device) string {
	var currDeviceInfo string

	if device == nil {
		currDeviceInfo = components.Content("Current Selected Device: "+"none"+"\n\n"+"Type: "+"none").Prepend('\n', 2).String()
	} else {
		currDeviceInfo = components.Content("Current Selected Device: "+device.Name+"\n\n"+"Type: "+device.Type).Prepend('\n', 2).String()
	}

	link := "https://i.pinimg.com/736x/7f/cb/55/7fcb55a037d93681c7396e50b6f074aa.jpg"

	cd, _ := os.UserCacheDir()
	img := components.Image{FilePath: filepath.Join(cd, "spogo", "temp500.jpeg")}
	img.Update(link)

	t := components.NewDefaultTable()
	t.AppendRow(table.Row{
		components.Content(dv.ListModel.View()).Prepend('\n', 1),
		img.AsciiSmall().Content().PadLinesLeft(0).PadLinesLeft(15),
	})
	vs := ViewStatus{}
	vs.Update(DEVICE_VIEW)

	return components.Join([]string{
		components.Content(t.Render()).Prepend('\n', 6).String(),
		currDeviceInfo,
		vs.Content().Prepend('\n', 2).String(),
	}, "\n").CenterVertical(term).CenterHorizontal(term).String()
}

func (dv *Device) GetDeviceFromChoice(choice string) *player.Device {
	return dv.deviceMap[choice]
}

func (dv *Device) GetSelectedDevice() *player.Device {
	return dv.deviceMap[dv.itemMap[dv.ListModel.list.SelectedItem()]]
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

func (_ DeviceListModel) Init() tea.Cmd {
	return nil
}
