package spotify

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/session"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/status"
	"github.com/dionvu/spogo/spotify/api/urls"
)

type PlaylistsResponse struct {
	Href     string     `json:"href"`
	Limit    int        `json:"limit"`
	Next     string     `json:"next"`
	Offset   int        `json:"offset"`
	Previous string     `json:"previous"`
	Total    int        `json:"total"`
	Items    []Playlist `json:"items"`
}

type Playlist struct {
	Images       []Image `json:"images"`
	Description  string  `json:"description"`
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Total        int     `json:"total"`
	Public       bool    `json:"public"`
	Tracks       Tracks  `json:"tracks"`
	URI          string  `json:"uri"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Owner struct {
		Followers   Followers `json:"followers"`
		ID          string    `json:"id"`
		Type        string    `json:"type"`
		URI         string    `json:"uri"`
		DisplayName string    `json:"display_name"`
	} `json:"owner"`
}
type Tracks struct {
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
	Items    []Track `json:"items"`
}

type Followers struct {
	Total int `json:"total"`
}

func UserPlaylists(s *session.Session) (*[]Playlist, error) {
	req, err := http.NewRequest(http.MethodGet, urls.PLAYLISTS, nil)
	if err != nil {
		return nil, errors.HTTPRequest.Wrap(err, "failed to make request for playlists")
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTP.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BadToken {
		return nil, errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode >= 400 {
		return nil, errors.HTTP.New("bad request")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.HTTP.Wrap(err, "failed to read response body")
	}

	pr := &PlaylistsResponse{}

	err = json.Unmarshal(b, pr)
	if err != nil {
		return nil, errors.JSONUnmarshal.Wrap(err, "failed to unmarshal playlists response")
	}

	return &pr.Items, nil
}
