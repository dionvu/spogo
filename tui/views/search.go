package views

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/err"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/spotify/auth"
	comp "github.com/dionvu/spogo/tui/views/components"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	MAX_RESULT_WIDTH      = 48
	LEFT_WIDTH            = 21
	TEXT_INPUT_CHAR_LIMIT = 156
	SEARCH_RESULT_LIMIT   = 42
	TOP_MARGIN_SEARCH     = 8
	MAX_RESULT_ITEM_WIDTH = MAX_RESULT_WIDTH - 5
	SEARCH_VIEW_WIDTH     = LEFT_WIDTH + MAX_RESULT_WIDTH

	TRACK    = "track"
	ALBUM    = "album"
	EPISODE  = "episode"
	PLAYLIST = "playlist"
)

var SEARCH_TYPES = []string{TRACK, ALBUM, PLAYLIST}

type Search struct {
	Input    SearchQuery
	TypeList SearchTypeList
	Results  Results
	Config   *config.Config

	session *auth.Session

	typeMap map[list.Item]string
}

func NewSearch(session *auth.Session, cfg *config.Config) Search {
	searchTypeListItemMap := map[list.Item]string{}
	searchTypeListItems := make([]list.Item, len(SEARCH_TYPES))

	for i, searchType := range SEARCH_TYPES {
		item := comp.ListItem(searchType)
		searchTypeListItems[i] = item
		searchTypeListItemMap[item] = searchType
	}

	return Search{
		session:  session,
		Config:   cfg,
		Input:    NewSearchQuery(),
		TypeList: NewSearchTypeList(searchTypeListItems),
		typeMap:  searchTypeListItemMap,
		Results:  Results{},
	}
}

func (r Results) SelectedTrack() *spotify.Track {
	return r.trackMap[r.listTracks.SelectedItem()]
}

func (r Results) SelectedAlbum() *spotify.Album {
	return r.albumMap[r.listAlbums.SelectedItem()]
}

func (r Results) SelectedPlaylist() *spotify.Playlist {
	return r.playlistMap[r.listPlaylists.SelectedItem()]
}

func (s Search) SelectedType() string {
	return s.typeMap[s.TypeList.Selected()]
}

// Renders the search view, this includes, the text area,
// the type selection, and the list of results.
func (s Search) View(term comp.Terminal, currentView string) string {
	queryAndTypeContainer := comp.NewDefaultTable()
	mainContainer := comp.NewDefaultTable()

	s.TypeList = s.TypeList.UpdateSelected(currentView)

	queryAndTypeContainer.AppendRows([]table.Row{
		{
			s.Input.Content().PadLinesLeft(2),
		},
		{
			s.TypeList.View(),
		},
	})

	queryAndType := comp.Content(queryAndTypeContainer.Render())

	details := func() comp.Content {
		switch s.Results.CurrentType {
		case TRACK:
			mins, secs := MsToMinutesAndSeconds(s.Results.SelectedTrack().DurationMs)
			return comp.Join(
				[]string{
					color.HiGreenString("Artist:    ") + s.Results.SelectedTrack().Artists[0].Name,
					color.HiGreenString("Duration:  ") + mins + "m:" + secs + "s",
				}, "\n\n")

		case ALBUM:
			return comp.Join(
				[]string{
					color.HiGreenString("Artist:  ") + s.Results.SelectedAlbum().Artists[0].Name,
					color.HiGreenString("Tracks:  ") + fmt.Sprint(s.Results.SelectedAlbum().TotalTracks),
				}, "\n\n")

		case PLAYLIST:
			return comp.Join(
				[]string{
					color.HiGreenString("Owner:   ") + s.Results.SelectedPlaylist().Owner.DisplayName,
					color.HiGreenString("Tracks:  ") + fmt.Sprint(s.Results.SelectedPlaylist().Tracks.Total),
				}, "\n\n")

		default:
			return "\n\n\n"
		}
	}()

	details = details.PadLinesLeft(28).AdjustFit(76)

	// if term.HeightIsSmall() || term.WidthIsSmall() {
	// 	mainContainer.AppendRow(table.Row{
	// 		// Offset to match playlist list's position.
	// 		queryAndType.PadLinesLeft(3),
	// 		s.Results.Content(),
	// 	})
	//
	// 	return comp.Join([]comp.Content{
	// 		comp.InvisibleBarV(TOP_MARGIN_SEARCH),
	// 		comp.Content(mainContainer.Render()),
	// 		"",
	// 		details.Append('\n', 1),
	// 		comp.InvisibleBar(SEARCH_VIEW_WIDTH).Append('\n', 1),
	// 	}).CenterVertical(term).CenterHorizontal(term).String()
	// }

	mainContainer.AppendRow(table.Row{
		queryAndType,
		s.Results.Content(),
	})

	c := comp.Join([]comp.Content{
		"\n" + comp.Content(mainContainer.Render()),
		"\n" + details.PadLinesLeft(4),
	}).String()

	return comp.Content(Box.String("[ Spogo 󰝚 ] "+ViewStatus{CurrentView: SEARCH_VIEW_RESULTS}.Content(s.Config).String(), comp.InvisibleBar(80).String()+"\n"+c+"\n")).CenterHorizontal(term).CenterVertical(term).String()
}

