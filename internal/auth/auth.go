package auth

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"os"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/auth/tokens"
	"github.com/dionv/spogo/pkg/utils"
	"github.com/joho/godotenv"
)

var (
	ch    = make(chan string)
	state string
)

// Checks if the user has a valid access token, else retrieves a new access token
// from refresh token. Authenticates user if neither are valid.
func EnsureAuthenticated() error {
	accessToken, accessErr := tokens.GetAccessToken()
	refreshToken, refreshErr := tokens.GetRefreshToken()

	if accessErr == nil && refreshErr == nil {
		isValidToken, err := tokens.EnsureValidAccessToken(accessToken)
		if err != nil {
			return err
		}
		if isValidToken {
			return nil
		}

		newAccessToken, err := tokens.GetNewToken(refreshToken)
		if err != nil {
			return AuthenticateUser()
		}

		return tokens.SaveToken(newAccessToken)

	}

	return AuthenticateUser()
}

// Fetches new authentication code, exchaning for new access and refresh tokens.
func AuthenticateUser() error {
	code := getAuthCode()
	newToken, newRefreshToken, err := exchangeForToken(code)
	if err != nil {
		return err
	}

	if err = tokens.SaveToken(newToken); err != nil {
		return err
	}

	return tokens.SetRefreshToken(newRefreshToken)
}

// Fetches authentication code.
func getAuthCode() string {
	godotenv.Load()

	http.HandleFunc("/", startAuthentication)
	http.HandleFunc("/callback", completeAuthentication)

	startServer()

	uri := "http://localhost:" + os.Getenv("PORT")
	utils.OpenURL(uri)

	code := <-ch

	return code
}

// Exchanges a valid authentication code for an access token and a refresh token.
func exchangeForToken(code string) (string, string, error) {
	clientID := os.Getenv("SPOTIFY_ID")
	spotifySecret := os.Getenv("SPOTIFY_SECRET")

	// Addes required query parameters to the /api/token endpoint.
	url := func() string {
		spotifytokenurl := "https://accounts.spotify.com/api/token"

		redirectUri := os.Getenv("REDIRECT_URI")
		query := url.Values{}
		query.Add("grant_type", "authorization_code")
		query.Add("code", code)
		query.Add("redirect_uri", redirectUri)
		return spotifytokenurl + "?" + query.Encode()
	}()

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
