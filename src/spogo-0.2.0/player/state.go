package player

import (
	"encoding/json"
	"net/http"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/urls"
)

const (
	TRACK_TYPE   = "track"
	EPISODE_TYPE = "episode"
	AD_TYPE      = "ad"
	UNKNOWN_TYPE = "unknown"
)

type PlayerState struct {
	CurrentPlayingType string `json:"currently_playing_type"`

	Device       *Device `json:"device"`
	ProgressMs   int     `json:"progress_ms"`
	IsPlaying    bool    `json:"is_playing"`
	ShuffleState bool    `json:"shuffle_state"`

	// Each PlayerState will have either a nil track or episode,
	// depending on what the user is playing.
	Track   *spotify.Track
	Episode *spotify.Episode

	// Only used to retrieve the current item and
	// transfer into either a track or episode.
	Item interface{} `json:"item"`
}

func (p *Player) State(s *auth.Session) (*PlayerState, error) {
	ps := &PlayerState{}

	req, err := http.NewRequest(http.MethodGet, spotifyurls.PLAYER, nil)
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

	var track spotify.Track
	var episode spotify.Episode

	if err := json.Unmarshal(itemBytes, &track); err == nil {
		ps.Track = &track

		return ps, nil
	}

	if err := json.Unmarshal(itemBytes, &episode); err == nil {
		ps.Episode = &episode

		return ps, nil
	}

	return ps, errors.HTTP.New("response body is neither type track or episode")
}
