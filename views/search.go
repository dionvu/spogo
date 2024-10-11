package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/components"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	SEARCH_VIEW_WIDTH    = 86
	SEARCH_RESULTS_WIDTH = 40
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
func (s Search) View(term components.Terminal, curView string) string {
	if curView == SEARCH_VIEW_TYPE {
		style := lipgloss.NewStyle().Underline(true)
		s.TypeList.list.Title = style.Render("Select a search type:")
	} else {
		s.TypeList.list.Title = "Select a search type:"
	}

	left := components.NewDefaultTable()
	container := components.NewDefaultTable()

	left.AppendRows([]table.Row{
		{components.Content(s.Input.View()).PadLinesLeft(2).String()},
		{s.TypeList.View()},
	})

	container.AppendRows([]table.Row{
		{
			// Offsets the left to align center, as the right requires more
			// allocated space for ~40 character result names.
			components.Content(left.Render()).PadLinesLeft(16).String(),

			components.Join([]string{
				s.Results.view(),
				components.Content("").Append(' ', SEARCH_RESULTS_WIDTH).String(),
			}, "\n"),
		},
	})

	extraInfo := func() string {
		switch s.Results.CurrentType {
		case "track":
			mins, secs := utils.MsToMinutesAndSeconds(s.Results.SelectedTrack().DurationMs)
			return components.Join(
				[]string{
					s.Results.SelectedTrack().Artists[0].Name,
					mins + "m:" + secs + "s",
				}, "\n\n").Prepend('\n', 0).String()

		case "album":
			return components.Join(
				[]string{
					s.Results.SelectedAlbum().Artists[0].Name,
					" Tracks: " + fmt.Sprint(s.Results.SelectedAlbum().TotalTracks),
				}, "\n\n").Prepend('\n', 0).String()

		default:
			return "\n\n\n"
		}
	}()

	vs := ViewStatus{}
	vs.Update(SEARCH_VIEW_RESULTS)

	c := components.Join([]components.Content{
		components.Content("").Append('\n', 6),
		components.Content(container.Render()).Append('\n', 1),
		components.Content(extraInfo),
		components.Content(" ").Append(' ', SEARCH_VIEW_WIDTH),
		components.Content("").Append('\n', 0),
		vs.Content(),
	}, "\n")

	return c.CenterVertical(term).CenterHorizontal(term).String()
}

// Returns the selected search type as a string
func (s Search) SelectedType() string {
	return s.typeMap[s.TypeList.Selected()]
}

type Results struct {
	CurrentType string
	Items       spotify.SearchResult

	listTracks list.Model
	trackMap   map[list.Item]*spotify.Track

	listAlbums list.Model
	albumMap   map[list.Item]*spotify.Album
}

// Called whenever the user has finished inputing a search query and selected the search type
// of the results to be displayed. This updates the state of result to match the desired
// specified content.
func (r Results) Refresh(query string, currentSelectedType string, s *auth.Session) Results {
	results, _ := spotify.Search(query, searchTypes, SEARCH_LIMIT, s)

	r.CurrentType = currentSelectedType

	switch r.CurrentType {
	case "track":
		r.trackMap = map[list.Item]*spotify.Track{}

		items := make([]list.Item, len(results.Tracks))
		for i, track := range results.Tracks {
			items[i] = components.UniqueItem{
				Name: components.Content(track.Name).AdjustFit(SEARCH_RESULTS_WIDTH - 5).String(),
				Id:   track.ID,
			}

			r.trackMap[items[i]] = track
		}

		r.listTracks = components.NewDefaultUniqueItemList(items, "Tracks")

	case "album":
		r.albumMap = map[list.Item]*spotify.Album{}

		items := make([]list.Item, len(results.Albums))
		for i, album := range results.Albums {
			items[i] = components.UniqueItem{
				Name: components.Content(album.Name).AdjustFit(SEARCH_RESULTS_WIDTH - 5).String(),
				Id:   album.ID,
			}

			r.albumMap[items[i]] = album
		}

		r.listAlbums = components.NewDefaultUniqueItemList(items, "Albums")
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

	// User's first time in the search view,
	// hasn't selected a type, so just hide
	// results.
	return ""
}

func (r Results) init() tea.Cmd {
	return nil
}

func (r Results) SelectedTrack() *spotify.Track {
	return r.trackMap[r.listTracks.SelectedItem()]
}

func (r Results) SelectedAlbum() *spotify.Album {
	return r.albumMap[r.listAlbums.SelectedItem()]
}
