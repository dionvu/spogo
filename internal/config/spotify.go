package config

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dionv/spogo/errors"
)

type Spotify struct {
	clientID     string `yaml:"client_id"`
	clientSecret string `yaml:"client_secret"`
}

// Attempts to do the "client credentials" authentication flow
// to test validity of spotify client ID and client secret.
func (s *Spotify) Valid() (bool, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	ep := "https://accounts.spotify.com/api/token"
	req, err := http.NewRequest(http.MethodPost, ep, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("unable to create new http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.ClientID(), s.ClientSecret())

	if res, err := http.DefaultClient.Do(req); err != nil || res.StatusCode != http.StatusOK {
		return false, errors.HTTPError.Wrap(err, fmt.Sprintf("unable to do http request: %v", err))
	}

	return true, nil
}

func (s *Spotify) ClientID() string {
	return s.clientID
}

func (s *Spotify) ClientSecret() string {
	return s.clientSecret
}

func (s *Spotify) setID(str string) {
	s.clientID = str
}

func (s *Spotify) setSecret(str string) {
	s.clientSecret = str
}
