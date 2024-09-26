package player

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/urls"
	"github.com/dionvu/spogo/utils"
)

// ContextUri can be the uri of an album or playlist. Uri should be a track
// contained in the album or playlist.
func (p *Player) Play(contextUri string, uri string, s *auth.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	payload := struct {
		Context_uri string `json:"context_uri"`
		Offset      struct {
			Uri string `json:"uri"`
		} `json:"offset"`
	}{
		Context_uri: contextUri,
		Offset: struct {
			Uri string `json:"uri"`
		}{
			Uri: uri,
		},
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return errors.JSONMarshal.WrapWithNoMessage(err)
	}

	req, err := http.NewRequest(http.MethodPut, spotifyurls.PLAYERPLAY, strings.NewReader(string(j)))
	if err != nil {
		return errors.HTTPRequest.Wrap(err, "failed to make new request to play: %v", uri)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	utils.PrintResponseBody(res.Body)

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode == 404 {
		return errors.NoDevice.New("playback device is not active")
	}

	if res.StatusCode >= http.StatusBadRequest {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Resume uses the "transfer playback device" endpoint instead of the
// "resume playback" to ensure playback is always transfered to
// selected device before the players resumes playback.
func (p *Player) Resume(s *auth.Session, play bool) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	data := map[string]interface{}{}

	data["device_ids"] = []string{p.device.ID}
	data["play"] = play

	j, err := json.Marshal(data)
	if err != nil {
		return errors.JSONMarshal.WrapWithNoMessage(err)
	}

	req, err := http.NewRequest(http.MethodPut, spotifyurls.PLAYER, strings.NewReader(string(j)))
	if err != nil {
		return errors.HTTPRequest.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode == http.StatusBadRequest {
		return errors.NoDevice.New("playback device is not active")
	}

	if res.StatusCode >= http.StatusOK {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Skips the the next track in the queue.
func (p *Player) SkipNext(s *auth.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	req, err := http.NewRequest(http.MethodPost, spotifyurls.PLAYERNEXT, nil)
	if err != nil {
		return errors.HTTPRequest.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode >= http.StatusOK {
		return errors.HTTP.New("bad request")
	}

	return nil
}

func (p *Player) SkipPrev(s *auth.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	req, err := http.NewRequest(http.MethodPost, spotifyurls.PLAYERPREV, nil)
	if err != nil {
		return errors.HTTPRequest.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode >= http.StatusOK {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Pauses playback on the current device.
func (p *Player) Pause(s *auth.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	req, err := http.NewRequest(http.MethodPut, spotifyurls.PLAYERPAUSE, nil)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	// Spotify returns 403 for some reason if track is already paused
	if res.StatusCode >= 403 {
		return nil
	}

	if res.StatusCode >= http.StatusOK {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Seeks to given position in milliseconds to user's current playing track.
func (p *Player) SeekToPosition(s *auth.Session, pos int) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	query := url.Values{}
	query.Set("position_ms", strconv.Itoa(pos))
	query.Set("device_id", p.device.ID)

	req, err := http.NewRequest(http.MethodPut, spotifyurls.PLAYERSEEK+"?"+query.Encode(), nil)
	if err != nil {
		return errors.HTTPRequest.WrapWithNoMessage(err)
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode >= http.StatusOK {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Enables or disables shuffling of tracks in current playlist or album.
func (p *Player) Shuffle(state bool, s *auth.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	query := &url.Values{}
	query.Set("state", strconv.FormatBool(state))

	url := spotifyurls.PLAYERSHUFFLE + "?" + query.Encode()
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return errors.HTTPRequest.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode >= http.StatusOK {
		return errors.HTTP.New("bad request")
	}

	return nil
}

// Sets the player volume to a value between [0-100] percent.
func (p *Player) SetVolume(s *auth.Session, val int) error {
	val = min(max(0, val), 100)

	query := url.Values{}
	query.Set("volume_percent", strconv.Itoa(val))

	req, err := http.NewRequest(http.MethodPut, spotifyurls.PLAYERVOLUME+"?"+query.Encode(), nil)
	if err != nil {
		return errors.HTTPRequest.Wrap(err, "failed to make request to change player volume")
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())
	req.Header.Set(headers.ContentType, headers.ApplicationJson)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 401 {
		return errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode >= http.StatusOK {
		return errors.HTTP.New("Bad request, likely invalid player")
	}

	return nil
}
