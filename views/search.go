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
	// Used to map the selected search type item to the search type as a string.
	typeMap map[list.Item]string
	Results Results

	session *auth.Session
}

var searchTypes = []string{
	"album",
	"track",
}

var typeMap = map[list.Item]string{
	components.ListItem("album"): "album",
	components.ListItem("track"): "track",
}

// Creates a new search view.
func NewSearch(session *auth.Session) Search {
	typeItems := make([]list.Item, len(searchTypes))
	for i, t := range searchTypes {
		typeItems[i] = components.ListItem(t)
	}

	s := Search{
		Input:    NewSearchQuery(),
		TypeList: NewSearchTypeList(typeItems),
		typeMap:  typeMap,
		Results:  Results{},
	}

	s.Results.listTracks = components.NewDefaultList([]list.Item{components.ListItem("")}, "Tracks")

	return s
}

// Renders the search view, this includes, the text area,
// the type selection, and the list of results.
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
			components.Content(s.Results.view()).PadLinesLeft(0).String() + components.Content("\n-").Append('-', 40).String(),
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

// Called whenever the user has finished inputing a search query and selected the search type
// of the results to be displayed. This updates the state of result to match the desired
// specified content.
func (r Results) Refresh(query string, currentSelectedType string, s *auth.Session) Results {
	results, _ := spotify.Search(query, searchTypes, s)

	r.CurrentType = currentSelectedType

	switch r.CurrentType {
	case "track":
		items := make([]list.Item, len(results.Tracks))
		for i, track := range results.Tracks {
			items[i] = components.ListItem(components.Content(track.Name).AdjustFit(35))
		}
		r.listTracks = components.NewDefaultList(items, "Tracks")

	case "album":
		items := make([]list.Item, len(results.Albums))
		for i, album := range results.Albums {
			items[i] = components.ListItem(components.Content(album.Name).AdjustFit(35))
		}
		r.listAlbums = components.NewDefaultList(items, "Albums")
	}

	return r
}

// Updates the state of the result when a key is pressed,
// handling every search type's result list. Ensure that
// result.CurrentType has been set before (by calling
// result.Refresh() ) before Update is called.
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

	switch r.CurrentType {
	case "track":
		r.listTracks, cmd = r.listTracks.Update(msg)
	case "album":
		r.listAlbums, cmd = r.listAlbums.Update(msg)
	}

	return r, cmd
}

// Renders the result view based on the
// current selected display type. Ensure
// result.Refresh() was once before we display
// the result content.
func (r Results) view() string {
	switch r.CurrentType {
	case "track":
		return r.listTracks.View()
	case "album":
		return r.listAlbums.View()
	}

	return "unknown search type"
}

func (r Results) init() tea.Cmd {
	return nil
}
