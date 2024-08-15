package player

import (
	"encoding/json"
	"net/http"

	"github.com/dionv/spogo/device"
	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/session"
	"github.com/dionv/spogo/spotify/headers"
	"github.com/dionv/spogo/spotify/urls"
)

type PlayerState struct {
	Device       *device.Device `json:"device"`
	ProgressMs   int            `json:"progress_ms"`
	IsPlaying    bool           `json:"is_playing"`
	ShuffleState bool           `json:"shuffle_state"`
}

func (p *Player) State(s *session.Session) (*PlayerState, error) {
	ps := &PlayerState{}

	req, err := http.NewRequest(http.MethodGet, urls.PLAYER, nil)
	if err != nil {
		return nil, errors.HTTPError.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.AUTH, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTPError.WrapWithNoMessage(err)
	}

	if res.StatusCode == 204 {
		return nil, errors.DeviceError.New("playback device is not active")
	}

	if res.StatusCode > 204 {
		return nil, errors.HTTPError.New("bad request")
	}

	if err := json.NewDecoder(res.Body).Decode(ps); err != nil {
		return nil, errors.JSONError.Wrap(err, "failed to decode player state response body")
	}
	defer res.Body.Close()

	return ps, nil
}
