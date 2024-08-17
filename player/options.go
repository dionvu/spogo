package player

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/session"
	"github.com/dionv/spogo/spotify/api/headers"
	"github.com/dionv/spogo/spotify/api/status"
	"github.com/dionv/spogo/spotify/api/urls"
)

// Enables or disables shuffling of tracks in current playlist or album.
func (p *Player) Shuffle(state bool, s *session.Session) error {
	if p.device == nil {
		return errors.DeviceError.New("no selected playback device")
	}

	query := &url.Values{}
	query.Set("state", strconv.FormatBool(state))

	url := urls.PLAYERSHUFFLE + "?" + query.Encode()
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BadToken {
		return errors.ReauthenticationError.NewWithNoMessage()
	}

	if res.StatusCode >= 400 {
		return errors.HTTPError.New("bad request")
	}

	return nil
}
