package player

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/errors"
)

// Helper function for player new function to
// see if the "device.json" cache file exists.
func deviceCacheExist(c *config.Config) bool {
	if _, err := os.ReadFile(filepath.Join(c.CachePath(), config.DEVICEFILE)); err != nil {
		return false
	}
	return true
}

// Helper function for the player new function
// to creates the "device.json" cache file.
func createCache(c *config.Config) error {
	file, err := os.Create(filepath.Join(c.CachePath(), config.DEVICEFILE))
	if err != nil {
		return errors.FileCreate.Wrap(err, fmt.Sprintf("creating file %v", c.FilePath()))
	}
	file.Close()

	return nil
}
