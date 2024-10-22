package spotify

import (
	"encoding/json"
	"net/http"

	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/spotify/auth"
)

type User struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	ID          string `json:"id"`

	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`

	Followers struct {
		Href string `json:"href"`

		Total int `json:"total"`
	} `json:"followers"`

	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
}

func New(s *auth.Session) (*User, error) {
	ep := "https://api.spotify.com/v1/me"
	req, _ := http.NewRequest(http.MethodGet, ep, nil)

	req.Header.Add("Authorization", "Bearer "+s.AccessToken.String())

	res, _ := http.DefaultClient.Do(req)

	u := &User{}

	err := json.NewDecoder(res.Body).Decode(u)
	if err != nil {
		err = errors.JSONDecode.Wrap(err, "failed to decode user response")
		errors.Log(err)
		return nil, err
	}

	return u, nil
}
