package player

import (
	"encoding/json"
	"net/http"

	"github.com/dionv/spogo/device"
	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/session"
	"github.com/dionv/spogo/spotify"
	"github.com/dionv/spogo/spotify/api/headers"
	"github.com/dionv/spogo/spotify/api/urls"
)

type PlayerState struct {
	Device       *device.Device `json:"device"`
	ProgressMs   int            `json:"progress_ms"`
	IsPlaying    bool           `json:"is_playing"`
	ShuffleState bool           `json:"shuffle_state"`
	Item         interface{}    `json:"item"`
	Track        spotify.Track
	Episode      spotify.Episode
}

func (p *Player) State(s *session.Session) (*PlayerState, error) {
	ps := &PlayerState{}

	req, err := http.NewRequest(http.MethodGet, urls.PLAYER, nil)
	if err != nil {
		return nil, errors.HTTPRequest.Wrap(err, "failed to make new request for player state")
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == 204 {
		return nil, errors.NoDevice.New("playback device is not active")
	}

	if res.StatusCode > 204 {
		return nil, errors.HTTP.New("bad request")
	}

	if err := json.NewDecoder(res.Body).Decode(ps); err != nil {
		return nil, errors.JSONDecode.Wrap(err, "failed to decode player state response body")
	}
	defer res.Body.Close()

	itemMap, _ := ps.Item.(map[string]interface{})

	itemBytes, err := json.Marshal(itemMap)
	if err != nil {
		return nil, errors.JSONMarshal.Wrap(err, "failed to marshaling response: %v", itemMap)
	}

	// Type of item is a track
	if err := json.Unmarshal(itemBytes, &ps.Track); err == nil {
		return ps, nil
	}

	// Type of item is an episode
	if err := json.Unmarshal(itemBytes, &ps.Episode); err == nil {
		return ps, nil
	}

	return ps, errors.HTTP.New("response body is neither type track or episode")
}
