package ui

import (
	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

type Ascii struct {
	ImageUrl string
	FilePath string
}

// Renders the ascii as a string.
func (a *Ascii) Render(flags aic_package.Flags) string {
	ascii, err := aic_package.Convert(a.FilePath, flags)
	if err != nil {
		return ""
	}

	return ascii
}

// Renders the ascii as a string centered in the given terminal size.
func (a *Ascii) Center(flags aic_package.Flags, terminal Terminal) string {
	return CenterString(a.Render(flags), terminal)
}

// Updates the ascii image url, and caches the image if it is not the same.
func (a *Ascii) UpdateImage(url string) {
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
