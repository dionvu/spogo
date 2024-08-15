package session

import (
	"github.com/dionv/spogo/config"
	"github.com/dionv/spogo/tokens"
)

type Session struct {
	AccessToken  *tokens.AccessToken
	RefreshToken *tokens.RefreshToken
}

// Creates a new session, loading tokens from respective files, and authenticating.
func New(c *config.Config) (*Session, error) {
	s := &Session{
		AccessToken:  &tokens.AccessToken{},
		RefreshToken: &tokens.RefreshToken{},
	}

	// Loads possible access token and refresh token from respective token files.
	s.AccessToken.Load(c)
	s.RefreshToken.Load(c)

	// Authenticates valid access token, or valid access token and refresh token.
	err := s.Authenticate(c)
	if err != nil {
		return nil, err
	}

	return s, err
}