type Results struct {
	CurrentType string
	Items       spotify.SearchResult

	listTracks list.Model
	trackMap   map[list.Item]*spotify.Track

	listAlbums list.Model
	albumMap   map[list.Item]*spotify.Album

	listPlaylists list.Model
	playlistMap   map[list.Item]*spotify.Playlist
}

// Called whenever the user has finished inputing a search query and selected the search type
// of the results to be displayed. This updates the state of result to match the desired
// specified content.
func (r Results) Refresh(query string, selectedType string, s *auth.Session) Results {
	r.CurrentType = selectedType

	searchResults, err := spotify.Search(query, SEARCH_TYPES, SEARCH_RESULT_LIMIT, s)
	if err != nil {
		errors.Log(err)
	}

	switch r.CurrentType {
	case TRACK:
		listItems := make([]list.Item, len(searchResults.Tracks))
		r.trackMap = map[list.Item]*spotify.Track{}

		for i, track := range searchResults.Tracks {
			listItems[i] = comp.UniqueItem{
				Name: comp.Content(track.Name).AdjustFit(MAX_RESULT_ITEM_WIDTH).String(),
				Id:   track.ID,
			}
			r.trackMap[listItems[i]] = track
		}

		r.listTracks = comp.NewDefaultUniqueItemList(listItems, "Tracks")

	case ALBUM:
		listItems := make([]list.Item, len(searchResults.Albums))
		r.albumMap = map[list.Item]*spotify.Album{}

		for i, album := range searchResults.Albums {
			listItems[i] = comp.UniqueItem{
				Name: comp.Content(album.Name).AdjustFit(MAX_RESULT_ITEM_WIDTH).String(),
				Id:   album.ID,
			}
			r.albumMap[listItems[i]] = album
		}

		r.listAlbums = comp.NewDefaultUniqueItemList(listItems, "Albums")

	case PLAYLIST:
		listItems := make([]list.Item, len(searchResults.Playlists))
		r.playlistMap = map[list.Item]*spotify.Playlist{}

		for i, playlist := range searchResults.Playlists {
			listItems[i] = comp.UniqueItem{
				Name: comp.Content(playlist.Name).AdjustFit(MAX_RESULT_ITEM_WIDTH).String(),
				Id:   playlist.ID,
			}
			r.playlistMap[listItems[i]] = playlist
		}

		r.listPlaylists = comp.NewDefaultUniqueItemList(listItems, "Playlist")
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
		switch msg.String() {
		case "ctrl+c":
			return r, tea.Quit

		case "esc":
			return r, nil
		}
	}

	switch r.CurrentType {
	case TRACK:
		r.listTracks, cmd = r.listTracks.Update(msg)
	case ALBUM:
		r.listAlbums, cmd = r.listAlbums.Update(msg)
	case PLAYLIST:
		r.listPlaylists, cmd = r.listPlaylists.Update(msg)
	}

	return r, cmd
}

