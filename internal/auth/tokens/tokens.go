package tokens

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	errx "github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/pkg/utils"
)

// Checks if provided access token is valid by sending a spotify api call.
func EnsureValidAccessToken(accessToken string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return false, errx.HTTPRequestError.Wrap(err, "Failed to create new http request")
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errx.HTTPRequestError.Wrap(err, "Failed to get http result")
	}

	if res.StatusCode != 200 {
		return false, nil
	}

	return true, nil
}

// Given the refresh token, returns a new access token.
func GetNewToken(refreshToken string) (string, error) {
	clientID := os.Getenv("SPOTIFY_ID")
	spotifySecret := os.Getenv("SPOTIFY_SECRET")

	// Addes required query parameters to the /api/token endpoint.
	url := func() string {
		spotifyTokenUrl := "https://accounts.spotify.com/api/token"
		query := url.Values{}
		query.Add("grant_type", "refresh_token")
		query.Add("refresh_token", refreshToken)

		return spotifyTokenUrl + "?" + query.Encode()
	}()

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", errx.HTTPRequestError.Wrap(err, "Unable to create new http request")
	}

	encodedImportantStuff := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + spotifySecret))

	// Addes required headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+encodedImportantStuff)

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != 200 || err != nil {
		return "", errx.ReauthenticationError.Wrap(err, "Bad refresh token.")
	}

	data, err := utils.ParseJsonResponse(res)
	if err != nil {
		return "", errx.JSONError.Wrap(err, "Failed to parse response")
	}

	return data["access_token"].(string), nil
}

// Saves the given access token to .config/spogo/tokens/access_token.json,
// overriding any previously saved access token.
func SaveToken(accessToken string) error {
	homeDir, err := os.UserHomeDir()
	tokenDir := filepath.Join(homeDir, ".config", "spogo", "tokens")

	if err = os.MkdirAll(tokenDir, os.ModePerm); err != nil {
		return errx.FileError.Wrap(err, "Failed to create config directory")
	}

	file, err := os.Create(tokenDir + "/access_token.json")
	if err != nil {
		return errx.FileError.Wrap(err, "Failed to save tokens")
	}
	defer file.Close()

	tok := map[string]string{"access_token": accessToken}

	jsonBody, err := json.Marshal(&tok)
	if err != nil {
		return errx.JSONError.Wrap(err, "Failed to marshal token data")
	}

	_, err = file.Write(jsonBody)
	if err != nil {
		return errx.FileError.Wrap(err, "Failed to write json to file")
	}

	return nil
}

// Saves the given refresh token to .config/spogo/tokens/refresh_token.json,
// overriding any previously saved refresh token.
func SetRefreshToken(refreshToken string) error {
	homeDir, err := os.UserHomeDir()
	tokenDir := filepath.Join(homeDir, ".config", "spogo", "tokens")

	if err = os.MkdirAll(tokenDir, os.ModePerm); err != nil {
		return errx.FileError.Wrap(err, "Failed to create config directory")
	}

	file, err := os.Create(tokenDir + "/refresh_token.json")
	if err != nil {
		return errx.FileError.Wrap(err, "Failed to save tokens")
	}
	defer file.Close()

	tok := map[string]string{"refresh_token": refreshToken}

	jsonBody, err := json.Marshal(&tok)
	if err != nil {
		return errx.JSONError.Wrap(err, "Failed to marshal token data")
	}

	_, err = file.Write(jsonBody)
	if err != nil {
		return errx.FileError.Wrap(err, "Failed to write json to file")
	}

	return nil
}

func GetAccessToken() (string, error) {
	homeDir, err := os.UserHomeDir()

	tokenDir := filepath.Join(homeDir, ".config", "spogo", "tokens")

	res, err := os.ReadFile(tokenDir + "/access_token.json")
	if err != nil {
		return "", errx.FileError.Wrap(err, "Unable to read access_token.json")
	}

	var data map[string]string

	err = json.Unmarshal(res, &data)
	if err != nil {
		return "", errx.JSONError.Wrap(err, "Failed to unmarshal token data")
	}

	return data["access_token"], nil
}

func GetRefreshToken() (string, error) {
	homeDir, err := os.UserHomeDir()

	tokenDir := filepath.Join(homeDir, ".config", "spogo", "tokens")

	res, err := os.ReadFile(tokenDir + "/refresh_token.json")
	if err != nil {
		return "", errx.FileError.Wrap(err, "Unable to read refresh_token.json")
	}

	var data map[string]string

	err = json.Unmarshal(res, &data)
	if err != nil {
		return "", errx.JSONError.Wrap(err, "Failed to unmarshal token data")
	}

	return data["refresh_token"], nil
}
