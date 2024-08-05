package tokens

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
)

type AccessToken struct {
	Token       string    `json:"access_token"`
	TimeCreated time.Time `json:"time_created"`
}

func NewAccessToken(tok string) *AccessToken {
	return &AccessToken{
		Token:       tok,
		TimeCreated: time.Now(),
	}
}

// Returns the token as a string
func (t *AccessToken) String() string {
	return t.Token
}

// Loads the access token from token file
func (t *AccessToken) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to read token file: %v", path))
	}

	// var data map[string]string
	err = json.Unmarshal(b, t)
	if err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal token")
	}

	return nil
}

// Refreshes the access token via valid refresh token.
func (tok *AccessToken) Refresh(refreshToken *RefreshToken, c *config.Config) error {
	id := c.Spotify.ClientID()
	secret := c.Spotify.ClientSecret()

	spotifyUrl := "https://accounts.spotify.com/api/token"
	query := url.Values{}
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", refreshToken.String())

	req, err := http.NewRequest(http.MethodPost, spotifyUrl, strings.NewReader(query.Encode()))
	if err != nil {
		return errors.HTTPRequestError.Wrap(err, "Failed to make a request for new access token")
	}

	encodedImportantStuff := base64.StdEncoding.EncodeToString([]byte(id + ":" + secret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedImportantStuff)

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != 200 || err != nil {
		return errors.ReauthenticationError.Wrap(err,
			"Failed to get response for new access token, bad refresh token")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to read response body")
	}

	// data := &struct {
	// 	Access_token string
	// }{}
	err = json.Unmarshal(body, tok)
	if err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal response body")
	}

	tok.Update(tok.String(), time.Now(), c)

	return nil
}

// Tests validity of the token by making a request.
func (tok *AccessToken) IsValid() (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return false, errors.HTTPRequestError.Wrap(err, "Failed to create new http request")
	}
	req.Header.Set("Authorization", "Bearer "+tok.String())

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
func (t *AccessToken) Update(tok string, timeCreated time.Time, c *config.Config) error {
	path, err := os.UserConfigDir()
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to get user's config dir")
	}

	path = filepath.Join(c.Path(), config.TOKENSDIRECTORY)

	os.MkdirAll(path, os.ModePerm)

	path = filepath.Join(path, config.ACCESSTOKENFILE)

	t.TimeCreated = timeCreated

	return t.saveToFile(path)
}

func (t *AccessToken) saveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	body, err := json.Marshal(t)
	if err != nil {
		return errors.JSONError.Wrap(err, fmt.Sprintf("Failed to marshal token body: %v", *t))
	}

	_, err = file.Write(body)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to write new token to file: %v", path)
	}

	return nil
}
