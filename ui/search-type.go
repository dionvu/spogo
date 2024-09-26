package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/auth"
)

type SearchTypeView struct {
	Session   *auth.Session
	ListModel *SearchTypeListModel
	Types     []string
}

type SearchTypeListModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewSearchTypeView(s *auth.Session) *SearchTypeView {
	types := []string{
		"album",
		"track",
	}

	items := []list.Item{
		Item("album"),
		Item("track"),
	}

	stv := SearchTypeView{
		Session:   s,
		ListModel: NewSearchTypeListModel(items),
		Types:     types,
	}

	if len(stv.ListModel.list.Items()) > 0 {
		stv.ListModel.choice = stv.Types[0]
	}

	return &stv
}

func (st *SearchTypeView) View(playerView *PlayerView, terminal Terminal) string {
	mainControls := MainControlsRender(SEARCH_TYPE_VIEW)

	if terminal.Height < TERMINALSIZE.Small {
		return "\n\n" + st.ListModel.View()
	}

	return "\n\n" + mainControls + "\n\n" + st.ListModel.View()
}

func NewSearchTypeListModel(items []list.Item) *SearchTypeListModel {
	l := list.New(items, itemDelegate{}, DEFAULT_WIDTH, LIST_HEIGHT)
	l.SetFilteringEnabled(false)
	l.Title = padLines("Select a search type: ", 2)
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
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