// Renders the result view based on the
// current selected display type. Ensure
// result.Refresh() was once before we display
// the result content.
func (r Results) view() string {
	switch r.CurrentType {
	case TRACK:
		return r.listTracks.View()
	case ALBUM:
		return r.listAlbums.View()
	case PLAYLIST:
		return r.listPlaylists.View()
	}

	// User's first time in the search view,
	// hasn't selected a type, so just hide
	// results.
	return ""
}

func (r Results) Content() comp.Content {
	return comp.Join([]comp.Content{
		comp.Content(r.view()),
		comp.InvisibleBar(MAX_RESULT_WIDTH),
	})
}

func (_ Results) init() tea.Cmd {
	return nil
}

type SearchTypeList struct {
	list     list.Model
	choice   string
	quitting bool
}

// The selected type as a list item.
func (stl SearchTypeList) Selected() list.Item {
	return stl.list.SelectedItem()
}

func NewSearchTypeList(items []list.Item) SearchTypeList {
	lm := SearchTypeList{
		list: comp.NewCustomList(items, "Select a search type:",
			comp.DEFAULT_WIDTH+4, comp.LIST_HEIGHT_SMALL-1),
	}

	return lm
}

// Highlights the title of the list if the user is currently making a selection.
func (stl SearchTypeList) UpdateSelected(currentView string) SearchTypeList {
	stl.list.Title = func() string {
		title := "Select a search type:"
		if currentView == SEARCH_VIEW_TYPE {
			return lg.NewStyle().Underline(true).Render(title)
		}
		return title
	}()

	return stl
}

func (m SearchTypeList) Update(msg tea.Msg) (SearchTypeList, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SearchTypeList) View() string {
	return m.list.View()
}

func (_ SearchTypeList) Init() tea.Cmd {
	return nil
}

type SearchQuery struct {
	Text textinput.Model
	err  error
}

func NewSearchQuery() SearchQuery {
	ti := textinput.New()
	ti.Placeholder = "What's on your mind?"
	ti.Focus()
	ti.CharLimit = TEXT_INPUT_CHAR_LIMIT
	ti.Width = LEFT_WIDTH - 1
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
		case tea.KeyCtrlC:
			return sq, tea.Quit

		case tea.KeyEsc:
			return sq, nil
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

func (sq SearchQuery) Content() comp.Content {
	return comp.Content(sq.View())
}

func (sq SearchQuery) Query() string {
	return sq.Text.Value()
}

func (sq SearchQuery) HideCursor() SearchQuery {
	sq.Text.Blur()
	return sq
}

type SearchResultView struct {
	items      *spotify.SearchResult
	query      string
	searchType string
}

func NewSearchResultView(searchQuery string, searchType string, s *auth.Session) *SearchResultView {
	srv := SearchResultView{
		query:      searchQuery,
		searchType: searchType,
	}

	if searchQuery == "" {
		log.Fatal("Empty search query")
	}

	fmt.Println(searchType)
	switch searchType {
	case "track":
		results, err := spotify.Search(searchQuery, []string{"track"}, SEARCH_RESULT_LIMIT, s)
		if err != nil {
			log.Fatal("Error getting results")
		}

		srv.items = results

		fmt.Println(srv.items.Tracks)
	default:
		log.Fatal("Unknown search type passed")
	}

	return &srv
}

// Converts the number of milliseconds into two string values
// of minutes and addittional seconds.
func MsToMinutesAndSeconds(ms int) (minutes string, seconds string) {
	m := ms / 60000
	s := (ms % 60000) / 1000

	minutes = fmt.Sprint(m)
	seconds = fmt.Sprint(s)

	if s < 10 {
		seconds = "0" + seconds
	}

	return minutes, seconds
}
