package tokens

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/app/config"
)

type AccessToken struct {
	Token
}

// Loads the access token from "config.yaml"
func (tok *AccessToken) Load(path string) error {
	return tok.load(path, "access_token")
}

// Refreshes the access token via valid refresh token.
func (tok *AccessToken) Refresh(refreshToken *RefreshToken, id string, secret string) error {
	url := func() string {
		spotifyTokenUrl := "https://accounts.spotify.com/api/token"
		query := url.Values{}
		query.Add("grant_type", "refresh_token")
		query.Add("refresh_token", refreshToken.String())

		return spotifyTokenUrl + "?" + query.Encode()
	}()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return errors.HTTPRequestError.Wrap(err, "Failed to make a request for new access token")
	}

	encodedImportantStuff := base64.StdEncoding.EncodeToString([]byte(id + ":" + secret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+encodedImportantStuff)

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != 200 || err != nil {
		return errors.ReauthenticationError.Wrap(err, "Failed to get response for new access token, bad refresh token")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to read response body")
	}

	data := &struct {
		Access_token string
	}{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal response body")
	}

	tok.token = data.Access_token

	return nil
}

// Tests validity of the token by making a request.
func (tok *AccessToken) IsValid() (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return false, errors.HTTPRequestError.Wrap(err, "Failed to create new http request")
	}
	req.Header.Add("Authorization", "Bearer "+tok.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errors.HTTPRequestError.Wrap(err, "Failed to get http result")
	}

	if res.StatusCode != 200 {
		return false, nil
	}

	return true, nil
}

// Updates token file and sets new token.
func (t *AccessToken) Update(tok string, c *config.Config) error {
	path, err := os.UserConfigDir()
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to get user's config dir")
	}

	path = filepath.Join(c.Path(), config.ACCESSTOKENFILE)

	return t.update(tok, path, "access_token")
}
