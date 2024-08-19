package device

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/icons"
	"github.com/dionv/spogo/session"
	"github.com/dionv/spogo/spotify/api/headers"
	"github.com/dionv/spogo/spotify/api/urls"
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
func GetDevices(s *session.Session) (*[]Device, error) {
	req, err := http.NewRequest(http.MethodGet, urls.PLAYERDEVICES, nil)
	if err != nil {
		return nil, errors.HTTPRequest.Wrap(err, "failed to create http request for playback devices")
	}
	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTP.Wrap(err, "failed to get response for playback devices")
	}
	defer res.Body.Close()

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

func (d *Device) String() string {
	return d.Name
}

func (d *Device) StringDetailed() string {
	return fmt.Sprintf("%v, %v, %v %v", d.Name, d.Type, d.VolumePercent, icons.VolumeMax)
}
