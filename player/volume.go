package player

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/session"
	"github.com/dionv/spogo/spotify/headers"
	"github.com/dionv/spogo/spotify/statuses"
	"github.com/dionv/spogo/spotify/urls"
)

func (p *Player) SetVolume(s *session.Session, val int) error {
	query := url.Values{}
	query.Set("volume_percent", strconv.Itoa(val))

	req, err := http.NewRequest(http.MethodPut, urls.PLAYERVOLUME+"?"+query.Encode(), nil)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.AUTH, "Bearer "+s.AccessToken.String())
	req.Header.Set(headers.CONTENTTYPE, headers.JSON)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BADTOKEN {
		return errors.ReauthenticationError.NewWithNoMessage()
	}

	if res.StatusCode >= 400 {
		return errors.HTTPError.New("Bad request, likely invalid player")
	}

	return nil
}
