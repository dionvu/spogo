package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

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
func Authenticate() string {
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

// Does some wacky shit lmao
// lol fsafdsfadsfdsa fdafdsaafsdj;kjl fdsallf;dsajfdsa fdsljkafdjs;al;dsf
// fsdafadsfds fdsafdssdfds fdafdsafdsaj; fdsasfdas
func startAuthentication(w http.ResponseWriter, r *http.Request) {
	state = uuid.New().String()

	authUrl := getAuthUrl()

	// Sends our authentication url
	req, err := http.NewRequest(http.MethodGet, authUrl, nil)
	if err != nil {
		utils.LogError("Creating request", nil)
	}

	// Spotify sets up the authentication url
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		utils.LogError("Unable to setup spotify auth url", err)
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
		utils.LogError("Invalid state.", nil)
	}
	if err != "" {
		utils.LogError("Failed to complete authentication", fmt.Errorf(err))
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
