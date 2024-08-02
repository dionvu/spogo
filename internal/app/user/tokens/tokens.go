package tokens

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dionv/spogo/errors"
)

type Token struct {
	token string
}

// Returns the token as a string
func (t *Token) String() string {
	return t.token
}

// Loads the token from the "config.yaml" file.
func (t *Token) load(path string, key string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to read token file: %v", path))
	}

	var data map[string]string
	err = json.Unmarshal(b, &data)
	if err != nil {
		return errors.JSONError.Wrap(err, "Failed to unmarshal token")
	}

	t.token = data[key]

	return nil
}

func (t *Token) update(tok string, path string, key string) error {
	t.token = tok

	file, err := os.Create(path)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to open token file path: %v", path))
	}
	defer file.Close()

	data := map[string]string{}
	data[key] = t.String()

	body, err := json.Marshal(&data)
	if err != nil {
		return errors.JSONError.Wrap(err, fmt.Sprintf("Failed to marshal token body: %v", &data))
	}

	_, err = file.Write(body)
	if err != nil {
		return errors.FileError.Wrap(err, "Failed to write new token to file: %v", path)
	}

	return nil
}
