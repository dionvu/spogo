package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/status"
	"github.com/dionvu/spogo/spotify/api/urls"
)

type Response struct {
	Tracks    TracksResponse    `json:"tracks"`
	Artists   ArtistsResponse   `json:"artists"`
	Albums    AlbumsResponse    `json:"albums"`
	Playlists PlaylistsResponse `json:"playlists"`
	Shows     ShowsResponse     `json:"shows"`
	Episodes  EpisodesResponse  `json:"episodes"`
}

type ArtistsResponse struct {
	Href     string   `json:"href"`
	Limit    int      `json:"limit"`
	Next     string   `json:"next"`
	Offset   int      `json:"offset"`
	Previous string   `json:"previous"`
	Total    int      `json:"total"`
	Items    []Artist `json:"items"`
}

type ShowsResponse struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []Show `json:"items"`
}

type EpisodesResponse struct {
	Href     string    `json:"href"`
	Limit    int       `json:"limit"`
	Next     string    `json:"next"`
	Offset   int       `json:"offset"`
	Previous string    `json:"previous"`
	Total    int       `json:"total"`
	Items    []Episode `json:"items"`
}

func Search(input string, searchType []string, s *auth.Session) (*Response, error) {
	r := &Response{}

	query := url.Values{}
	query.Set("q", input)
	query.Set("type", strings.Join(searchType, ","))
	query.Set("limit", "20")

	req, err := http.NewRequest(http.MethodGet, urls.SEARCH+"?"+query.Encode(), nil)
	if err != nil {
		return nil, errors.HTTPRequest.Wrap(err, fmt.Sprintf("failed to make request for search query: %v", input))
	}
	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.HTTPRequest.WrapWithNoMessage(err)
	}

	if res.StatusCode == status.BadToken {
		return nil, errors.Reauthentication.NewWithNoMessage()
	}

	if res.StatusCode != status.Ok {
		return nil, errors.HTTP.New("bad request")
	}

	if err = json.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
