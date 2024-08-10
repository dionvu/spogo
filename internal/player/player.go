package player

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/device"
	"github.com/dionv/spogo/internal/session"
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

	if !cacheExist(c) {
		if err := createCache(c); err != nil {
			return nil, err
		}
	}

	d, err := getCachedPlaybackDevice(c)
	if errorx.GetTypeName(err) == errors.PLAYBACKERROR.String() {
		return p, nil
	}
	if err != nil {
		return nil, err
	}

	p.device = d

	return p, nil
}

// Prompts the user with all avaiable playback devices,
// caches user's choice, and sets player device to choice.
func (p *Player) UserSelectDevice(s *session.Session, c *config.Config) error {
	devices, err := device.GetDevices(s)
	if err != nil {
		return err
	}

	if len(*devices) == 0 {
		return errors.PLAYBACKERROR.New("No playback devices were found.")
	}

	deviceNames := []string{}

	for _, d := range *devices {
		deviceNames = append(deviceNames, d.Name)
	}

	prompt := promptui.Select{
		Label: "Select a playback device",
		Items: deviceNames,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return errors.PROMPTTUIERROR.Wrap(err, "Prompt failed")
	}

	err = p.SetDevice(&(*devices)[i], c)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func cacheExist(c *config.Config) bool {
	if _, err := os.ReadFile(filepath.Join(c.CachePath(), config.DEVICEFILE)); err != nil {
		return false
	}
	return true
}

func createCache(c *config.Config) error {
	file, err := os.Create(filepath.Join(c.CachePath(), config.DEVICEFILE))
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Creating file %v", c.FilePath()))
	}
	file.Close()

	return nil
}

// Sets the player device and caches the device.
func (p *Player) SetDevice(d *device.Device, c *config.Config) error {
	p.device = d
	return cachePlaybackDevice(d, c)
}

func (p *Player) GetDevice() *device.Device {
	if p.device == nil {
		return nil
	}

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

// Gets a playback device stored in device cache file.
// This function will error if no device is found.
func getCachedPlaybackDevice(c *config.Config) (*device.Device, error) {
	f, err := os.Open(c.DeviceFile())
	if err != nil {
		return nil, errors.FileError.Wrap(err, "Failed to open device cache file")
	}
	defer f.Close()

	d := &device.Device{}

	err = json.NewDecoder(f).Decode(d)
	if err != nil {

		if err == io.EOF {
			return nil, errors.PLAYBACKERROR.Wrap(err, "No device found.")
		}

		return nil, errors.FileError.Wrap(err, "Failed to marshal device")
	}

	return d, nil
}
