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
	APP_NAME         = "spogo"
	CONFIG_FILE_NAME = "config.yaml"
)

type Config struct {
	Spotify struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"spotify"`
}

func New() *Config {
	return &Config{}
}

// Loads "spogo/config.yaml" from user's config dir.
func (c *Config) Load() error {
	cfp, e := c.FilePath()
	if e != nil {
		return e
	}

	file, e := os.Open(cfp)

	if e != nil {
		return errors.FileError.Wrap(e, fmt.Sprintf("Missing config file: %v", cfp))
	}

	defer file.Close()

	buf, e := io.ReadAll(file)

	if e != nil {
		return errors.FileError.Wrap(e, fmt.Sprintf("Failed to read config file: %v", cfp))
	}

	e = yaml.Unmarshal(buf, c)

	if e != nil {
		return errors.YAMLError.Wrap(e, fmt.Sprintf("Unmarshal failed"))
	}

	return nil
}

// Creates "spogo/config.yaml" in user's config dir.
func (c *Config) Create() error {
	cp, e := c.Root()
	if e != nil {
		return e
	}

	if e := os.MkdirAll(cp, os.ModePerm); e != nil {
		return errors.FileError.Wrap(e, fmt.Sprintf("Creating file path %v", cp))
	}

	cfp := filepath.Join(cp, CONFIG_FILE_NAME)

	file, e := os.Create(cfp)

	if e != nil {
		return errors.FileError.Wrap(e, fmt.Sprintf("Creating file %v", cfp))
	}

	defer file.Close()

	content := "spotify:\n  client_id: \"\"\n  client_secret: \"\""

	_, e = file.WriteString(content)
	if e != nil {
		return errors.FileError.Wrap(e, fmt.Sprintf("Writing to file: %v", content))
	}

	return nil
}

// Returns true if config file exists.
func (c *Config) Exists() (bool, error) {
	cfp, e := c.FilePath()

	if e != nil {
		return false, e
	}

	if _, e := os.ReadFile(cfp); e != nil {
		return false, nil
	}

	return true, nil
}

// Returns config root path.
func (c *Config) Root() (string, error) {
	ucd, e := os.UserConfigDir()

	if e != nil {
		return "", errors.FileError.Wrap(e, "Failed to get user's home directory")
	}

	return filepath.Join(ucd, APP_NAME), nil
}

// Returns config path including the config file.
func (c *Config) FilePath() (string, error) {
	ucd, e := os.UserConfigDir()

	if e != nil {
		return "", errors.FileError.Wrap(e, "Failed to get user's home directory")
	}

	return filepath.Join(ucd, APP_NAME, CONFIG_FILE_NAME), nil
}
