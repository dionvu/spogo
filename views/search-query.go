package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	CHAR_LIMIT         = 156
	SEARCH_QUERY_WIDTH = 30
)

type SearchQuery struct {
	textInput textinput.Model
	err       error
}

func NewSearchQuery() SearchQuery {
	ti := textinput.New()
	ti.Placeholder = "What do you want to listen to?"
	ti.Focus()
	ti.CharLimit = CHAR_LIMIT
	ti.Width = SEARCH_QUERY_WIDTH

	return SearchQuery{
		textInput: ti,
		err:       nil,
	}
}

func (sq SearchQuery) Init() tea.Cmd {
	return textinput.Blink
}

func (sq SearchQuery) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return sq, tea.Quit
		}
	}

	sq.textInput, cmd = sq.textInput.Update(msg)
	return sq, cmd
}

func (sq SearchQuery) View() string {
	s := fmt.Sprintf(
		"Search\n\n%s\n\n%s",
		sq.textInput.View(),
		"(esc to quit)",
	) + "\n"

	return s
}

func (sq SearchQuery) Query() string {
	return sq.textInput.Value()
}
