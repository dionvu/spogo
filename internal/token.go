package auth

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/dionv/spogo/utils"
)

// Exchanges authentication code for an access token and refresh token
func ExchangeToken(code string) (string, string) {
	clientID := os.Getenv("SPOTIFY_ID")
	spotifySecret := os.Getenv("SPOTIFY_SECRET")

	url := getTokenUrl(code)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		utils.LogError("Failed to get token URL", err)
	}

	// Encodes in base64 and formates in required format
	a := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + spotifySecret))
	a = "Basic " + a

	// Add required headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", a)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogError("Failed to get token URL", err)
	}
	if res.StatusCode != 200 {
		utils.LogError("Failed to get token URL", nil)
	}

	// Parse Json
	data := map[string]interface{}{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		utils.LogError("Failed to read response", err)
	}

	json.Unmarshal(body, &data)

	return data["access_token"].(string), data["refresh_token"].(string)
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
