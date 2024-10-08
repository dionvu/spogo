package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
	"github.com/jedib0t/go-pretty/v6/table"
)

// TEMP
// Just a random image from a youtube channel.
const DEFAULT_PLAYLIST_IMAGE_URL = "https://yt3.googleusercontent.com/Rrn4HYgcjL1us1TVcr2MjePqJ6fPvKwCNZ6STs4-8oUyJg0Z86xJs1FFEnVQ8mshVeY31nIjxw=s160-c-k-c0x00ffffff-no-rj"

// The view that handles everything related to
// the user's playlist library.
type PlaylistView struct {
	// The content to be displayed when the view
	// is being displayed.
	Content Content

	// The playlists of the user.
	UserPlaylists *[]spotify.Playlist

	// The list selection for user to select
	// playlists through hovering and selecting.
	PlaylistList PlaylistList

	// The images of all the user's playlists.
	Images   []Image
	imageMap map[list.Item]*Image

	// The detailed information about the selected
	// playlist.
	PlaylistInfo *PlaylistInfo

	// Displays the alternative main views, with the
	// current view (player view) highlighted.
	ViewStatus *ViewStatus

	// Used to access playlists from a selected item or
	// selected playlist name.
	playlistsMap map[list.Item]*spotify.Playlist

	Session *auth.Session
}

// Creates the new playlist view by fetching the user's spotify playlists, determining
// their images to be displayed. Appending all playlists to a bubbletea list.
func NewPlaylistView(s *auth.Session, initialTerm Terminal) *PlaylistView {
	// Will holds our playlist names as "list items" to be displayed visually.
	items := []list.Item{}

	// Necessary to cache playlist images.
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "playlist")

	pv := &PlaylistView{
		Images:       []Image{},
		PlaylistInfo: &PlaylistInfo{},
		ViewStatus:   &ViewStatus{},
		playlistsMap: map[list.Item]*spotify.Playlist{},
		imageMap:     map[list.Item]*Image{},
		Session:      s,
	}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)

	for i, playlist := range *pv.UserPlaylists {
		pv.Images = append(pv.Images, Image{FilePath: path + fmt.Sprint(i) + ".jpeg"})

		if len(playlist.Images) != 0 {
			pv.Images[i].Update(playlist.Images[0].Url)
		} else {
			pv.Images[i].Update(DEFAULT_PLAYLIST_IMAGE_URL)
		}

		item := Item(playlist.Name)
		items = append(items, item)

		pv.imageMap[item] = &pv.Images[i]
		pv.playlistsMap[item] = &playlist
	}

	pv.PlaylistList = NewPlaylistListModel(items, PlaylistViewStyle.Title.Render("Playlists"), initialTerm)

	// Sets the initial choice, else the list will not move.
	if len(items) > 0 {
		pv.PlaylistList.choice = (*pv.UserPlaylists)[0].Name
	}

	return pv
}

// Updates the content to be displayed based on the dimensions of the terminal.
func (pv *PlaylistView) UpdateContent(term Terminal) {
	var ascii Content

	pv.PlaylistInfo.Update(pv.GetSelectedPlaylist())
	pv.ViewStatus.Update(PLAYLIST_VIEW)

	container := table.NewWriter()
	container.Style().Options.DrawBorder = false
	container.Style().Options.SeparateColumns = false

	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	pv.Content = func() Content {
		t.AppendRow(table.Row{
			Join([]Content{
				pv.PlaylistList.Content().Prepend('\n', 2).Append('\n', 1),
				pv.PlaylistInfo.Content(term).PadLinesLeft(2),
				Content("").Append(' ', 40),
			}, "\n"),
		})

		if term.IsSizeSmall() {
			pv.PlaylistList.list.SetHeight(LIST_HEIGHT_SMALL)

			ascii = pv.SelectedImage().AsciiSmallBW().Content().
				Prepend('\n', 1).Append('\n', 1).
				PadLinesLeft(2).CenterVertical(term)
		} else {
			pv.PlaylistList.list.SetHeight(LIST_HEIGHT_NORMAL)

			ascii = pv.SelectedImage().AsciiNormalBW().Content().
				Prepend('\n', 1).Append('\n', 1).
				PadLinesLeft(2).CenterVertical(term)
		}

		container.AppendRow(table.Row{
			Content(t.Render()).CenterVertical(term),
			ascii,
		})

		return Content(container.Render()).CenterHorizontal(term)
	}()
}

// Gets the playlist struct corresponding to the playlist that the user
// is hovering or selected.
func (pv *PlaylistView) GetSelectedPlaylist() *spotify.Playlist {
	return pv.playlistsMap[pv.PlaylistList.list.SelectedItem()]
}

// The selected playlist's image.
func (pv *PlaylistView) SelectedImage() *Image {
	return pv.imageMap[pv.PlaylistList.list.SelectedItem()]
}

// Updates the content and renders the view as a string.
func (pv *PlaylistView) View(playerView *PlayerView, term Terminal) string {
	pv.UpdateContent(term)
	return pv.Content.String()
}

// The detailed infomation about a playlist to
// be displayed when a playlist is selected
// or hovered.
type PlaylistInfo struct {
	Name        PlaylistName
	TotalTracks int
	Owner       string
}

func (pi *PlaylistInfo) Update(playlist *spotify.Playlist) {
	pi.Name = PlaylistName(playlist.Name)
	pi.TotalTracks = playlist.Tracks.Total
	pi.Owner = playlist.Owner.DisplayName
}

// Renders the playlistInfo as a content string.
func (pi PlaylistInfo) Content(term Terminal) Content {
	style := lg.NewStyle().Bold(true)
	// .Foreground(lg.Color("#458588"))

	return Join([]string{
		style.Render(pi.Name.AdjustFit(term).String()),
		"Tracks: " + fmt.Sprint(pi.TotalTracks),
		"Owner: " + fmt.Sprint(pi.Owner),
	}, "\n")
}

type PlaylistName string

// Adjusts the playlist string to fit
// within the terminal if it is too big.
func (pn PlaylistName) AdjustFit(term Terminal) PlaylistName {
	c := Content(pn)

	return PlaylistName(c.AdjustFit(term.Width))
}

func (pn PlaylistName) String() string {
	return string(pn)
}

// The list displaying the playlist names.
type PlaylistList struct {
	list     list.Model
	choice   string
	quitting bool
}

// A new playlist list model sized dynamically based on the initial terminal dimensions.
func NewPlaylistListModel(items []list.Item, title string, initialTerm Terminal) PlaylistList {
	w := DEFAULT_WIDTH
	h := LIST_HEIGHT_NORMAL

	if initialTerm.IsSizeSmall() {
		h = LIST_HEIGHT_SMALL
	}

	l := list.New(items, itemDelegate{}, w, h)
	l.SetFilteringEnabled(false)
	l.Title = title
	l.Styles.Title = lg.NewStyle().MarginLeft(0)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	lm := PlaylistList{list: l}

	return lm
}

func (pll PlaylistList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			pll.quitting = true
			return pll, tea.Quit

		case "esc":
			return pll, nil
		}
	}

	pll.list, cmd = pll.list.Update(msg)

	return pll, cmd
}

// The list as a content string.
func (pl PlaylistList) Content() Content {
	return Content(pl.list.View())
}

func (pl PlaylistList) View() string {
	return pl.list.View()
}

func (pl PlaylistList) Init() tea.Cmd {
	return nil
}
