package ui

import (
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

// Image is a struct that allows caching of the Image
// and conversion to ascii. It also works with the Content
// type to format the ascii content around a terminal.
type Image struct {
	// The url of the image.
	ImageUrl string

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

// Renders the ascii as a string centered in the given terminal size.
func (a Ascii) Center(terminal Terminal) Ascii {
	return Ascii(Content(a).CenterHorizontal(terminal))
}

func (a Ascii) CenterV(terminal Terminal) Ascii {
	return Ascii(Content(a).CenterVertical(terminal))
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
func (a *Image) UpdateImage(url string) {
	if AsciiNewUrl := url; AsciiNewUrl != a.ImageUrl {
		cacheImage(AsciiNewUrl, a.FilePath)
		a.ImageUrl = AsciiNewUrl
	}
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
	flags.Dimensions = []int{30, 15}
	flags.Braille = true
	flags.Threshold = 20
	return flags
}
