package config

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dionvu/spogo/errors"
)

type Credentials struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// Attempts to do the "client credentials" authentication flow
// to test validity of spotify client ID and client secret.
func (c *Credentials) Valid() (bool, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	ep := "https://accounts.spotify.com/api/token"
	req, err := http.NewRequest(http.MethodPost, ep, strings.NewReader(data.Encode()))
	if err != nil {
		err = errors.HTTPRequest.Wrap(err, "failed to make new request with token")
		errors.LogError(err)
		return false, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.ClientID, c.ClientSecret)

	if res, err := http.DefaultClient.Do(req); err != nil || res.StatusCode != http.StatusOK {
		err = errors.HTTP.Wrap(err, fmt.Sprintf("unable to do http request: %v", err))
		errors.LogError(err)
		return false, err
	}

	return true, nil
}
