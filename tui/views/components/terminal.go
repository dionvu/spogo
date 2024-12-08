package components

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/term"
)

const (
	MIN_TERMINAL_HEIGHT            = 22
	MIN_TERMINAL_WIDTH             = 42
	MAX_TERMINAL_HEIGHT_VERY_SMALL = 22
	MAX_TERMINAL_HEIGHT_SMALL      = 30
	MAX_TERMINAL_WIDTH_SMALL       = 76
	MIN_TERMINAL_HEIGHT_NORMAL     = 40
)

type Terminal struct {
	Height int
	Width  int
}

// Asyncronously updates the terminal dimensions.
func (terminal *Terminal) UpdateSize() {
	// Channel to receive terminal size change signals (SIGWINCH)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)

	// If their is a change in terminal dimensions, updates terminal.
	go func() {
		for range sigCh {
			w, h := GetTerminalSize()
			if w != terminal.Width || h != terminal.Height {
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()

			}

			terminal.Width, terminal.Height = w, h
		}
	}()
}

// Gets the current dimensions of the user's terminal.
func GetTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return -1, -1
	}

	return width, height
}

// If the terminal is within the minimum dimensions.
func (t Terminal) IsValid() bool {
	return t.Height >= MIN_TERMINAL_HEIGHT && t.Width >= MIN_TERMINAL_WIDTH
}

// If the terminal height is in within minimum dimensions
// to be considered small.
func (t Terminal) HeightIsSmall() bool {
	return MAX_TERMINAL_HEIGHT_VERY_SMALL < t.Height && t.Height <= MAX_TERMINAL_HEIGHT_SMALL
}

func (t Terminal) HeightIsVerySmall() bool {
	return MAX_TERMINAL_HEIGHT_VERY_SMALL >= t.Height
}

// If the terminal height is in within minimum dimensions
// to be considered small.
func (t Terminal) WidthIsSmall() bool {
	return t.Width <= MAX_TERMINAL_WIDTH_SMALL
}

// If the terminal exceeds the minimum dimensions to be considered normal.
func (t Terminal) IsSizeNormal() bool {
	return t.Height >= MIN_TERMINAL_HEIGHT_NORMAL
}

// Returns the error message associated with the terminal being
// below the required dimensions.
func (terminal *Terminal) WarningString() string {
	return color.RedString(
		fmt.Sprint(
			"Terminal of size ",
			terminal.Height, "x", terminal.Width,
			" is prone to visual glitches.\nMinimum required height is ",
			MIN_TERMINAL_HEIGHT, "x", MIN_TERMINAL_WIDTH, "."))
}
