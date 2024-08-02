package tokens

import (
	"os"
	"path/filepath"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/app/config"
)

type RefreshToken struct {
	Token
}

func (tok *RefreshToken) Load(path string) error {
	return tok.load(path, "refresh_token")
}

func (t *RefreshToken) Update(tok string, c *config.Config) error {
	path, err := os.UserConfigDir()
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to get user's config dir")
	}

	path = filepath.Join(c.Path(), config.REQUESTTOKENFILE)

	return t.update(tok, path, "refresh_token")
}
