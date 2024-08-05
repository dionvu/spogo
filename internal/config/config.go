package config

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dionv/spogo/errors"
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
	configExists, err := c.exists()
	if err != nil {
		return c, err
	}

	if !configExists {
		err := c.create()
		if err != nil {
			return c, err
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

// Attempts to do the "client credentials" authentication flow
// to verify valid spotify client ID and client secret.
func (c *Config) ValidSpotifyCredentials() (bool, error) {
	spotifyUrl := "https://accounts.spotify.com/api/token"
	id := c.Spotify.ClientID()
	secret := c.Spotify.ClientSecret()

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, spotifyUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("unable to create new http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	encodedImportantStuff := base64.StdEncoding.EncodeToString([]byte(id + ":" + secret))
	req.Header.Set("Authorization", "Basic "+encodedImportantStuff)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, errors.HTTPError.Wrap(err, fmt.Sprintf("unable to do http request: %v", err))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, errors.HTTPError.Wrap(err, "invalid client ID or secret")
	}

	return true, nil
}

// Returns config root path.
func (c *Config) Path() string {
	return c.path
}

// Returns config path including the config file.
func (c *Config) FilePath() string {
	return filepath.Join(c.Path(), CONFIGFILE)
}

// Creates the root config directory and "config.yaml".
func (c *Config) create() error {
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

	fmt.Printf("Please enter your spotify client ID & client secret: %v\n", color.GreenString(cfp))
	os.Exit(0)

	return nil
}

// Returns true if config file exists.
func (c *Config) exists() (bool, error) {
	cfp := c.FilePath()

	if _, err := os.ReadFile(cfp); err != nil {
		return false, nil
	}

	return true, nil
}
