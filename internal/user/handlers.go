package user

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/dionv/spogo/errors"
	"github.com/google/uuid"
)

func startAuth(w http.ResponseWriter, r *http.Request) {
	state = uuid.NewString()

	authUrl := func() string {
		spotifyauthurl := "https://accounts.spotify.com/authorize"

		scope := "user-read-private user-read-email"

		query := url.Values{}
		query.Set("client_id", clientID)
		query.Set("response_type", "code")
		query.Set("redirect_uri", REDIRECT_URI)
		query.Set("scope", scope)
		query.Set("state", state)

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

	fmt.Fprintln(w, "Login success!")
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	responseState := r.URL.Query().Get("state")
	err := r.URL.Query().Get("error")

	if responseState != state {
		log.Fatal("Invalid state")
	}
	if err != "" {
		log.Fatal("Failed to complete authentication")
	}

	fmt.Fprintln(w, "Login success!")

	ch <- code
}
