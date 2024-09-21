package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listHeight          = 10
	defaultWidth        = 40
	DEFAULT_LIST_HEIGHT = 20 // 8 Selections per page
	SMALL_LIST_HEIGHT   = 9  // 5 Selections per page
)

type PlaylistListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewPlaylistListModel(items []list.Item, title string) *PlaylistListModel {
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetFilteringEnabled(false)
	l.Title = padLines(title, 2)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	lm := &PlaylistListModel{list: l}

	return lm
}

func (m PlaylistListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m PlaylistListModel) View() string {
	return m.list.View()
}

func (m PlaylistListModel) Init() tea.Cmd {
	return nil
}
