package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
)

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

func NewRefreshToken(tok string) *RefreshToken {
	t := &RefreshToken{
		Token: tok,
	}

	return t
}

// Loads the token fields from the refresh token file.
func (t *RefreshToken) Load(c *config.Config) error {
	path := filepath.Join(c.CachePath(), config.REQUESTTOKENFILE)

	file, err := os.Open(path)
	if err != nil {
		err = errors.FileOpen.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
		errors.LogError(err)
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		err = errors.FileRead.Wrap(err, fmt.Sprintf("Failed to read token file: %v", path))
		errors.LogError(err)
		return err
	}

	err = json.Unmarshal(b, t)
	if err != nil {
		err = errors.JSONUnmarshal.Wrap(err, "Failed to unmarshal token from file body: %v", string(b))
		errors.LogError(err)
		return err
	}

	return nil
}

// Updates the token file with new token.
func (t *RefreshToken) Update(tok string, c *config.Config) error {
	t.Token = tok

	filePath := filepath.Join(c.CachePath(), config.REQUESTTOKENFILE)
	file, err := os.Create(filePath)
	if err != nil {
		err = errors.FileCreate.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", filePath))
		errors.LogError(err)
		return err
	}
	defer file.Close()

	b, err := json.Marshal(t)
	if err != nil {
		err = errors.JSONMarshal.Wrap(err, fmt.Sprintf("Failed to marshal token body: %v", *t))
		errors.LogError(err)
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		err = errors.FileWrite.Wrap(err, fmt.Sprintf("Failed to write new token to file: %v", filePath))
		errors.LogError(err)
		return err
	}

	return nil
}

// The token as a string
func (t *RefreshToken) String() string {
	return t.Token
}
