package session

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/pkg/utils"
	"github.com/fatih/color"
)

const (
	REDIRECT_URI = "http://localhost:42069/callback"
	URI          = "http://localhost:42069"
	PORT         = "42069"
)

var (
	ch           = make(chan string)
	state        string
	clientID     string
	clientSecret string
)

// Authenticate is set to only run checks after the access token expiry
// period has elapsed. This is for faster runtime, should be perfectly okay
// unless token files are externally tappered.
// Checks if the access token is valid. If not, refreshes the access token.
// If the access token is not valid, reauthenticates s. Updating the token
// file.
func (s *Session) Authenticate(c *config.Config) error {
	if time.Now().After(s.AccessToken.Expiry) {
		validCred, _ := c.Spotify.Valid()
		if !validCred {
			fmt.Printf("%v %v %v\n", color.RedString("Error"), "Invalid spotify client credentials:", color.YellowString(c.FilePath()))
			os.Exit(0)
		}

		if err := s.AccessToken.Refresh(s.RefreshToken, c); err != nil {
			if err := getNewTokens(s, c); err != nil {
				return err
			}
		}
	}

	return nil
}

// Uses client ID and secret to retrieve an authentication code.
// Exchanges code for an access token and a refresh token.
// Updates session tokens and respective token files.
func getNewTokens(s *Session, c *config.Config) error {
	code := func() string {
		// For handlers access.
		clientID = c.Spotify.ClientID()
		clientSecret = c.Spotify.ClientSecret()

		http.HandleFunc("/", startAuth)
		http.HandleFunc("/callback", completeAuth)

		startServer()

		utils.OpenURL(URI)

		// Awaits the authentication code from handlers.
		return <-ch
	}()

	query := url.Values{}
	query.Set("grant_type", "authorization_code")
	query.Set("redirect_uri", REDIRECT_URI)
	query.Set("code", code)

	ep := "https://accounts.spotify.com/api/token"
	req, err := http.NewRequest(http.MethodPost, ep, strings.NewReader(query.Encode()))
	if err != nil {
		return errors.HTTPRequestError.Wrap(err, "Unable to create new http request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Spotify.ClientID(), c.Spotify.ClientSecret())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTPRequestError.Wrap(err, "Unable to get http response")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to read response body")
	}

	data := map[string]interface{}{}
	if err = json.Unmarshal(body, &data); err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal response body")
	}

	s.AccessToken.Update(data["access_token"].(string), c)
	s.RefreshToken.Update(data["refresh_token"].(string), c)

	return nil
}
