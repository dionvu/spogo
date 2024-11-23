package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/spotify/api/scopes"
	"github.com/dionvu/spogo/spotify/api/urls"
	"github.com/fatih/color"
	"github.com/google/uuid"
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
// If the access token is not valid, reauthenticates. Updating the token
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
				errors.Log(err)
				return err
			}
		}
	}

	return nil
}

// Forces reauthentication.
func (s *Session) Reauth(c *config.Config) error {
	if err := s.AccessToken.Refresh(s.RefreshToken, c); err != nil {
		if err := getNewTokens(s, c); err != nil {
			errors.Log(err)
			return err
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

		return <-ch
	}()

	query := url.Values{}
	query.Set("grant_type", "authorization_code")
	query.Set("redirect_uri", REDIRECT_URI)
	query.Set("code", code)

	ep := "https://accounts.spotify.com/api/token"
	req, err := http.NewRequest(http.MethodPost, ep, strings.NewReader(query.Encode()))
	if err != nil {
		err = errors.HTTPRequest.Wrap(err, "unable to create new http request for new token")
		errors.Log(err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Spotify.ClientID, c.Spotify.ClientSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTP.Wrap(err, "unable to get http response")
		errors.Log(err)
		return err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.HTTP.Wrap(err, "failed to read response body")
		errors.Log(err)
		return err
	}

	data := map[string]interface{}{}
	if err = json.Unmarshal(b, &data); err != nil {
		err = errors.JSONUnmarshal.Wrap(err, "failed to unmarshal response body: %v", string(b))
		errors.Log(err)
		return err
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
		errors.Log(err)
		return err
	}

	fmt.Println(color.HiGreenString("Opening -> " + url))

	return nil
}

func startServer() {
	go func() {
		err := http.ListenAndServe(":"+PORT, nil)
		if err != nil {
			err := errors.HTTP.Wrap(err, fmt.Sprintf("failed to start server on port: %v", PORT))
			errors.Log(err)
			log.Fatal(err)
		}
	}()
}

//

// Redirects user to the spotify authentication url and awaits callback.
func startAuth(w http.ResponseWriter, r *http.Request) {
	state = uuid.NewString()

	query := url.Values{}
	query.Set("redirect_uri", REDIRECT_URI)
	query.Set("response_type", "code")
	query.Set("client_id", clientID)
	query.Set("scope", strings.Join([]string{
		scopes.UserReadPrivate,
		scopes.UserReadEmail,
		scopes.UserReadPlaybackState,
		scopes.UserModifyPlaybackState,
		scopes.UserPlaylistRead,
		scopes.UserReadCollab,
	}, " "))
	query.Set("state", state)

	req, err := http.NewRequest(http.MethodGet, spotifyurls.SPOTIFYAUTHURL, strings.NewReader(query.Encode()))
	if err != nil {
		log.Fatal(errors.HTTPRequest.Wrap(err, "unable to create new http request for spotify authentication url"))
	}

	if _, err = http.DefaultClient.Do(req); err != nil {
		log.Fatal(errors.HTTP.Wrap(err, "unable to do http request"))
	}

	http.Redirect(w, r, fmt.Sprintf("%s?%s", spotifyurls.SPOTIFYAUTHURL, query.Encode()), http.StatusTemporaryRedirect)
}

// After user is redirected to the redirect uri, ensures valid state
// and fetches the authentication code.
func completeAuth(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		log.Fatal("invalid state")
	}

	if r.URL.Query().Get("error") != "" {
		log.Fatal("failed to complete authentication")
	}

	ch <- r.URL.Query().Get("code")

	fmt.Fprintln(w, "authentication success!")
}

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
		err = errors.FileOpen.Wrap(err, fmt.Sprintf("failed to open token file path: %v", path))
		errors.Log(err)
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		err = errors.FileRead.Wrap(err, fmt.Sprintf("failed to read token file: %v", path))
		errors.Log(err)
		return err
	}

	err = json.Unmarshal(b, t)
	if err != nil {
		err = errors.JSONUnmarshal.Wrap(err, fmt.Sprintf("failed to unmarshal token body %v", string(b)))
		errors.Log(err)
		return err
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
		err = errors.HTTPRequest.Wrap(err, "failed to make a request for new access token")
		errors.Log(err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Spotify.ClientID, c.Spotify.ClientSecret)

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK || err != nil {
		err = errors.Reauthentication.Wrap(err, "bad refresh token")
		errors.Log(err)
		return err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		err = errors.FileRead.Wrap(err, fmt.Sprintf("failed to read response body"))
		errors.Log(err)
		return err
	}

	if err = json.Unmarshal(b, t); err != nil {
		err = errors.JSONUnmarshal.Wrap(err, fmt.Sprintf("failed to unmarshal response body: %v", string(b)))
		errors.Log(err)
		return err
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
		return errors.FileCreate.Wrap(err, fmt.Sprintf("failed to open token file path: %v", path))
	}
	defer file.Close()

	b, err := json.Marshal(t)
	if err != nil {
		return errors.JSONMarshal.Wrap(err, fmt.Sprintf("failed to marshal token body: %v", *t))
	}

	_, err = file.Write(b)
	if err != nil {
		return errors.FileWrite.Wrap(err, fmt.Sprintf("failed to write new token to file: %v", path))
	}

	return nil
}

