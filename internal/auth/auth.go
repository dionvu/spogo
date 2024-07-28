// Deal with everything up until recieving an authentication code
// and exchanging for a token.

package auth

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/auth/tokens"
	"github.com/dionv/spogo/utils"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const (
	// Hides port in .env, also me shows port
	SPOTIFYAUTHURL  = "https://accounts.spotify.com/authorize"
	SPOTIFYTOKENURL = "https://accounts.spotify.com/api/token"
)

var (
	ch    = make(chan string)
	state string
)

// Sets up & redirects to a Spotify authentication URL.
// Returns an authentication code to exchange for an access token.
func GetAuthCode() string {
	godotenv.Load()

	http.HandleFunc("/", startAuthentication)
	http.HandleFunc("/callback", completeAuthentication)

	startServer()

	uri := "http://localhost:" + os.Getenv("PORT")
	utils.OpenURL(uri)

	// Await authentication to finish
	code := <-ch

	return code
}

// Exchanges authentication code for an access token and refresh token
func ExchangeForToken(code string) (string, string, error) {
	clientID := os.Getenv("SPOTIFY_ID")
	spotifySecret := os.Getenv("SPOTIFY_SECRET")

	url := getTokenUrl(code)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", "", errors.HTTPRequestError.Wrap(err, "Unable to create new http request")
	}

	encodedImportantStuff := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + spotifySecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+encodedImportantStuff)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", errors.HTTPRequestError.Wrap(err, "Unable to get http response")
	}
	if res.StatusCode != 200 {
		return "", "", errors.ReauthenticationError.Wrap(err, "Bad authentication code")
	}

	data, err := utils.ParseJsonResponse(res)

	return data["access_token"].(string), data["refresh_token"].(string), nil
}

func startAuthentication(w http.ResponseWriter, r *http.Request) {
	state = uuid.New().String()

	authUrl := getAuthUrl()

	// Sends our authentication url
	req, err := http.NewRequest(http.MethodGet, authUrl, nil)
	if err != nil {
		err = errors.HTTPRequestError.Wrap(err, "Unable to create new http request")
		log.Fatal(err)
	}

	// Spotify sets up the authentication url
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTPRequestError.Wrap(err, "Unable to do http request")
		log.Fatal(err)
	}

	// Redirects user to authentication url then to callback
	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}

// Verifies the response. Sends authentication code in
// response through awaiting channel.
func completeAuthentication(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	responseState := r.URL.Query().Get("state")
	err := r.URL.Query().Get("error")

	if responseState != state {
		log.Println("Invalid state")
		os.Exit(1)
	}
	if err != "" {
		log.Fatal("Failed to complete authentication")
	}

	fmt.Fprintln(w, "Login success!")

	ch <- code
}

// Returns an spotify authentication url in the form:
// AUTHURL?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s
func getAuthUrl() string {
	scope := "user-read-private user-read-email"
	clientID := os.Getenv("SPOTIFY_ID")
	redirectUri := os.Getenv("REDIRECT_URI")

	query := url.Values{}
	query.Add("client_id", clientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirectUri)
	query.Add("scope", scope)
	query.Add("state", state)

	return fmt.Sprintf("%s?%s", SPOTIFYAUTHURL, query.Encode())
}

func getTokenUrl(code string) string {
	redirectUri := os.Getenv("REDIRECT_URI")
	query := url.Values{}
	query.Add("grant_type", "authorization_code")
	query.Add("code", code)
	query.Add("redirect_uri", redirectUri)

	return SPOTIFYTOKENURL + "?" + query.Encode()
}

func AuthenticateUser() error {
	accessToken, accessErr := tokens.GetAccessToken()
	refreshToken, refreshErr := tokens.GetRefreshToken()

	// Able to access both tokens
	if accessErr == nil && refreshErr == nil {

		isValidToken, err := tokens.EnsureValidAccessToken(accessToken)
		if err != nil {
			return err
		}

		// If access token is still valid we leave
		if isValidToken {
			return nil
		}

		// If access token is not valid, try to get a new one using the refresh token
		newAccessToken, err := tokens.GetNewToken(refreshToken)
		// If getting a new token with the refresh token fails, reauthenticate
		// Likely due to bad refresh token
		if err != nil {
			return reauthenticate()
		}

		return tokens.SaveToken(newAccessToken)

	} else {
		// If tokens are not available or errors occurred, reauthenticate
		// Likely due to first time user
		return reauthenticate()
	}
}

func reauthenticate() error {
	code := GetAuthCode()
	newToken, newRefreshToken, err := ExchangeForToken(code)
	if err != nil {
		return err
	}

	if err = tokens.SaveToken(newToken); err != nil {
		return err
	}

	return tokens.SetRefreshToken(newRefreshToken)
}
