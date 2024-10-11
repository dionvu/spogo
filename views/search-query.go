package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	CHAR_LIMIT         = 156
	SEARCH_QUERY_WIDTH = 20
)

type SearchQuery struct {
	Text textinput.Model
	err  error
}

func NewSearchQuery() SearchQuery {
	ti := textinput.New()
	ti.Placeholder = "What's on your mind?"
	ti.Focus()
	ti.CharLimit = CHAR_LIMIT
	ti.Width = SEARCH_QUERY_WIDTH
	ti.Cursor.SetMode(cursor.CursorBlink)

	return SearchQuery{
		Text: ti,
		err:  nil,
	}
}

func (sq SearchQuery) Init() tea.Cmd {
	return textinput.Blink
}

func (sq SearchQuery) Update(msg tea.Msg) (SearchQuery, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return sq, tea.Quit
		}
	}

	sq.Text, cmd = sq.Text.Update(msg)

	return sq, cmd
}

func (sq SearchQuery) View() string {
	s := fmt.Sprintf(
		"Search\n\n%s",
		sq.Text.View(),
	) + "\n"

	return s
}

func (sq SearchQuery) Query() string {
	return sq.Text.Value()
}

func (sq SearchQuery) HideCursor() SearchQuery {
	// sq.Text.Cursor.SetMode(cursor.CursorHide)
	sq.Text.Blur()

	return sq
}
