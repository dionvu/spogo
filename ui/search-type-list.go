package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchTypeListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewSearchTypeListModel(items []list.Item) *SearchTypeListModel {
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetFilteringEnabled(false)
	l.Title = padLines("Select a search type: ", 2)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	lm := &SearchTypeListModel{list: l}

	return lm
}

func (m SearchTypeListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m SearchTypeListModel) View() string {
	return m.list.View()
}

func (m SearchTypeListModel) Init() tea.Cmd {
	return nil
}
