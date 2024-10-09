package views

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/components"
)

type Search struct {
	Input    SearchQuery
	TypeList SearchTypeList
	typeMap  map[list.Item]string

	Session *auth.Session
}

var searchTypes = []string{
	"album",
	"track",
}

func NewSearchView() *Search {
	typeItems := make([]list.Item, len(searchTypes))
	for i, t := range searchTypes {
		typeItems[i] = components.ListItem(t)
	}

	s := &Search{
		Input:    NewSearchQuery(),
		TypeList: NewSearchTypeList(typeItems),
		typeMap: map[list.Item]string{
			components.ListItem("album"): "album",
			components.ListItem("track"): "track",
		},
	}

	return s
}

// Returns the selected search type as a string
func (s *Search) SelectedType() string {
	return s.typeMap[s.TypeList.Selected()]
}
