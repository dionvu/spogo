package player

import (
	"net/http"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/api/headers"
	"github.com/dionv/spogo/internal/api/urls"
	"github.com/dionv/spogo/internal/session"
)

func (p *Player) Resume(s *session.Session) error {
	req, err := http.NewRequest(http.MethodPut, urls.PLAYERPLAY, nil)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.AUTH, "Bearer "+s.AccessToken.String())
	req.Header.Set(headers.CONTENTTYPE, headers.JSON)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	if res.StatusCode != http.StatusOK {
		return errors.ReauthenticationError.WrapWithNoMessage(err)
	}

	return nil
}

func (p *Player) Pause(s *session.Session) error {
	req, err := http.NewRequest(http.MethodPut, urls.PLAYERPAUSE, nil)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.AUTH, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	if res.StatusCode != http.StatusOK {
		return errors.ReauthenticationError.WrapWithNoMessage(err)
	}

	return nil
}

func (p *Player) SkipNext(s *session.Session) error {
	req, err := http.NewRequest(http.MethodPost, urls.PLAYERNEXT, nil)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.AUTH, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}

	if res.StatusCode != http.StatusOK {
		return errors.ReauthenticationError.New("Bad access token")
	}

	return nil
}
