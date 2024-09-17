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

var TITLE_PLAYLIST_STYLE lg.Style = lg.NewStyle().Bold(true).Background(lg.Color("#d65d0e")).Foreground(lg.Color("#282828")).PaddingLeft(1).PaddingRight(1)

var ITEM_SELECTED_STYLE lg.Style = lg.NewStyle().Bold(true).Foreground(lg.Color("#d65d0e"))

func padLines(s string, padding int) string {
	pad := strings.Repeat(" ", padding)
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		lines[i] = pad + line
	}

	return strings.Join(lines, "\n")
}
