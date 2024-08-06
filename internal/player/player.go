package player

import (
	"github.com/dionv/spogo/internal/device"
)

type Player struct {
	device *device.Device
}

func New(d *device.Device) *Player {
	return &Player{
		device: d,
	}
}

func (p *Player) SetDevice(d *device.Device) {
	p.device = d
}

func (p *Player) GetDevice() *device.Device {
	return p.device
}
