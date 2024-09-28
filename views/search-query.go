package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchQuery struct {
	textInput textinput.Model
	err       error
}

func NewSearchQuery() SearchQuery {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

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
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		sq.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (sq SearchQuery) Query() string {
	return sq.textInput.Value()
}
