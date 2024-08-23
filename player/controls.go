package player

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/session"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/status"
	"github.com/dionvu/spogo/spotify/api/urls"
)

func (p *Player) Play(ctxUri string, uri string, s *session.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	data := map[string]interface{}{}

	data["device_id"] = p.device.ID
	data["position_ms"] = 0

	if ctxUri != "" {
		data["context_uri"] = ctxUri
	}

	if uri != "" {
		data["uris"] = []string{uri}
	}

	j, err := json.Marshal(data)
	if err != nil {
		return errors.JSONMarshal.WrapWithNoMessage(err)
	}

	req, err := http.NewRequest(http.MethodPut, urls.PLAYERPLAY, strings.NewReader(string(j)))
	if err != nil {
		return errors.HTTPRequest.Wrap(err, "failed to make new request to play: %v", uri)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BadToken {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode == 404 {
		return errors.NoDevice.New("playback device is not active")
	}

	if res.StatusCode >= 400 {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Resume uses the "transfer playback device" endpoint instead of the
// "resume playback" to ensure playback is always transfered to
// selected device before the players resumes playback.
func (p *Player) Resume(s *session.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	data := map[string]interface{}{}

	data["device_ids"] = []string{p.device.ID}
	data["play"] = true

	j, err := json.Marshal(data)
	if err != nil {
		return errors.JSONMarshal.WrapWithNoMessage(err)
	}

	req, err := http.NewRequest(http.MethodPut, urls.PLAYER, strings.NewReader(string(j)))
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

	if res.StatusCode == 404 {
		return errors.NoDevice.New("playback device is not active")
	}

	if res.StatusCode >= 400 {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Skips the the next track in the queue.
func (p *Player) SkipNext(s *session.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	req, err := http.NewRequest(http.MethodPost, urls.PLAYERNEXT, nil)
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

func (p *Player) SkipPrev(s *session.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	req, err := http.NewRequest(http.MethodPost, urls.PLAYERPREV, nil)
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

// Pauses playback on the current device.
func (p *Player) Pause(s *session.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	req, err := http.NewRequest(http.MethodPut, urls.PLAYERPAUSE, nil)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BadToken {
		return errors.Reauthentication.NewWithNoMessage()
	}

	// Spotify returns 403 for some reason if track is already paused
	if res.StatusCode >= 403 {
		return nil
	}

	if res.StatusCode >= 400 {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Seeks to given position in milliseconds to user's current playing track.
func (p *Player) SeekToPosition(s *session.Session, pos int) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	query := url.Values{}
	query.Set("position_ms", strconv.Itoa(pos))
	query.Set("device_id", p.device.ID)

	req, err := http.NewRequest(http.MethodPut, urls.PLAYERSEEK+"?"+query.Encode(), nil)
	if err != nil {
		return errors.HTTPRequest.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

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