// Returns the token as a string
func (t *AccessToken) String() string {
	return t.Token
}

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

func NewRefreshToken(tok string) *RefreshToken {
	t := &RefreshToken{
		Token: tok,
	}

	return t
}

// Loads the token fields from the refresh token file.
func (t *RefreshToken) Load(c *config.Config) error {
	path := filepath.Join(c.CachePath(), config.REQUESTTOKENFILE)

	file, err := os.Open(path)
	if err != nil {
		err = errors.FileOpen.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
		errors.Log(err)
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		err = errors.FileRead.Wrap(err, fmt.Sprintf("Failed to read token file: %v", path))
		errors.Log(err)
		return err
	}

	err = json.Unmarshal(b, t)
	if err != nil {
		err = errors.JSONUnmarshal.Wrap(err, "Failed to unmarshal token from file body: %v", string(b))
		errors.Log(err)
		return err
	}

	return nil
}

// Updates the token file with new token.
func (t *RefreshToken) Update(tok string, c *config.Config) error {
	t.Token = tok

	filePath := filepath.Join(c.CachePath(), config.REQUESTTOKENFILE)
	file, err := os.Create(filePath)
	if err != nil {
		err = errors.FileCreate.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", filePath))
		errors.Log(err)
		return err
	}
	defer file.Close()

	b, err := json.Marshal(t)
	if err != nil {
		err = errors.JSONMarshal.Wrap(err, fmt.Sprintf("Failed to marshal token body: %v", *t))
		errors.Log(err)
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		err = errors.FileWrite.Wrap(err, fmt.Sprintf("Failed to write new token to file: %v", filePath))
		errors.Log(err)
		return err
	}

	return nil
}

// The token as a string
func (t *RefreshToken) String() string {
	return t.Token
}

type Session struct {
	AccessToken  *AccessToken
	RefreshToken *RefreshToken
}

// Creates a new session, loading tokens from respective files, and authenticating.
func New(c *config.Config) (*Session, error) {
	s := &Session{
		AccessToken:  &AccessToken{},
		RefreshToken: &RefreshToken{},
	}

	// Loads possible access token and refresh token from respective token files.
	s.AccessToken.Load(c)
	s.RefreshToken.Load(c)

	// Authenticates valid access token, or valid access token and refresh token.
	err := s.Authenticate(c)
	if err != nil {
		errors.Log(err)
		return nil, err
	}

	return s, nil
}
