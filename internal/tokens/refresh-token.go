package tokens

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/internal/config"
)

type RefreshToken struct {
	Token string
}

func NewRefreshToken(tok string) *RefreshToken {
	t := &RefreshToken{
		Token: tok,
	}

	return t
}

// Returns the token as a string
func (t *RefreshToken) String() string {
	return t.Token
}

// Loads the token from the token file.
func (t *RefreshToken) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to read token file: %v", path))
	}

	// var data map[string]string
	err = json.Unmarshal(b, t)
	if err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal token")
	}

	return nil
}

// Updates the token file with new token.
func (t *RefreshToken) Update(tok string, c *config.Config) error {
	path, err := os.UserConfigDir()
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to get user's config dir")
	}

	path = filepath.Join(c.Path(), config.TOKENSDIRECTORY)

	os.MkdirAll(path, os.ModePerm)

	path = filepath.Join(path, config.REQUESTTOKENFILE)

	return t.saveToFile(path)
}

func (t *RefreshToken) saveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	body, err := json.Marshal(t)
	if err != nil {
		return errors.JSONError.Wrap(err, fmt.Sprintf("Failed to marshal token body: %v", *t))
	}

	_, err = file.Write(body)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to write new token to file: %v", path)
	}

	return nil
}
