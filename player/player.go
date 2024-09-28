package player

import (
	"encoding/json"
	"os"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/joomcode/errorx"
)

type Player struct {
	device *Device
}

// Creates a new player, getting device to any cached device
// or nil if no devices were found in cache.
func New(c *config.Config) (*Player, error) {
	p := &Player{
		device: nil,
	}

	if !deviceCacheExist(c) {
		if err := createCache(c); err != nil {
			err = err
			errors.LogError(err)
			return nil, err
		}
	}

	d, err := getCachedPlaybackDevice(c)
	if errorx.GetTypeName(err) == errors.JSONDecode.String() {
		return p, nil
	}
	if err != nil {
		return nil, err
	}

	p.device = d

	return p, nil
}

// Sets the playback device for player and saves playback device into
// cache file "device.json".
func (p *Player) SetDevice(d *Device, c *config.Config) error {
	p.device = d

	f, err := os.Create(c.DeviceFile())
	if err != nil {
		err = errors.FileCreate.Wrap(err, "failed to open device cache file")
		errors.LogError(err)
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(d)
	if err != nil {
		err = errors.JSONEncode.Wrap(err, "failed to marshal device")
		errors.LogError(err)
		return err
	}

	return nil
}

func (p *Player) Device() *Device {
	return p.device
}
