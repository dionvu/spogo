package player

import (
	"encoding/json"
	"net/http"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/api/headers"
	"github.com/dionv/spogo/internal/api/urls"
	"github.com/dionv/spogo/internal/device"
	"github.com/dionv/spogo/internal/session"
)

type PlayerState struct {
	Device     *device.Device `json:"device"`
	ProgressMs int            `json:"progress_ms"`
}

func (p *Player) GetPlayerState(s *session.Session) (*PlayerState, error) {
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
		return nil, errors.NoPlaybackError.New("playback not available or active")
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
