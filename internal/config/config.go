package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dionv/spogo/errors"
	"github.com/dionv/spogo/templates"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

const (
	APPNAME          = "spogo"
	CONFIGFILE       = "config.yaml"
	TOKENSDIRECTORY  = ".tokens"
	ACCESSTOKENFILE  = "access-token.json"
	REQUESTTOKENFILE = "refresh-token.json"
)

// The struct that holds configuration options from "config.yaml",
// including spotify client information.
type Config struct {
	path    string
	Spotify Spotify `yaml:"spotify"`
}

// Creates app root directory, and "config.yaml".
func New() (*Config, error) {
	c := &Config{}

	// Sets the root config path.
	path, err := os.UserConfigDir()
	if err != nil {
		return nil, errors.FileError.Wrap(err, "Failed to get user's home directory")
	}

	c.path = filepath.Join(path, APPNAME)

	// Creates the app root dir and "config.yaml"
	configExists, err := c.Exists()
	if err != nil {
		return c, err
	}

	if !configExists {
		if err := c.create(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Loads all config options and client ID & client secret from "config.yaml".
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

// Creates the root config directory and "config.yaml".
func (c *Config) create() error {
	if err := os.MkdirAll(c.Path(), os.ModePerm); err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Creating file path %v", c.Path()))
	}

	file, err := os.Create(c.FilePath())
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Creating file %v", c.FilePath()))
	}
	defer file.Close()

	// This couldn't possibily go wrong!
	wd, _ := os.Getwd()
	confileFile, _ := os.Open(filepath.Join(wd, templates.DIRECTORY, templates.CONFIGFILE))
	b, _ := io.ReadAll(confileFile)

	_, err = file.WriteString(string(b))
	if err != nil {
		return errors.FileError.Wrap(err, fmt.Sprintf("Writing to file: %v", file.Name()))
	}

	fmt.Printf("Please enter your spotify client ID & client secret: %v\n", color.YellowString(c.FilePath()))
	os.Exit(0)

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

// Returns true if config file exists.
func (c *Config) Exists() (bool, error) {
	if _, err := os.ReadFile(c.FilePath()); err != nil {
		return false, nil
	}

	return true, nil
}
