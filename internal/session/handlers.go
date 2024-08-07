package session

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/session/auth/scopes"
	"github.com/google/uuid"
)

const (
	SPOTIFYAUTHURL = "https://accounts.spotify.com/authorize"
)

// Redirects user to the spotify authentication url and awaits callback.
func startAuth(w http.ResponseWriter, r *http.Request) {
	state = uuid.NewString()

	query := url.Values{}
	query.Set("redirect_uri", REDIRECT_URI)
	query.Set("response_type", "code")
	query.Set("client_id", clientID)
	query.Set("scope", strings.Join([]string{
		scopes.UserReadPrivate,
		scopes.UserReadEmail,
		scopes.UserReadPlaybackState,
		scopes.UserModifyPlaybackState,
	}, " "))
	query.Set("state", state)

	req, err := http.NewRequest(http.MethodGet, SPOTIFYAUTHURL, strings.NewReader(query.Encode()))
	if err != nil {
		log.Fatal(errors.HTTPRequestError.Wrap(err, "Unable to create new http request"))
	}

	if _, err = http.DefaultClient.Do(req); err != nil {
		log.Fatal(errors.HTTPRequestError.Wrap(err, "Unable to do http request"))
	}

	http.Redirect(w, r, fmt.Sprintf("%s?%s", SPOTIFYAUTHURL, query.Encode()), http.StatusTemporaryRedirect)
}

// After user is redirected to the redirect uri, ensures valid state
// and fetches the authentication code.
func completeAuth(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		log.Fatal("Invalid state")
	}

	if r.URL.Query().Get("error") != "" {
		log.Fatal("Failed to complete authentication")
	}

	ch <- r.URL.Query().Get("code")

	fmt.Fprintln(w, "Authentication success!")
}