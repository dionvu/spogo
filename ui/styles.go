package ui

import lg "github.com/charmbracelet/lipgloss"

var NOW_PLAYING_STYLE lg.Style = lg.NewStyle().Bold(true).Foreground(lg.Color("#282828")).
	Background(lg.Color("#98971a")).PaddingLeft(1).PaddingRight(1)

var PAUSED_STYLE lg.Style = lg.NewStyle().Bold(true).Foreground(lg.Color("#282828")).
	Background(lg.Color("#d79921")).PaddingLeft(1).PaddingRight(1)

var NO_PLAYER_STYLE lg.Style = lg.NewStyle().Bold(true).Foreground(lg.Color("#282828")).
	Background(lg.Color("#cc241d")).PaddingLeft(1).PaddingRight(1)

var MAIN_CONTROLS_SELECTED_STYLE lg.Style = lg.NewStyle()

var MAIN_CONTROLS_STYLE lg.Style = lg.NewStyle().Faint(true)

var TITLE_PLAYLIST_STYLE lg.Style = lg.NewStyle().Bold(true).Background(lg.Color("#d65d0e")).Foreground(lg.Color("#282828")).PaddingLeft(1).PaddingRight(1)

var ITEM_SELECTED_STYLE lg.Style = lg.NewStyle().Bold(true).Foreground(lg.Color("#d65d0e"))
