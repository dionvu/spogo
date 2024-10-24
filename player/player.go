package player

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/urls"
	"github.com/dionvu/spogo/spotify/auth"
	"github.com/joomcode/errorx"
)

type Player struct {
	device *Device
}

// Creates a new player, getting device to any cached device
// or nil if no devices were found in cache.
func New(c *config.Config) (*Player, error) {
	p := &Player{
		device: nil,
	}

	if !deviceCacheExist(c) {
		if err := createCache(c); err != nil {
			errors.Log(err)
			return nil, err
		}
	}

	d, err := getCachedPlaybackDevice(c)
	if errorx.GetTypeName(err) == errors.JSONDecode.String() {
		return p, nil
	}
	if err != nil {
		return nil, err
	}

	p.device = d

	return p, nil
}

// Sets the playback device for player and saves playback device into
// cache file "device.json".
func (p *Player) SetDevice(d *Device, c *config.Config) error {
	p.device = d

	f, err := os.Create(c.DeviceFile())
	if err != nil {
		err = errors.FileCreate.Wrap(err, "failed to open device cache file")
		errors.Log(err)
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(d)
	if err != nil {
		err = errors.JSONEncode.Wrap(err, "failed to marshal device")
		errors.Log(err)
		return err
	}

	return nil
}

func (p *Player) Device() *Device {
	return p.device
}

const (
	UPDATE_RATE_SEC          = time.Second
	POLLING_RATE_STATE_SEC   = time.Second * 5
	VOLUME_INCREMENT_PERCENT = 5
)

// ContextUri can be the uri of an album or playlist. Uri should be a track
// contained in the album or playlist.
func (p *Player) Play(contextUri string, uri string, s *auth.Session) error {
	if p.device == nil {
		err := errors.NoDevice.New("no selected playback device")
		errors.Log(err)
		return err
	}

	var payload interface{}

	if contextUri != "" && uri != "" {
		payload = struct {
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
	} else if uri == "" {
		payload = struct {
			Context_uri string `json:"context_uri"`
		}{
			Context_uri: contextUri,
		}
	} else {
		payload = struct {
			Uris []string `json:"uris"`
		}{
			Uris: []string{uri},
		}
	}

	j, err := json.Marshal(payload)
	if err != nil {
		err = errors.JSONMarshal.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	req, err := http.NewRequest(http.MethodPut, spotifyurls.PLAYERPLAY, strings.NewReader(string(j)))
	if err != nil {
		err = errors.HTTPRequest.Wrap(err, "failed to make new request to play: %v", uri)
		errors.Log(err)
		return err
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYERPLAY, res.StatusCode)

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return err
	}

	if res.StatusCode == 404 {
		err = errors.NoDevice.New("playback device is not active")
		errors.Log(err)
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return err
	}

	return nil
}

// Resume uses the "transfer playback device" endpoint instead of the
// "resume playback" to ensure playback is always transfered to
// selected device before the players resumes playback.
func (p *Player) Resume(s *auth.Session, play bool) error {
	if p.device == nil {
		err := errors.NoDevice.New("no selected playback device")
		errors.Log(err)
		return err
	}

	data := map[string]interface{}{}

	data["device_ids"] = []string{p.device.ID}
	data["play"] = play

	j, err := json.Marshal(data)
	if err != nil {
		err = errors.JSONMarshal.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	req, err := http.NewRequest(http.MethodPut, spotifyurls.PLAYER, strings.NewReader(string(j)))
	if err != nil {
		err = errors.HTTPRequest.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYER, res.StatusCode)

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return err
	}

	if res.StatusCode == http.StatusBadRequest {
		err = errors.NoDevice.New("playback device is not active")
		errors.Log(err)
		return err
	}

	if res.StatusCode >= http.StatusOK {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return err
	}

	return nil
}

// Skips the the next track in the queue.
func (p *Player) SkipNext(s *auth.Session) error {
	if p.device == nil {
		err := errors.NoDevice.New("no selected playback device")
		errors.Log(err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, spotifyurls.PLAYERNEXT, nil)
	if err != nil {
		err = errors.HTTPRequest.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYERNEXT, res.StatusCode)

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return err
	}

	if res.StatusCode >= http.StatusOK {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return err
	}

	return nil
}

func (p *Player) SkipPrev(s *auth.Session) error {
	if p.device == nil {
		err := errors.NoDevice.New("no selected playback device")
		errors.Log(err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, spotifyurls.PLAYERPREV, nil)
	if err != nil {
		err = errors.HTTPRequest.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYERPREV, res.StatusCode)

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return err
	}

	if res.StatusCode >= http.StatusOK {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return err
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
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYERPAUSE, res.StatusCode)

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return err
	}

	// Spotify returns 403 for some reason if track is already paused
	if res.StatusCode >= 403 {
		return nil
	}

	if res.StatusCode >= http.StatusOK {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		errors.LogApiCall(ResponseBody(res.Body), 1)
		return err
	}

	return nil
}

// Seeks to given position in milliseconds to user's current playing track.
func (p *Player) Seek(positionMs int, s *auth.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	query := url.Values{}
	query.Set("position_ms", strconv.Itoa(positionMs))
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

	errors.LogApiCall(spotifyurls.PLAYERSEEK, res.StatusCode)

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
		err = errors.HTTPRequest.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYERSHUFFLE, res.StatusCode)

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return err
	}

	if res.StatusCode >= http.StatusOK {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return err
	}

	return nil
}

// Toggles repeating on the current context.
func (p *Player) Repeat(state bool, s *auth.Session) error {
	if p.device == nil {
		return errors.NoDevice.New("no selected playback device")
	}

	query := &url.Values{}
	if state == true {
		query.Set("state", "context")
	} else {
		query.Set("state", "off")
	}

	url := spotifyurls.PLAYERREPEAT + "?" + query.Encode()
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		err = errors.HTTPRequest.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYERREPEAT, res.StatusCode)

	if res.StatusCode > 204 {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return err
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
		err = errors.HTTPRequest.Wrap(err, "failed to make request to change player volume")
		return err
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())
	req.Header.Set(headers.ContentType, headers.ApplicationJson)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return err
	}

	errors.LogApiCall(spotifyurls.PLAYERVOLUME, res.StatusCode)

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return err
	}

	if res.StatusCode >= http.StatusOK {
		err = errors.HTTP.New("Bad request, likely invalid player")
		errors.Log(err)
		return err
	}

	return nil
}

// Use to debug wacky api json confuzzling.
func PrintResponseBody(r io.ReadCloser) {
	b, _ := io.ReadAll(r)
	defer r.Close()
	fmt.Println(string(b))
}

// Use to debug wacky api json confuzzling.
func ResponseBody(r io.ReadCloser) string {
	b, _ := io.ReadAll(r)
	defer r.Close()
	return string(b)
}
