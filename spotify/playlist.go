package spotify

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/errors"
	"github.com/dionvu/spogo/spotify/api/headers"
	"github.com/dionvu/spogo/spotify/api/urls"
)

type Playlist struct {
	Images      []Image   `json:"images"`
	Description string    `json:"description"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Followers   Followers `json:"followers"`
	Public      bool      `json:"public"`
	Tracks      struct {
		Total int `json:"total"`
	} `json:"tracks"`
	URI          string `json:"uri"`
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

type Followers struct {
	Total int `json:"total"`
}

func UserPlaylists(s *auth.Session) (*[]Playlist, error) {
	req, err := http.NewRequest(http.MethodGet, spotifyurls.PLAYLISTS, nil)
	if err != nil {
		return nil, errors.HTTPRequest.Wrap(err, "failed to make request for playlists")
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
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.HTTP.Wrap(err, "failed to read response body")
		errors.Log(err)
		return nil, err
	}

	pr := &playlistsSearchResponse{}

	err = json.Unmarshal(b, pr)
	if err != nil {
		err = errors.JSONUnmarshal.Wrap(err, "failed to unmarshal playlists response")
		errors.Log(err)
		return nil, err
	}

	return &pr.Items, nil
}

func PlaylistTracks(s *auth.Session, playlistID string) (*[]Track, error) {
	type PlaylistTrack struct {
		Track Track `json:"track"`
	}

	url := "https://api.spotify.com/v1/playlists/" + playlistID + "/tracks"

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
		err = errors.HTTP.New("bad request")
		errors.Log(err)
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.HTTP.Wrap(err, "failed to read response body")
		errors.Log(err)
		return nil, err
	}

	var response struct {
		Items []PlaylistTrack `json:"items"`
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		err = errors.JSONUnmarshal.Wrap(err, "failed to unmarshal playlists response")
		errors.Log(err)
		return nil, err
	}

	tracks := []Track{}

	for _, playlistTrack := range response.Items {
		tracks = append(tracks, playlistTrack.Track)
	}

	return &tracks, nil
}
