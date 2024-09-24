package spotify

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/status"
)

type AlbumsResponse struct {
	Href     string  `json:"href"`
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
	Items    []Album `json:"items"`
}

type Album struct {
	Images      []Image  `json:"images"`
	Artists     []Artist `json:"artists"`
	TotalTracks int      `json:"total_tracks"`
	AlbumType   string   `json:"album_type"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"release_date"`
	Type        string   `json:"type"`
	Uri         string   `json:"uri"`
}

func AlbumTracks(s *auth.Session, albumID string) (*[]AlbumTrack, error) {
	url := "https://api.spotify.com/v1/albums/" + albumID + "/tracks"

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var response struct {
		Items []AlbumTrack `json:"items"`
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, errors.JSONUnmarshal.Wrap(err, "failed to unmarshal playlists response")
	}

	return &response.Items, nil
}
