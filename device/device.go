package device

import (
	"encoding/json"
	"net/http"

	"github.com/dionv/spogo/errors"
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
		return nil, errors.HTTPError.Wrap(err, "Failed to create http request for playback devices")
	}
	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTPError.Wrap(err, "Failed to get response for playback devices")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.ReauthenticationError.Wrap(err, "Bad token")
	}

	data := &struct {
		Devices []Device `json:"devices"`
	}{}

	if err = json.NewDecoder(res.Body).Decode(data); err != nil {
		return nil, errors.JSONError.Wrap(err, "Failed to decode json response for playback devices")
	}

	return &data.Devices, nil
}
