package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/player"
	"github.com/jedib0t/go-pretty/v6/table"
)

type DeviceView struct {
	Session   *auth.Session
	ListModel *DeviceListModel
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
func NewDeviceView(s *auth.Session) *DeviceView {
	items := []list.Item{}
	deviceMap := map[string]*player.Device{}
	itemMap := map[list.Item]string{}

	devices, _ := player.GetDevices(s)

	for _, device := range *devices {
		item := Item(device.Name)
		items = append(items, item)
		deviceMap[device.Name] = &device
		itemMap[item] = device.Name
	}

	dv := DeviceView{
		Session:   s,
		ListModel: NewDeviceListModel(items),
		deviceMap: deviceMap,
		itemMap:   itemMap,
	}

	if len((*devices)) > 0 {
		dv.ListModel.choice = (*devices)[0].Name
	}

	return &dv
}

func (dv *DeviceView) UpdateDevices() {
	items := []list.Item{}
	dv.deviceMap = map[string]*player.Device{}
	dv.itemMap = map[list.Item]string{}

	devices, _ := player.GetDevices(dv.Session)

	for _, device := range *devices {
		item := Item(device.Name)
		items = append(items, item)
		dv.deviceMap[device.Name] = &device
		dv.itemMap[item] = device.Name
	}

	dv.ListModel = NewDeviceListModel(items)
}

func (dv *DeviceView) View(terminal Terminal, device *player.Device) string {
	if terminal.IsSizeSmall() {
		return "\n\n" + RenderDeviceView(dv, device)
	}

	return "\n\n" + MainControlsRender(DEVICE_VIEW) + "\n\n" + RenderDeviceView(dv, device)
}

func (dv *DeviceView) GetDeviceFromChoice(choice string) *player.Device {
	return dv.deviceMap[choice]
}

func (dv *DeviceView) GetSelectedDevice() *player.Device {
	return dv.deviceMap[dv.itemMap[dv.ListModel.list.SelectedItem()]]
}

// Renders the list of devices, and current device information in a
// single row, two column table.
func RenderDeviceView(dv *DeviceView, device *player.Device) string {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	t.AppendRow(table.Row{
		dv.ListModel.View(),
		"Current Device: " + device.Name + "\n\n" + "Type: " + device.Type,
	})

	return t.Render()
}

func NewDeviceListModel(items []list.Item) *DeviceListModel {
	l := list.New(items, itemDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT)
	l.SetFilteringEnabled(false)
	l.Title = DeviceViewStyle.Title.Render("Devices")
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	lm := &DeviceListModel{list: l}

	return lm
}

func (m DeviceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

		return m, cmd
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
