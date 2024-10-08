package ui

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
	MIN_TERMINAL_HEIGHT = 16
	MIN_TERMINAL_WIDTH  = 42

	MAX_TERMINAL_HEIGHT_SMALL  = 30
	MAX_TERMINAL_WIDTH_SMALL   = 70
	MIN_TERMINAL_HEIGHT_NORMAL = 40
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
			w, h := getTerminalSize()
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
func getTerminalSize() (int, int) {
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

// If the terminal is in minimum dimensions to be considered small.
func (t Terminal) IsSizeSmall() bool {
	return t.Height <= MAX_TERMINAL_HEIGHT_SMALL || t.Width <= MAX_TERMINAL_WIDTH_SMALL
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
