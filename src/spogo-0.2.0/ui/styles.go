package ui

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

const TAB_WIDTH = 4

var CommonStyle = struct {
	MainControls struct {
		Selected lg.Style
		Normal   lg.Style
	}
}{
	MainControls: struct {
		Selected lg.Style
		Normal   lg.Style
	}{
		Normal:   lg.NewStyle().Faint(true),
		Selected: lg.NewStyle(),
	},
}

var PlayerViewStyle = struct {
	StatusBar struct {
		NowPlaying lg.Style
		Paused     lg.Style
		NoPlayer   lg.Style
	}
}{
	StatusBar: struct {
		NowPlaying lg.Style
		Paused     lg.Style
		NoPlayer   lg.Style
	}{
		NowPlaying: lg.NewStyle().
			Bold(true).
			Foreground(lg.Color("#282828")).
			Background(lg.Color("#98971a")).
			PaddingLeft(1).
			PaddingRight(1),

		Paused: lg.NewStyle().
			Bold(true).
			Foreground(lg.Color("#282828")).
			Background(lg.Color("#d79921")).
			PaddingLeft(1).
			PaddingRight(1),

		NoPlayer: lg.NewStyle().
			Bold(true).
			Foreground(lg.Color("#282828")).
			Background(lg.Color("#cc241d")).
			PaddingLeft(1).
			PaddingRight(1),
	},
}

var PlaylistViewStyle = struct {
	Title        lg.Style
	ItemSelected lg.Style
	Item         lg.Style
}{
	Title:        lg.NewStyle().Bold(true).Background(lg.Color("#a89984")).Foreground(lg.Color("#282828")).PaddingLeft(1).PaddingRight(1),
	ItemSelected: lg.NewStyle().PaddingLeft(2),

	Item: lg.NewStyle().PaddingLeft(4).Faint(true),
}

var DeviceViewStyle = struct {
	Title        lg.Style
	ItemSelected lg.Style
	Item         lg.Style
}{
	Title:        lg.NewStyle().Bold(true).Background(lg.Color("#a89984")).Foreground(lg.Color("#282828")).PaddingLeft(1).PaddingRight(1),
	ItemSelected: lg.NewStyle().PaddingLeft(2),

	Item: lg.NewStyle().PaddingLeft(4).Faint(true),
}

func padLines(s string, padding int) string {
	pad := strings.Repeat(" ", padding)
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		lines[i] = pad + line
	}

	return strings.Join(lines, "\n")
}
