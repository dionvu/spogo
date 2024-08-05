package session

import (
	"path/filepath"

	"github.com/dionv/spogo/internal/config"
	"github.com/dionv/spogo/internal/tokens"
)

type Session struct {
	AccessToken  *tokens.AccessToken
	RefreshToken *tokens.RefreshToken
}

func New(c *config.Config) (*Session, error) {
	s := &Session{
		AccessToken:  &tokens.AccessToken{},
		RefreshToken: &tokens.RefreshToken{},
	}

	// Loads possible access token and refresh token from respective token files.
	s.AccessToken.Load(filepath.Join(c.Path(), config.TOKENSDIRECTORY, config.ACCESSTOKENFILE))
	s.RefreshToken.Load(filepath.Join(c.Path(), config.TOKENSDIRECTORY, config.REQUESTTOKENFILE))

	// Authenticates valid access token, or valid access token and refresh token.
	err := s.Authenticate(c)
	if err != nil {
		return nil, err
	}

	return s, err
}
