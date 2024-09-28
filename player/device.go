package player

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/urls"
)

// Device represents a spotify playback device. There is no
// guarantee that the device is still valid or active on the user's
// device since it will be cached.
type Device struct {
	ID            string `json:"id"`
	IsActive      bool   `json:"is_active"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	VolumePercent int    `json:"volume_percent"`
}

// Retrieves currently available playback devices, or an empty slice
// if none are available.
func GetDevices(s *auth.Session) (*[]Device, error) {
	req, err := http.NewRequest(http.MethodGet, spotifyurls.PLAYERDEVICES, nil)
	if err != nil {
		return nil, errors.HTTPRequest.Wrap(err, "failed to create http request for playback devices")
	}
	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTP.Wrap(err, "failed to get response for playback devices")
	}
	defer res.Body.Close()

	errors.LogApiCall(spotifyurls.PLAYERDEVICES, res.StatusCode)

	if res.StatusCode != 200 {
		return nil, errors.Reauthentication.Wrap(err, "bad token")
	}

	data := &struct {
		Devices []Device `json:"devices"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(data); err != nil {
		return nil, errors.JSONDecode.Wrap(err, "failed to decode json response for playback devices")
	}

	return &data.Devices, nil
}

// Helper function for player new function to
// see if the "device.json" cache file exists.
func deviceCacheExist(c *config.Config) bool {
	if _, err := os.ReadFile(filepath.Join(c.CachePath(), config.DEVICEFILE)); err != nil {
		return false
	}
	return true
}

// Helper function for the player new function
// to creates the "device.json" cache file.
func createCache(c *config.Config) error {
	file, err := os.Create(filepath.Join(c.CachePath(), config.DEVICEFILE))
	if err != nil {
		return errors.FileCreate.Wrap(err, fmt.Sprintf("creating file %v", c.FilePath()))
	}
	file.Close()

	return nil
}

// Gets playback device stored in device cache file
// "device.json". This function will error if no device
// is found.
func getCachedPlaybackDevice(c *config.Config) (*Device, error) {
	d := &Device{}

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
