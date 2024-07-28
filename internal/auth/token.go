// Deals with everything after recieving a token.

package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/dionv/spogo/utils"
)

const TOKEN_FILE_PATH = "/home/dionv/repos/spogo/token.json"

func GetTokens() (string, string, error) {
	res, err := os.ReadFile(TOKEN_FILE_PATH)
	if err != nil {
		return "", "", err
	}

	data := TokenStorage{}

	err = json.Unmarshal(res, &data)

	return data.AccessToken, data.RefreshToken, err
}

func SaveToken(token string, refreshToken string) {
	file, err := os.Create(TOKEN_FILE_PATH)
	if err != nil {
		utils.LogError("Failed to save token", err)
	}
	defer file.Close()

	jsonBody, err := json.Marshal(&TokenStorage{AccessToken: token, RefreshToken: refreshToken})
	if err != nil {
		utils.LogError("Failed to Marshal refresh token", err)
	}

	_, err = file.Write(jsonBody)
	if err != nil {
		utils.LogError("Failed to write token to file", err)
	}
}

// Checks if tokens are valid, returns either new access token or original token.
func EnsureValidTokens(token string, refreshToken string) string {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		utils.LogError("", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)

	// Expired token
	if res.StatusCode == 401 {
		newToken := getNewToken(refreshToken)
		SaveToken(token, refreshToken)
		return newToken
	}

	// Bad token or refreshToken, likely from user goofiness lol
	if err != nil || res.StatusCode != 200 {
		code := Authenticate()
		token, refreshToken = ExchangeToken(code)
		SaveToken(token, refreshToken)
	}

	return token
}

// Given the refresh token, returns a new access token.
func getNewToken(refreshToken string) string {
	clientID := os.Getenv("SPOTIFY_ID")
	spotifySecret := os.Getenv("SPOTIFY_SECRET")

	url := getRefreshTokenUrl(refreshToken)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		utils.LogError("Failed to refresh token", err)
	}

	// Encodes in base64 and formates in required format
	encodedImportantStuff := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + spotifySecret))
	encodedImportantStuff = "Basic " + encodedImportantStuff

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", encodedImportantStuff)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogError("Failed to get token URL", err)
	}

	data := utils.ParseJsonResponse(res)

	return data["access_token"].(string)
}

// Adds required query parameters onto token endpoint
func getTokenUrl(code string) string {
	redirectUri := os.Getenv("REDIRECT_URI")
	query := url.Values{}
	query.Add("grant_type", "authorization_code")
	query.Add("code", code)
	query.Add("redirect_uri", redirectUri)

	return SPOTIFYTOKENURL + "?" + query.Encode()
}

// Adds required query parameters onto token endpoint
func getRefreshTokenUrl(refreshToken string) string {
	query := url.Values{}
	query.Add("grant_type", "refresh_token")
	query.Add("refresh_token", refreshToken)

	return SPOTIFYTOKENURL + "?" + query.Encode()
}
