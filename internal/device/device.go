package device

import (
	"encoding/json"
	"net/http"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/api/headers"
	"github.com/dionv/spogo/internal/api/urls"
	"github.com/dionv/spogo/internal/session"
)

type Device struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int    `json:"volume_percent"`
	SupportsVolume   bool   `json:"supports_volume"`
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

func GetDevices(s *session.Session) (*[]Device, error) {
	req, err := http.NewRequest(http.MethodGet, urls.PLAYERDEVICES, nil)
	if err != nil {
		return nil, errors.HTTPError.Wrap(err, "Failed to create http request for playback devices")
	}
	req.Header.Set(headers.AUTH, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTPError.Wrap(err, "Failed to get response for playback devices")
	}

	if res.StatusCode != 200 {
		return nil, errors.ReauthenticationError.Wrap(err, "Bad token")
	}

	data := &DevicesResponse{}

	if err = json.NewDecoder(res.Body).Decode(data); err != nil {
		return nil, errors.JSONError.Wrap(err, "Failed to decode json response for playback devices")
	}

	return &data.Devices, nil
}
