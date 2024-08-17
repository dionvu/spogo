package search

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/session"
	"github.com/dionv/spogo/spotify"
	"github.com/dionv/spogo/spotify/api/headers"
	"github.com/dionv/spogo/spotify/api/status"
	"github.com/dionv/spogo/spotify/api/urls"
)

type Response struct {
	Tracks    TracksResponse    `json:"tracks"`
	Artists   ArtistsResponse   `json:"artists"`
	Albums    AlbumsResponse    `json:"albums"`
	Playlists PlaylistsResponse `json:"playlists"`
	Shows     ShowsResponse     `json:"shows"`
	Episodes  EpisodesResponse  `json:"episodes"`
}

type TracksResponse struct {
	Href     string          `json:"href"`
	Limit    int             `json:"limit"`
	Next     string          `json:"next"`
	Offset   int             `json:"offset"`
	Previous string          `json:"previous"`
	Total    int             `json:"total"`
	Items    []spotify.Track `json:"items"`
}

type ArtistsResponse struct {
	Href     string           `json:"href"`
	Limit    int              `json:"limit"`
	Next     string           `json:"next"`
	Offset   int              `json:"offset"`
	Previous string           `json:"previous"`
	Total    int              `json:"total"`
	Items    []spotify.Artist `json:"items"`
}

type PlaylistsResponse struct {
	Href     string             `json:"href"`
	Limit    int                `json:"limit"`
	Next     string             `json:"next"`
	Offset   int                `json:"offset"`
	Previous string             `json:"previous"`
	Total    int                `json:"total"`
	Items    []spotify.Playlist `json:"items"`
}

type ShowsResponse struct {
	Href     string         `json:"href"`
	Limit    int            `json:"limit"`
	Next     string         `json:"next"`
	Offset   int            `json:"offset"`
	Previous string         `json:"previous"`
	Total    int            `json:"total"`
	Items    []spotify.Show `json:"items"`
}

type EpisodesResponse struct {
	Href     string          `json:"href"`
	Limit    int             `json:"limit"`
	Next     string          `json:"next"`
	Offset   int             `json:"offset"`
	Previous string          `json:"previous"`
	Total    int             `json:"total"`
	Items    []spotify.Album `json:"items"`
}

type AlbumsResponse struct {
	Href     string          `json:"href"`
	Limit    int             `json:"limit"`
	Next     string          `json:"next"`
	Offset   int             `json:"offset"`
	Previous string          `json:"previous"`
	Total    int             `json:"total"`
	Items    []spotify.Album `json:"items"`
}

func Search(input string, searchType []string, s *session.Session) (*Response, error) {
	r := &Response{}

	query := url.Values{}
	query.Set("q", input)
	query.Set("type", strings.Join(searchType, ","))
	query.Set("limit", "10")

	req, err := http.NewRequest(http.MethodGet, urls.SEARCH+"?"+query.Encode(), nil)
	if err != nil {
		return nil, errors.HTTPError.WrapWithNoMessage(err)
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTPError.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BadToken {
		return nil, errors.ReauthenticationError.NewWithNoMessage()
	}

	if res.StatusCode != status.Ok {
		return nil, errors.HTTPError.New("bad request")
	}

	if err = json.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
