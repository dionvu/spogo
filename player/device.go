package player

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/urls"
	"github.com/dionvu/spogo/spotify/auth"
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
		err = errors.HTTPRequest.Wrap(err, "failed to create http request for playback devices")
		errors.Log(err)
		return nil, err
	}
	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.Wrap(err, "failed to get response for playback devices")
		errors.Log(err)
		return nil, err
	}
	defer res.Body.Close()

	errors.LogApiCall(spotifyurls.PLAYERDEVICES, res.StatusCode)

	if res.StatusCode != 200 {
		err = errors.Reauthentication.Wrap(err, "bad token")
		errors.Log(err)
		return nil, err
	}

	data := &struct {
		Devices []Device `json:"devices"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(data); err != nil {
		err = errors.JSONDecode.Wrap(err, "failed to decode json response for playback devices")
		errors.Log(err)
		return nil, err
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
		err = errors.FileCreate.Wrap(err, fmt.Sprintf("creating file %v", c.FilePath()))
		errors.Log(err)
		return err
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
		err = errors.FileOpen.Wrap(err, "failed to open device cache file")
		errors.Log(err)
		return nil, err
	}
	defer f.Close()

	// Reached EOF before finished decoding into a device.
	if err = json.NewDecoder(f).Decode(d); err == io.EOF {
		err = errors.JSONDecode.Wrap(err, "invalid playback device")
		errors.Log(err)
		return nil, err
	}

	if err != nil {
		err = errors.JSONMarshal.Wrap(err, "failed to marshal device")
		errors.Log(err)
		return nil, err
	}

	return d, nil
}

func (d Device) IsMobile() bool {
	return d.Type == "Smartphone" || d.Type == "Tablet"
}

func IsValidVolume(vol int) bool {
	return 0 <= vol && vol <= 100
}
