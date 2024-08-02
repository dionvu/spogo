package user

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/app/config"
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
	// at := user.AccessToken.String()
	// rt := user.RefreshToken.String()

	code := getAuthCode(c)

	newAT, newRT, err := exchangeForToken(code)
	if err != nil {
		return err
	}

	fmt.Println(newAT)
	fmt.Println(newRT)

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

func exchangeForToken(code string) (string, string, error) {
	url := func() string {
		spotifytokenurl := "https://accounts.spotify.com/api/token"

		query := url.Values{}
		query.Add("grant_type", "authorization_code")
		query.Add("code", code)
		query.Add("redirect_uri", REDIRECT_URI)
		return spotifytokenurl + "?" + query.Encode()
	}()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", "", errors.HTTPRequestError.Wrap(err, "Unable to create new http request")
	}

	encodedImportantStuff := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+encodedImportantStuff)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", errors.HTTPRequestError.Wrap(err, "Unable to get http response")
	}
	if res.StatusCode != 200 {
		return "", "", errors.ReauthenticationError.Wrap(err, "Bad authentication code")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", errors.FileError.Wrap(err, "Failed to read response body")
	}

	data := &struct {
		Access_token  string
		Refresh_token string
	}{}

	err = json.Unmarshal(body, data)
	if err != nil {
		return "", "", errors.JSONError.Wrap(err, "Failed to unmarshal response body")
	}

	return data.Access_token, data.Refresh_token, nil
}
