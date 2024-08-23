package player

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/session"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/status"
	"github.com/dionvu/spogo/spotify/api/urls"
)

// Enables or disables shuffling of tracks in current playlist or album.
func (p *Player) Shuffle(state bool, s *session.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	query := &url.Values{}
	query.Set("state", strconv.FormatBool(state))

	url := urls.PLAYERSHUFFLE + "?" + query.Encode()
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return errors.HTTPRequest.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BadToken {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode >= 400 {
		return errors.HTTP.New("bad request")
	}

	return nil
}
