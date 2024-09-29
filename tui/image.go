package ui

import (
	"io"
	"net/http"
	"os"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/dionvu/spogo/errors"
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

// Shorthand for rendering image as ascii with size
// small flags.
func (i *Image) AsciiSmall() Ascii {
	return i.Ascii(AsciiFlagsSmall())
}

// Renders the ascii as a string.
func (a Image) Ascii(flags aic_package.Flags) Ascii {
	ascii, err := aic_package.Convert(a.FilePath, flags)
	if err != nil {
		return ""
	}

	return Ascii(ascii)
}

// Updates the ascii image url, and caches the image if it is not the same.
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
	flags.Colored = true
	flags.Dimensions = []int{40, 20}
	flags.Braille = true
	flags.Threshold = 20
	return flags
}

func AsciiFlagsSmall() aic_package.Flags {
	flags := aic_package.DefaultFlags()
	flags.Colored = true
	flags.Dimensions = []int{24, 12}
	flags.Braille = true
	flags.Threshold = 20
	return flags
}
