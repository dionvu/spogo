package player

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/session"
	"github.com/dionv/spogo/spotify/headers"
	"github.com/dionv/spogo/spotify/statuses"
	"github.com/dionv/spogo/spotify/urls"
)

// Resume uses the "transfer playback device" endpoint instead of the
// "resume playback" to ensure playback is always transfered to
// selected device before the players resumes playback.
func (p *Player) Resume(s *session.Session) error {
	if p.device == nil {
		return errors.DeviceError.New("no selected playback device")
	}

	data := map[string]interface{}{}

	data["device_ids"] = []string{p.device.ID}
	data["play"] = true

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

	if res.StatusCode == 404 {
		return errors.DeviceError.New("playback device is not active")
	}

	if res.StatusCode >= 400 {
		return errors.HTTPError.New("bad request")
	}

	return nil
}

// Skips the the next track in the queue.
func (p *Player) SkipNext(s *session.Session) error {
	if p.device == nil {
		return errors.DeviceError.New("no selected playback device")
	}

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
		return errors.HTTPError.New("bad request")
	}

	return nil
}

func (p *Player) SkipPrev(s *session.Session) error {
	if p.device == nil {
		return errors.DeviceError.New("no selected playback device")
	}

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
		return errors.HTTPError.New("bad request")
	}

	return nil
}

// Pauses playback on the current device.
func (p *Player) Pause(s *session.Session) error {
	if p.device == nil {
		return errors.DeviceError.New("no selected playback device")
	}

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
		return errors.HTTPError.New("bad request")
	}

	return nil
}

// Seeks to given position in milliseconds to user's current playing track.
func (p *Player) SeekToPosition(s *session.Session, pos int) error {
	if p.device == nil {
		return errors.DeviceError.New("no selected playback device")
	}

	query := url.Values{}
	query.Set("position_ms", strconv.Itoa(pos))
	query.Set("device_id", p.device.ID)

	req, err := http.NewRequest(http.MethodPut, urls.PLAYERSEEK+"?"+query.Encode(), nil)
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

	if res.StatusCode >= 400 {
		return errors.HTTPError.New("bad request")
	}

	return nil
}
