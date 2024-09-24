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

const (
	TRACK_TYPE    = "track"
	ALBUM_TYPE    = "album"
	PLAYLIST_TYPE = "playlist"
)

type Response struct {
	Tracks    TracksResponse    `json:"tracks"`
	Artists   ArtistsResponse   `json:"artists"`
	Albums    AlbumsResponse    `json:"albums"`
	Playlists PlaylistsResponse `json:"playlists"`
	Shows     ShowsResponse     `json:"shows"`
	Episodes  EpisodesResponse  `json:"episodes"`
}

type PlaylistsResponse struct {
	Href     string     `json:"href"`
	Limit    int        `json:"limit"`
	Next     string     `json:"next"`
	Offset   int        `json:"offset"`
	Previous string     `json:"previous"`
	Total    int        `json:"total"`
	Items    []Playlist `json:"items"`
}

type SearchResult struct {
	Tracks    []*Track
	Albums    []*Album
	Playlists []*Playlist
	Artists   []*Artist
	Shows     []*Show
	Episodes  []*Episode
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

func Search(input string, searchType []string, s *auth.Session) (*SearchResult, error) {
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

	searchResult := &SearchResult{
		Tracks:    []*Track{},
		Albums:    []*Album{},
		Playlists: []*Playlist{},
		Artists:   []*Artist{},
		Shows:     []*Show{},
		Episodes:  []*Episode{},
	}

	for _, track := range r.Tracks.Items {
		searchResult.Tracks = append(searchResult.Tracks, &track)
	}

	for _, album := range r.Albums.Items {
		searchResult.Albums = append(searchResult.Albums, &album)
	}

	for _, playlist := range r.Playlists.Items {
		searchResult.Playlists = append(searchResult.Playlists, &playlist)
	}

	for _, artist := range r.Artists.Items {
		searchResult.Artists = append(searchResult.Artists, &artist)
	}

	for _, show := range r.Shows.Items {
		searchResult.Shows = append(searchResult.Shows, &show)
	}

	for _, episode := range r.Episodes.Items {
		searchResult.Episodes = append(searchResult.Episodes, &episode)
	}

	return searchResult, nil
}
