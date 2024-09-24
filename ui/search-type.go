package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/dionvu/spogo/auth"
)

type SearchTypeView struct {
	Session   *auth.Session
	ListModel *SearchTypeListModel
	Types     []string
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
	mainControls := MainControlsView(SEARCH_TYPE_VIEW)

	if terminal.Height < TERMINALSIZE.Small {
		return "\n\n" + st.ListModel.View()
	}

	return "\n\n" + mainControls + "\n\n" + st.ListModel.View()
}
