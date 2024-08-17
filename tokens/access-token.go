package tokens

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dionv/spogo/config"
	"github.com/dionv/spogo/errors"
)

type AccessToken struct {
	Token  string    `json:"access_token"`
	Expiry time.Time `json:"time_created"`
}

func NewAccessToken(str string) *AccessToken {
	return &AccessToken{
		Token:  str,
		Expiry: time.Now().Add(time.Hour),
	}
}

// Loads the access token from token file
func (t *AccessToken) Load(c *config.Config) error {
	path := filepath.Join(c.CachePath(), config.ACCESSTOKENFILE)
	file, err := os.Open(path)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to read token file: %v", path))
	}

	err = json.Unmarshal(b, t)
	if err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal token")
	}

	return nil
}

// Refreshes the access token via valid refresh token.
// Then updates the token string and token file.
func (t *AccessToken) Refresh(refreshToken *RefreshToken, c *config.Config) error {
	query := url.Values{}
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", refreshToken.String())

	ep := "https://accounts.spotify.com/api/token"
	req, err := http.NewRequest(http.MethodPost, ep, strings.NewReader(query.Encode()))
	if err != nil {
		return errors.HTTPRequestError.Wrap(err, "Failed to make a request for new access token")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Spotify.ClientID(), c.Spotify.ClientSecret())

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != 200 || err != nil {
		return errors.ReauthenticationError.Wrap(err, "Bad refresh token")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to read response body")
	}

	if err = json.Unmarshal(body, t); err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal response body")
	}

	t.Update(t.String(), c)

	return nil
}

// Update Updates the token value, and replaces the contents of the token
// file with the new token and an updated expiry time.
func (t *AccessToken) Update(tok string, c *config.Config) error {
	t.Token = tok
	t.Expiry = time.Now().Add(time.Hour)

	path := c.CachePath()
	os.MkdirAll(path, os.ModePerm)

	file, err := os.Create(filepath.Join(path, config.ACCESSTOKENFILE))
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	b, err := json.Marshal(t)
	if err != nil {
		return errors.JSONError.Wrap(err, fmt.Sprintf("Failed to marshal token body: %v", *t))
	}

	_, err = file.Write(b)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to write new token to file: %v", path)
	}

	return nil
}

// Returns the token as a string
func (t *AccessToken) String() string {
	return t.Token
}