package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
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
			fmt.Printf("%v %v %v\n", color.RedString("Error:"),
				"invalid spotify client credentials:", color.YellowString(c.FilePath()))
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
		clientID = c.Spotify.ClientID
		clientSecret = c.Spotify.ClientSecret

		http.HandleFunc("/", startAuth)
		http.HandleFunc("/callback", completeAuth)

		startServer()

		if err := OpenURL(URI); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		}

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
		return errors.HTTPRequest.Wrap(err, "unable to create new http request for new token")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Spotify.ClientID, c.Spotify.ClientSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.HTTP.Wrap(err, "unable to get http response")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.FileRead.Wrap(err, "failed to read response body")
	}

	data := map[string]interface{}{}
	if err = json.Unmarshal(b, &data); err != nil {
		return errors.JSONUnmarshal.Wrap(err, "failed to unmarshal response body: %v", string(b))
	}

	s.AccessToken.Update(data["access_token"].(string), c)
	s.RefreshToken.Update(data["refresh_token"].(string), c)

	os.Exit(0)

	return nil
}

func OpenURL(url string) error {
	var cmd *exec.Cmd

	os := runtime.GOOS

	switch {
	case os == "windows":
		cmd = exec.Command("start", url)
	case os == "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println(color.HiGreenString("Opening -> ", url))

	return nil
}
