package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DeviceListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewDeviceListModel(items []list.Item) *DeviceListModel {
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetFilteringEnabled(false)
	l.Title = padLines("Select a search type: ", 2)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	lm := &DeviceListModel{list: l}

	return lm
}

func (m DeviceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (dlm DeviceListModel) View() string {
	return dlm.list.View()
}

func (_ DeviceListModel) Init() tea.Cmd {
	return nil
}
