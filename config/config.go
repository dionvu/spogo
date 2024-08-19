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
	ACCESSTOKENFILE  = "access-token.json"
	REQUESTTOKENFILE = "refresh-token.json"
	DEVICEFILE       = "device.json"
)

// The struct that holds configuration options from "config.yaml",
// including spotify client information, and all information
// about directory locations.
type Config struct {
	path      string
	cachePath string
	Spotify   Spotify `yaml:"spotify"`
}

// Creates spogo config root directory, "config.yaml",
// and spogo cache directory.
func New() (*Config, error) {
	c := &Config{}

	// Sets the root config path.
	path, err := os.UserConfigDir()
	if err != nil {
		return nil, errors.FileOpen.Wrap(err, "failed to get user's home directory")
	}

	cd, err := os.UserCacheDir()
	if err != nil {
		return nil, errors.FileOpen.Wrap(err, "failed to get user's cache directory")
	}

	// Ensures ".config/spogo" exists.
	c.path = filepath.Join(path, APPNAME)
	if err := os.MkdirAll(c.path, os.ModePerm); err != nil {
		return nil, errors.FileCreate.Wrap(err, fmt.Sprintf("creating file path %v", c.path))
	}

	// Ensures ".cache/spogo" exists.
	c.cachePath = filepath.Join(cd, APPNAME)
	if err := os.MkdirAll(c.cachePath, os.ModePerm); err != nil {
		return nil, errors.FileCreate.Wrap(err, fmt.Sprintf("creating file path %v", c.cachePath))
	}

	// Creates "config.yaml" if it doesn't exist.
	if !c.Exists() {
		if err := c.create(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Loads all config options and client ID & client secret from "config.yaml".
func (c *Config) Load() error {
	file, err := os.Open(c.FilePath())
	if err != nil {
		return errors.FileOpen.Wrap(err, fmt.Sprintf("missing config file: %v", c.FilePath()))
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return errors.FileRead.Wrap(err, fmt.Sprintf("failed to read config file: %v", c.FilePath()))
	}

	data := &struct {
		Spotify struct {
			ClientID     string `yaml:"client_id"`
			ClientSecret string `yaml:"client_secret"`
		} `yaml:"spotify"`
	}{}

	err = yaml.Unmarshal(b, data)
	if err != nil {
		return errors.YAML.Wrap(err, fmt.Sprintf("failed to unmarshal yaml: %v", string(b)))
	}

	c.Spotify.setID(data.Spotify.ClientID)
	c.Spotify.setSecret(data.Spotify.ClientSecret)

	return nil
}

// Creates the "config.yaml" file, assuming the config directory,
// ".config/spogo" (for unix), exists.
func (c *Config) create() error {
	file, err := os.Create(c.FilePath())
	if err != nil {
		return errors.FileCreate.Wrap(err, fmt.Sprintf("creating file %v", c.FilePath()))
	}
	defer file.Close()

	// This couldn't possibily go wrong! (I promise this will be fixed in at least the next 2 years)
	wd, _ := os.Getwd()
	confileFile, _ := os.Open(filepath.Join(wd, templates.DIRECTORY, templates.CONFIGFILE))
	b, _ := io.ReadAll(confileFile)

	_, err = file.WriteString(string(b))
	if err != nil {
		return errors.FileWrite.Wrap(err, fmt.Sprintf("writing to file: %v", file.Name()))
	}

	fmt.Printf("Please enter your spotify client ID & client secret: %v\n", color.YellowString(c.FilePath()))
	os.Exit(0)

	return nil
}

// Returns the config path, ".config/spogo" for unix.
func (c *Config) Path() string {
	return c.path
}

// Returns the config file, ".config/spogo/config.yaml" for unix.
func (c *Config) FilePath() string {
	return filepath.Join(c.Path(), CONFIGFILE)
}

// Returns the config file, ".cache/spogo" for unix.
func (c *Config) CachePath() string {
	return c.cachePath
}

// Returns the config file, ".cache/spogo/device.json" for unix.
func (c *Config) DeviceFile() string {
	return filepath.Join(c.CachePath(), DEVICEFILE)
}

// Returns true if the config file, "config.yaml", exists.
func (c *Config) Exists() bool {
	if _, err := os.ReadFile(c.FilePath()); err != nil {
		return false
	}
	return true
}
