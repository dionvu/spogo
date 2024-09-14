package ui

import lg "github.com/charmbracelet/lipgloss"

var NOWPLAYINGSTYLE lg.Style = lg.NewStyle().Bold(true).Foreground(lg.Color("#282828")).
	Background(lg.Color("#98971a")).PaddingLeft(1).PaddingRight(1)

var PAUSEDSTYLE lg.Style = lg.NewStyle().Bold(true).Foreground(lg.Color("#282828")).
	Background(lg.Color("#d79921")).PaddingLeft(1).PaddingRight(1)
