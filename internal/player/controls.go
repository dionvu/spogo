package player

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/api/headers"
	"github.com/dionv/spogo/internal/api/status"
	"github.com/dionv/spogo/internal/api/urls"
	"github.com/dionv/spogo/internal/session"
)

// Transfers playback to current device, then plays or pauses the device.
// Required to transfer playback when user first opens device. Thus, its
// an alternative that replaces the "player/play" and "player/pause"
// endpoints.
func (p *Player) TransferPlayback(s *session.Session, play bool) error {
	data := map[string]interface{}{}

	data["device_ids"] = []string{p.device.ID}
	data["play"] = play

	j, err := json.Marshal(data)
	if err != nil {
		return errors.JSONError.WrapWithNoMessage(err)
	}

	req, err := http.NewRequest(http.MethodPut, urls.PLAYER, strings.NewReader(string(j)))
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.AUTH, "Bearer "+s.AccessToken.String())

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

// Resumes playback on the current device.
func (p *Player) Resume(s *session.Session) error {
	return p.TransferPlayback(s, true)
}

func (p *Player) TogglePlayback(s *session.Session) error {
	return nil
}

// func (p *Player) Pause(s *session.Session) error {
// 	return p.TransferPlayback(s, false)
// }

// Skips the the next track in the queue.
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

	if res.StatusCode == status.BADTOKEN {
		return errors.ReauthenticationError.NewWithNoMessage()
	}

	if res.StatusCode >= 400 {
		return errors.HTTPError.New("Bad request, likely invalid player")
	}

	return nil
}

func (p *Player) SkipPrev(s *session.Session) error {
	req, err := http.NewRequest(http.MethodPost, urls.PLAYERPREV, nil)
	if err != nil {
		return errors.HTTPError.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.AUTH, "Bearer "+s.AccessToken.String())

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

// Pauses playback on the current device.
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

	if res.StatusCode == status.BADTOKEN {
		return errors.ReauthenticationError.NewWithNoMessage()
	}

	// Spotify returns 403 for some reason if track is already paused
	if res.StatusCode >= 403 {
		return nil
	}

	if res.StatusCode >= 400 {
		return errors.HTTPError.New("Bad request, likely invalid player")
	}

	return nil
}
