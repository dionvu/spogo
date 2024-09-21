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
}

func NewDeviceView(s *auth.Session, devices *[]player.Device) *DeviceView {
	items := []list.Item{}
	deviceMap := map[string]*player.Device{}

	for _, device := range *devices {
		items = append(items, Item(device.Name))
		deviceMap[device.Name] = &device
	}

	dv := DeviceView{
		Session:   s,
		ListModel: NewDeviceListModel(items),
		deviceMap: deviceMap,
	}

	return &dv
}

func (dv *DeviceView) View(playerView *PlayerView, terminal Terminal) string {
	mainControls := MainControlsView(SEARCH_TYPE_VIEW)

	if terminal.Height < TERMINALSIZE.Small {
		return "\n\n" + dv.ListModel.View()
	}

	return mainControls + "\n\n" + dv.ListModel.View()
}

func (dv *DeviceView) GetDeviceFromChoice(choice string) *player.Device {
	return dv.deviceMap[choice]
}
