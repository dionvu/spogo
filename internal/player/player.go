package player

import (
	"encoding/json"
	"os"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/device"
)

type Player struct {
	device *device.Device
}

// Creates a new player, getting device to any cached device
// or nil if no devices were found in cache.
func New(c *config.Config) (*Player, error) {
	d, err := getCachedPlaybackDevice(c)
	if err != nil {
		return nil, err
	}

	p := &Player{
		device: d,
	}

	return p, nil
}

// Sets the player device and caches the device.
func (p *Player) SetDevice(d *device.Device, c *config.Config) error {
	p.device = d
	return cachePlaybackDevice(d, c)
}

func (p *Player) GetDevice() *device.Device {
	return p.device
}

func cachePlaybackDevice(d *device.Device, c *config.Config) error {
	f, err := os.Create(c.DeviceFile())
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to open device cache file")
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(d)
	if err != nil {
		errors.FileError.Wrap(err, "Failed to marshal device")
	}

	return nil
}

func getCachedPlaybackDevice(c *config.Config) (*device.Device, error) {
	f, err := os.Open(c.DeviceFile())
	if err != nil {
		return nil, errors.FileError.Wrap(err, "Failed to open device cache file")
	}
	defer f.Close()

	d := &device.Device{}

	err = json.NewDecoder(f).Decode(d)
	if err != nil {
		return nil, errors.FileError.Wrap(err, "Failed to marshal device")
	}

	return d, nil
}
