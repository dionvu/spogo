package spotify

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify/api/headers"
)

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

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func AlbumTracks(s *auth.Session, albumID string) (*[]AlbumTrack, error) {
	url := "https://api.spotify.com/v1/albums/" + albumID + "/tracks"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = errors.HTTPRequest.Wrap(err, "failed to make request for playlists")
		errors.Log(err)
		return nil, err
	}

	req.Header.Add(headers.Auth, "Bearer "+s.AccessToken.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.WrapWithNoMessage(err)
		errors.Log(err)
		return nil, err
	}

	if res.StatusCode == 401 {
		err = errors.Reauthentication.NewWithNoMessage()
		errors.Log(err)
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.HTTP.New("bad request")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.HTTP.Wrap(err, "failed to read response body")
		errors.Log(err)
		return nil, err
	}

	var response struct {
		Items []AlbumTrack `json:"items"`
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		err = errors.JSONUnmarshal.Wrap(err, "failed to unmarshal playlists response")
		errors.Log(err)
		return nil, err
	}

	return &response.Items, nil
}
