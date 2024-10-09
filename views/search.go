package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/components"
	"github.com/dionvu/spogo/spotify"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Search struct {
	Input    SearchQuery
	TypeList SearchTypeList
	typeMap  map[list.Item]string
	Results  Results

	session *auth.Session
}

var searchTypes = []string{
	"album",
	"track",
}

func NewSearch(session *auth.Session) Search {
	typeItems := make([]list.Item, len(searchTypes))
	for i, t := range searchTypes {
		typeItems[i] = components.ListItem(t)
	}

	s := Search{
		Input:    NewSearchQuery(),
		TypeList: NewSearchTypeList(typeItems),
		typeMap: map[list.Item]string{
			components.ListItem("album"): "album",
			components.ListItem("track"): "track",
		},

		Results: Results{},
	}

	// s.Results = s.Results.Refresh("spotify", session)

	s.Results.listTracks = components.NewDefaultList([]list.Item{components.ListItem("")}, "Tracks")

	return s
}

func (s Search) View(term components.Terminal) string {
	l := table.NewWriter()
	l.Style().Options.DrawBorder = false
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	l.AppendRows([]table.Row{
		{components.Content(s.Input.View()).PadLinesLeft(2).String()},
		{s.TypeList.View()},
	})

	t.AppendRows([]table.Row{
		{
			l.Render(),
			components.Content(s.Results.View()).PadLinesLeft(0).String() + components.Content("\n-").Append('-', 40).String(),
		},
	})

	c := components.Content(t.Render()) + "\n" + components.Content("-").Append('-', 100)

	return c.CenterVertical(term).CenterHorizontal(term).String()
}

// Returns the selected search type as a string
func (s Search) SelectedType() string {
	return s.typeMap[s.TypeList.Selected()]
}

type Results struct {
	CurrentType string
	Items       spotify.SearchResult
	listTracks  list.Model
	listAlbums  list.Model
}

func (r Results) Refresh(query string, s *auth.Session) Results {
	results, _ := spotify.Search(query, searchTypes, s)

	items := make([]list.Item, len(results.Tracks))
	for i, track := range results.Tracks {
		items[i] = components.ListItem(components.Content(track.Name).AdjustFit(35))
	}
	r.listTracks = components.NewDefaultList(items, "Tracks")
	r.listAlbums = components.NewDefaultList(items, "Albums")

	return r
}

func (r Results) Update(msg tea.Msg) (Results, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return r, tea.Quit

		case "esc":
			return r, nil
		}
	}

	r.listAlbums, cmd = r.listAlbums.Update(msg)
	r.listTracks, cmd = r.listTracks.Update(msg)

	return r, cmd
}

func (r Results) View() string {
	return r.listTracks.View()
}

func (r Results) Init() tea.Cmd {
	return nil
}
