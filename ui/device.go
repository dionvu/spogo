package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/player"
)

type DeviceView struct {
	Session   *auth.Session
	ListModel *DeviceListModel
	deviceMap map[string]*player.Device
	itemMap   map[list.Item]string
}

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

func (dv *DeviceView) View(terminal Terminal, device *player.Device) string {
	mainControls := MainControlsView(DEVICE_VIEW)

	if terminal.Height < TERMINALSIZE.Small {
		return deviceView(dv, device)
	}

	return "\n\n" + mainControls + "\n\n" + deviceView(dv, device)
}

func (dv *DeviceView) GetDeviceFromChoice(choice string) *player.Device {
	return dv.deviceMap[choice]
}

func (dv *DeviceView) GetSelectedDevice() *player.Device {
	return dv.deviceMap[dv.itemMap[dv.ListModel.list.SelectedItem()]]
}
