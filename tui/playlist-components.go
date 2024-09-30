package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/spotify"
)

// The detailed infomation about a playlist to
// be displayed when a playlist is selected
// or hovered.
type PlaylistInfo struct {
	Name        PlaylistName
	TotalTracks int
	Owner       string
}

func (pi *PlaylistInfo) Update(playlist *spotify.Playlist) {
	pi.Name = PlaylistName(playlist.Name)
	pi.TotalTracks = playlist.Tracks.Total
	pi.Owner = playlist.Owner.DisplayName
}

// Renders the playlistInfo as a content string.
func (pi PlaylistInfo) Content(term Terminal) Content {
	style := lg.NewStyle().Bold(true).Foreground(lg.Color("#458588"))

	return Join([]string{
		style.Render(pi.Name.Adjust(term).String()),
		"Tracks: " + fmt.Sprint(pi.TotalTracks),
		"Owner: " + fmt.Sprint(pi.Owner),
	}, "\n")
}

type PlaylistName string

// Adjusts the playlist string to fit
// within the terminal if it is too big.
func (pn PlaylistName) Adjust(term Terminal) PlaylistName {
	c := Content(pn)

	return PlaylistName(c.AdjustFit(term.Width))
}

func (pn PlaylistName) String() string {
	return string(pn)
}

type PlaylistList struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewPlaylistListModel(items []list.Item, title string) PlaylistList {
	l := list.New(items, itemDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT)
	l.SetFilteringEnabled(false)
	l.Title = title
	// l.Title = lg.NewStyle().PaddingLeft(2).Render(title)
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	lm := PlaylistList{list: l}

	return lm
}

func (pll PlaylistList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			pll.quitting = true
			return pll, tea.Quit

		case "esc":
			return pll, nil
		}
	}

	pll.list, cmd = pll.list.Update(msg)

	return pll, cmd
}

func (pl PlaylistList) Content() Content {
	return Content(pl.list.View())
}

func (pl PlaylistList) View() string {
	return pl.list.View()
}

func (pl PlaylistList) Init() tea.Cmd {
	return nil
}
