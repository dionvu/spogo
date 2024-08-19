package player

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dionv/spogo/config"
	"github.com/dionv/spogo/device"
	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/session"
	"github.com/joomcode/errorx"
	"github.com/manifoldco/promptui"
)

type Player struct {
	device *device.Device
}

// Creates a new player, getting device to any cached device
// or nil if no devices were found in cache.
func New(c *config.Config) (*Player, error) {
	p := &Player{
		device: nil,
	}

	if !deviceCacheExist(c) {
		if err := createCache(c); err != nil {
			return nil, err
		}
	}

	d, err := getCachedPlaybackDevice(c)
	if errorx.GetTypeName(err) == errors.NoDevice.String() {
		return p, nil
	}
	if err != nil {
		return nil, err
	}

	p.device = d

	return p, nil
}

// Prompts the user with all avaiable playback devices and returns choice.
func (p *Player) UserSelectDevice(s *session.Session, c *config.Config, detailed bool) (*device.Device, error) {
	devices, err := device.GetDevices(s)
	if err != nil {
		return nil, err
	}

	if len(*devices) == 0 {
		return nil, errors.NoDevice.New("no active playback devices detected")
	}

	deviceNames := []string{}

	for _, d := range *devices {
		if detailed {
			deviceNames = append(deviceNames, d.StringDetailed())
		} else {
			deviceNames = append(deviceNames, d.String())
		}
	}

	prompt := promptui.Select{
		Label: "Select a playback device",
		Items: deviceNames,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, errors.PromptTui.Wrap(err, "devices prompt failed")
	}

	return &(*devices)[i], nil
}

// Sets the playback device for player and saves playback device into
// cache file "device.json".
func (p *Player) SetDevice(d *device.Device, c *config.Config) error {
	p.device = d

	f, err := os.Create(c.DeviceFile())
	if err != nil {
		return errors.FileCreate.Wrap(err, "failed to open device cache file")
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(d)
	if err != nil {
		errors.JSONEncode.Wrap(err, "failed to marshal device")
	}

	return nil
}

func (p *Player) GetDevice() *device.Device {
	if p.device == nil {
		return nil
	}

	return p.device
}

// Gets playback device stored in device cache file
// "device.json". This function will error if no device
// is found.
func getCachedPlaybackDevice(c *config.Config) (*device.Device, error) {
	d := &device.Device{}

	f, err := os.Open(c.DeviceFile())
	if err != nil {
		return nil, errors.FileOpen.Wrap(err, "failed to open device cache file")
	}
	defer f.Close()

	// Reached EOF before finished decoding into a device.
	if err = json.NewDecoder(f).Decode(d); err == io.EOF {
		return nil, errors.JSONDecode.Wrap(err, "invalid playback device")
	}

	if err != nil {
		return nil, errors.JSONMarshal.Wrap(err, "failed to marshal device")
	}

	return d, nil
}
