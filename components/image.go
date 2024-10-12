package components

import (
	"io"
	"net/http"
	"os"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/dionvu/spogo/errors"
)

const (
	ASCII_SMALL_HEIGHT  = 22
	ASCII_SMALL_WIDTH   = 11
	ASCII_MEDIUM_HEIGHT = 34
	ASCII_MEDIUM_WIDTH  = 16
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
func (i *Image) AsciiNormal() Ascii {
	return i.Ascii(AsciiFlagsNormal())
}

func (i *Image) AsciiNormalBW() Ascii {
	return i.Ascii(AsciiFlagsNormalBW())
}

// Shorthand for rendering image as ascii with size
// small flags.
func (i *Image) AsciiSmall() Ascii {
	return i.Ascii(AsciiFlagsSmall())
}

func (i *Image) AsciiSmallBW() Ascii {
	return i.Ascii(AsciiFlagsSmallBW())
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
			errors.LogError(errors.PlayerViewImageCache.Wrap(err, "failed to cache new image with url: %s", url))
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

func AsciiFlagsNormal() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{ASCII_MEDIUM_HEIGHT, ASCII_MEDIUM_WIDTH}
	flags.Colored = true
	flags.Braille = true
	flags.Threshold = 20
	return flags
}

func AsciiFlagsNormalBW() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	// flags.Colored = false
	flags.Grayscale = true
	flags.Dimensions = []int{ASCII_MEDIUM_HEIGHT, ASCII_MEDIUM_WIDTH}
	flags.Braille = true
	flags.Threshold = 80
	return flags
}

func AsciiFlagsSmall() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Colored = true
	flags.Dimensions = []int{ASCII_SMALL_HEIGHT, ASCII_SMALL_WIDTH}
	flags.Braille = true
	flags.Threshold = 0
	return flags
}

func AsciiFlagsSmallBW() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	// flags.Colored = false
	flags.Grayscale = true
	flags.Dimensions = []int{ASCII_SMALL_HEIGHT, ASCII_SMALL_WIDTH}
	flags.Braille = true
	flags.Threshold = 80
	return flags
}
