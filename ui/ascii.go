package ui

import "github.com/TheZoraiz/ascii-image-converter/aic_package"

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

func AsciiRender(filepath string, flags aic_package.Flags) (string, error) {
	ascii, err := aic_package.Convert(filepath, flags)
	if err != nil {
		return "", err
	}

	return ascii, nil
}
