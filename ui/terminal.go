package ui

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	MIN_TERMINAL_HEIGHT = 21
	MIN_TERMINAL_WIDTH  = 42

	MAX_TERMINAL_HEIGHT_SMALL  = 30
	MIN_TERMINAL_HEIGHT_NORMAL = 40
)

func (t Terminal) IsValid() bool {
	return t.Height >= MIN_TERMINAL_HEIGHT && t.Width >= MIN_TERMINAL_WIDTH
}

func (t Terminal) IsSizeSmall() bool {
	return t.Height <= MAX_TERMINAL_HEIGHT_SMALL
}

func (t Terminal) IsSizeNormal() bool {
	return t.Height >= MIN_TERMINAL_HEIGHT_NORMAL
}

// Returns the error message associated with the terminal being
// below the required dimensions.
func terminalWarningView(terminal Terminal) string {
	return color.RedString(
		fmt.Sprint(
			"Terminal of size ",
			terminal.Height, "x", terminal.Width,
			" is prone to visual glitches.\nMinimum required height is ",
			MIN_TERMINAL_HEIGHT, "x", MIN_TERMINAL_WIDTH, "."))
}
