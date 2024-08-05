package user

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/tokens"
	"github.com/dionv/spogo/pkg/utils"
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

func (u *User) Authenticate(c *config.Config) error {
	at := u.AccessToken
	rt := u.RefreshToken

	if at.String() == "" && rt.String() == "" {
		return u.getNewTokens(c)
	}

	isValid, err := at.IsValid()
	if err != nil {
		return err
	}
	if isValid {
		return nil
	}

	err = at.Refresh(rt, c)
	if err != nil {
		return u.getNewTokens(c)
	}

	return nil
}

func (u *User) getNewTokens(c *config.Config) error {
	code := getAuthCode(c)

	at, rt, err := exchangeForToken(code)
	if err != nil {
		return err
	}

	if err = u.AccessToken.Update(at.String(), at.TimeCreated, c); err != nil {
		return err
	}

	if err = u.RefreshToken.Update(rt.String(), c); err != nil {
		return err
	}

	return nil
}

func getAuthCode(c *config.Config) string {
	clientID = c.Spotify.ClientID()
	clientSecret = c.Spotify.ClientSecret()

	http.HandleFunc("/", startAuth)
	http.HandleFunc("/callback", completeAuth)

	startServer()

	utils.OpenURL(URI)

	code := <-ch

	return code
}

func exchangeForToken(code string) (*tokens.AccessToken, *tokens.RefreshToken, error) {
	spotifytokenurl := "https://accounts.spotify.com/api/token"

	query := url.Values{}
	query.Set("grant_type", "authorization_code")
	query.Set("code", code)
	query.Set("redirect_uri", REDIRECT_URI)

	req, err := http.NewRequest(http.MethodPost, spotifytokenurl, strings.NewReader(query.Encode()))
	if err != nil {
		return nil, nil, errors.HTTPRequestError.Wrap(err, "Unable to create new http request")
	}

	stuff := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+stuff)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, errors.HTTPRequestError.Wrap(err, "Unable to get http response")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, errors.FileError.Wrap(err, "Failed to read response body")
	}

	fmt.Println(string(body))

	data := map[string]interface{}{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, nil, errors.JSONError.Wrap(err, "Failed to unmarshal response body")
	}

	return tokens.NewAccessToken(data["access_token"].(string)), tokens.NewRefreshToken(data["refresh_token"].(string)), nil
}
