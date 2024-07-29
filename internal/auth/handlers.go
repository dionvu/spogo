package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dionv/spogo/errors"
	"github.com/google/uuid"
)

func startAuthentication(w http.ResponseWriter, r *http.Request) {
	state = uuid.New().String()

	// Addes required query parameters to the /authorize endpoint.
	authUrl := func() string {
		spotifyauthurl := "https://accounts.spotify.com/authorize"

		scope := "user-read-private user-read-email"
		clientID := os.Getenv("SPOTIFY_ID")
		redirectUri := os.Getenv("REDIRECT_URI")

		query := url.Values{}
		query.Add("client_id", clientID)
		query.Add("response_type", "code")
		query.Add("redirect_uri", redirectUri)
		query.Add("scope", scope)
		query.Add("state", state)

		return fmt.Sprintf("%s?%s", spotifyauthurl, query.Encode())
	}()

	req, err := http.NewRequest(http.MethodGet, authUrl, nil)
	if err != nil {
		err = errors.HTTPRequestError.Wrap(err, "Unable to create new http request")
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.HTTPRequestError.Wrap(err, "Unable to do http request")
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal("Invalid auth url")
	}

	http.Redirect(w, r, authUrl, http.StatusTemporaryRedirect)
}

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
