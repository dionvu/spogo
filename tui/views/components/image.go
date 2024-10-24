package components

import (
	"io"
	"net/http"
	"os"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
)

const (
	FILE_EXTENSION      = ".jpeg"
	ASCII_SMALL_HEIGHT  = 22
	ASCII_SMALL_WIDTH   = ASCII_SMALL_HEIGHT / 2
	ASCII_MEDIUM_HEIGHT = 32
	ASCII_MEDIUM_WIDTH  = ASCII_MEDIUM_HEIGHT / 2
)

// Image is a struct that allows caching of the Image
// and conversion to ascii. It also works with the Content
// type to format the ascii content around a terminal.
type Image struct {
	// The url of the image.
	Url string

	// The cached image's path.
	FilePath string
}

// A string of ascii.
type Ascii string

// Returns the ascii as a string.
func (a Ascii) String() string {
	return string(a)
}

// Returns the ascii as a content string.
func (a Ascii) Content() Content {
	return Content(a)
}

// Shorthand for rendering image as ascii with size
// normal flags.
func (i *Image) AsciiNormal(cfg *config.Config) Ascii {
	if !cfg.Ascii.Enabled {
		return Ascii(InvisibleBarV(ASCII_MEDIUM_HEIGHT/2 - 1).PadLinesLeft(ASCII_MEDIUM_WIDTH))
	}

	return i.Ascii(AsciiFlagsNormal(cfg))
}

// Shorthand for rendering image as ascii with size
// small flags.
func (i *Image) AsciiSmall(cfg *config.Config) Ascii {
	if !cfg.Ascii.Enabled {
		return Ascii(InvisibleBarV(ASCII_SMALL_HEIGHT/2 - 1).PadLinesLeft(ASCII_SMALL_WIDTH))
	}

	return i.Ascii(AsciiFlagsSmall(cfg))
}

// Renders the ascii as a string.
func (a Image) Ascii(flags aic_package.Flags) Ascii {
	ascii, err := aic_package.Convert(a.FilePath, flags)
	if err != nil {
		return ""
	}

	return Ascii(ascii)
}

// Updates the image url, and caches the image if it is not the same.
func (img *Image) Update(url string) {
	if AsciiNewUrl := url; AsciiNewUrl != img.Url {
		img.Url = AsciiNewUrl

		err := img.Cache()
		if err != nil {
			errors.Log(errors.PlayerViewImageCache.Wrap(err, "failed to cache new image with url: %s", url))
		}
	}
}

// Caches the image.
func (img *Image) Cache() error {
	res, err := http.Get(img.Url)
	if err != nil {
		return err
	}

	file, err := os.Create(img.FilePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func AsciiFlagsNormal(cfg *config.Config) aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{ASCII_MEDIUM_HEIGHT, ASCII_MEDIUM_WIDTH}
	flags.Threshold = cfg.Ascii.Threshold

	if cfg.Ascii.Grayscale {
		flags.Grayscale = true
	} else {
		flags.Colored = true
	}

	flags.Braille = true

	return flags
}

func AsciiFlagsSmall(cfg *config.Config) aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{ASCII_SMALL_HEIGHT, ASCII_SMALL_WIDTH}

	flags.Threshold = cfg.Ascii.Threshold

	if cfg.Ascii.Grayscale {
		flags.Grayscale = true
	} else {
		flags.Colored = true
	}

	flags.Braille = true

	return flags
}
