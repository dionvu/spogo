package player

import (
	"encoding/json"
	"net/http"

	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/urls"
	"github.com/dionvu/spogo/spotify/auth"
)

const (
	TRACK_TYPE   = "track"
	EPISODE_TYPE = "episode"
	AD_TYPE      = "ad"
	UNKNOWN_TYPE = "unknown"
)

type State struct {
	CurrentPlayingType string `json:"currently_playing_type"`

	Device       *Device `json:"device"`
	ProgressMs   int     `json:"progress_ms"`
	IsPlaying    bool    `json:"is_playing"`
	ShuffleState bool    `json:"shuffle_state"`
	RepeatState  string  `json:"repeat_state"`

	Context struct {
		Type string `json:"type"`
	} `json:"context"`

	// Each State will have either a nil track or episode,
	// depending on what the user is playing.
	Track   *spotify.Track
	Episode *spotify.Episode

	// Only used to retrieve the current item and
	// transfer into either a track or episode.
	Item interface{} `json:"item"`
}

func (p *Player) State(s *auth.Session) (*State, error) {
	ps := &State{}

	req, err := http.NewRequest(http.MethodGet, spotifyurls.PLAYER, nil)
	if err != nil {
		return nil, errors.HTTPRequest.Wrap(err, "failed to make new request for player state")
	}

	req.Header.Set(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return nil, err
	}

	errors.LogApiCall(spotifyurls.PLAYER, res.StatusCode)

	if res.StatusCode == 204 {
		err = errors.NoDevice.New("playback device is not active")
		errors.Log(err)
		return nil, err
	}

	if res.StatusCode > 204 {
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return nil, err
	}

	if err := json.NewDecoder(res.Body).Decode(ps); err != nil {
		err = errors.JSONDecode.Wrap(err, "failed to decode player state response body")
		errors.Log(err)
		return nil, err
	}
	defer res.Body.Close()

	itemMap, _ := ps.Item.(map[string]interface{})

	itemBytes, err := json.Marshal(itemMap)
	if err != nil {
		err = errors.JSONMarshal.Wrap(err, "failed to marshaling response: %v", itemMap)
		errors.Log(err)
		return nil, err
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

	err = errors.HTTP.New("response body is neither type track or episode")
	errors.Log(err)
	return nil, err
}
