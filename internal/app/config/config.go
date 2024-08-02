package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dionv/spogo/errors"
	"gopkg.in/yaml.v3"
)

const (
	APPNAME          = "spogo"
	CONFIGFILE       = "config.yaml"
	ACCESSTOKENFILE  = "tokens/access-token.json"
	REQUESTTOKENFILE = "tokens/refresh-token.json"
)

type Config struct {
	path    string
	Spotify Spotify `yaml:"spotify"`
}

func New() (*Config, error) {
	c := &Config{}

	path, err := os.UserConfigDir()
	if err != nil {
		return nil, errors.FileError.Wrap(err, "Failed to get user's home directory")
	}

	c.path = filepath.Join(path, APPNAME)

	return c, nil
}

// Creates "spogo/config.yaml" in user's config dir.
func (c *Config) Create() error {
	cp := c.Path()

	if err := os.MkdirAll(cp, os.ModePerm); err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Creating file path %v", cp))
	}

	cfp := filepath.Join(cp, CONFIGFILE)

	file, err := os.Create(cfp)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Creating file %v", cfp))
	}
	defer file.Close()

	content := "spotify:\n  client_id: \"\"\n  client_secret: \"\""

	_, err = file.WriteString(content)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Writing to file: %v", content))
	}

	return nil
}

// Returns true if config file exists.
func (c *Config) Exists() (bool, error) {
	cfp := c.FilePath()

	if _, err := os.ReadFile(cfp); err != nil {
		return false, nil
	}

	return true, nil
}

// Loads "spogo/config.yaml" from user's config dir.
func (c *Config) Load() error {
	cfp := c.FilePath()

	file, err := os.Open(cfp)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Missing config file: %v", cfp))
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Failed to read config file: %v", cfp))
	}

	data := &struct {
		Spotify struct {
			ClientID     string `yaml:"client_id"`
			ClientSecret string `yaml:"client_secret"`
		} `yaml:"spotify"`
	}{}

	err = yaml.Unmarshal(buf, data)
	if err != nil {
		return errors.YAMLError.Wrap(err, fmt.Sprintf("Unmarshal failed"))
	}

	c.Spotify.setID(data.Spotify.ClientID)
	c.Spotify.setSecret(data.Spotify.ClientSecret)

	return nil
}

// Returns config root path.
func (c *Config) Path() string {
	return c.path
}

// Returns config path including the config file.
func (c *Config) FilePath() string {
	return filepath.Join(c.Path(), CONFIGFILE)
}
