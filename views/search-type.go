package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/components"
)

type SearchTypeList struct {
	list     list.Model
	choice   string
	quitting bool
}

// The selected type as a list item.
func (stl SearchTypeList) Selected() list.Item {
	return stl.list.SelectedItem()
}

func NewSearchTypeList(items []list.Item) SearchTypeList {
	l := components.NewDefaultList(items, "Select a search type: ")

	lm := SearchTypeList{list: l}

	return lm
}

func (m SearchTypeList) Update(msg tea.Msg) (SearchTypeList, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m SearchTypeList) View() string {
	return m.list.View()
}

func (m SearchTypeList) Init() tea.Cmd {
	return nil
}
