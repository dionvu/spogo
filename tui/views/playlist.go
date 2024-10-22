package views

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dionvu/spogo/config"
	"github.com/dionvu/spogo/spotify"
	"github.com/dionvu/spogo/spotify/auth"
	comp "github.com/dionvu/spogo/tui/views/components"
	"github.com/jedib0t/go-pretty/v6/table"
)

// TEMP
// Just a random image from a youtube channel.
const DEFAULT_PLAYLIST_IMAGE_URL = "https://yt3.googleusercontent.com/Rrn4HYgcjL1us1TVcr2MjePqJ6fPvKwCNZ6STs4-8oUyJg0Z86xJs1FFEnVQ8mshVeY31nIjxw=s160-c-k-c0x00ffffff-no-rj"

// The view that handles everything related to
// the user's playlist library.
type Playlist struct {
	// The content to be displayed when the view
	// is being displayed.
	Content comp.Content

	// The playlists of the user.
	UserPlaylists *[]spotify.Playlist

	// The list selection for user to select
	// playlists through hovering and selecting.
	PlaylistList PlaylistList

	// The images of all the user's playlists.
	Images   []comp.Image
	imageMap map[list.Item]*comp.Image

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
func NewPlaylistView(s *auth.Session, initialTerm comp.Terminal) *Playlist {
	// Will holds our playlist names as "list items" to be displayed visually.
	items := []list.Item{}

	// Necessary to cache playlist images.
	cd, _ := os.UserCacheDir()
	path := filepath.Join(cd, config.APPNAME, "playlist")

	pv := &Playlist{
		Images:       []comp.Image{},
		PlaylistInfo: &PlaylistInfo{},
		ViewStatus:   &ViewStatus{},
		playlistsMap: map[list.Item]*spotify.Playlist{},
		imageMap:     map[list.Item]*comp.Image{},
		Session:      s,
	}

	pv.UserPlaylists, _ = spotify.UserPlaylists(pv.Session)
	if pv.UserPlaylists == nil {
		return pv
	}

	for i, playlist := range *pv.UserPlaylists {
		pv.Images = append(pv.Images, comp.Image{FilePath: path + fmt.Sprint(i) + ".jpeg"})

		if len(playlist.Images) != 0 {
			pv.Images[i].Update(playlist.Images[0].Url)
		} else {
			pv.Images[i].Update(DEFAULT_PLAYLIST_IMAGE_URL)
		}

		item := comp.ListItem(playlist.Name)
		items = append(items, item)

		pv.imageMap[item] = &pv.Images[i]
		pv.playlistsMap[item] = &playlist
	}

	pv.PlaylistList = PlaylistList{list: comp.NewDefaultList(items, "Playlists")}

	// Sets the initial choice, else the list will not move.
	if len(items) > 0 {
		pv.PlaylistList.choice = (*pv.UserPlaylists)[0].Name
	}

	return pv
}

// Updates the content to be displayed based on the dimensions of the terminal.
func (pv *Playlist) UpdateContent(term comp.Terminal) {
	var ascii comp.Content
	container := comp.NewDefaultTable()
	t := comp.NewDefaultTable()

	pv.PlaylistInfo.Update(pv.GetSelectedPlaylist())

	vs := ViewStatus{CurrentView: PLAYLIST_VIEW}

	pv.Content = func() comp.Content {
		t.AppendRow(table.Row{
			comp.Join([]comp.Content{
				pv.PlaylistList.Content().
					Prepend('\n', 1),

				comp.Content("").Append(' ', 35),
			}, "\n"),
		})

		ascii = pv.SelectedImage().AsciiSmall().Content().
			Append('\n', 1).
			PadLinesLeft(2)

		container.AppendRow(table.Row{
			comp.Content(t.Render()),
			ascii,
		})

		if term.WidthIsSmall() || term.HeightIsSmall() {
			return comp.Join([]comp.Content{
				"\n",
				comp.Content(container.Render()).Append('\n', 1).CenterHorizontal(term, -2),
				pv.PlaylistInfo.Content(term).CenterHorizontal(term),
				comp.Content("").Append('\n', 2),
			}, "\n").Prepend('\n', 5).CenterVertical(term)
		}

		return comp.Join([]comp.Content{
			"\n",
			comp.Content(container.Render()).Append('\n', 2).CenterHorizontal(term),
			pv.PlaylistInfo.Content(term).CenterHorizontal(term),
			"",
			vs.Content().CenterHorizontal(term),
		}, "\n").Prepend('\n', 5).CenterVertical(term)
	}()
}

// Gets the playlist struct corresponding to the playlist that the user
// is hovering or selected.
func (pv *Playlist) GetSelectedPlaylist() *spotify.Playlist {
	return pv.playlistsMap[pv.PlaylistList.list.SelectedItem()]
}

// The selected playlist's image.
func (pv *Playlist) SelectedImage() *comp.Image {
	return pv.imageMap[pv.PlaylistList.list.SelectedItem()]
}

// Updates the content and renders the view as a string.
func (pv *Playlist) View(playerView *Player, term comp.Terminal) string {
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
func (pi PlaylistInfo) Content(term comp.Terminal) comp.Content {
	return comp.Join(
		[]string{
			pi.Name.AdjustFit(term).String() + "\n",
			"Tracks: " + fmt.Sprint(pi.TotalTracks),
		}, "\n")
}

type PlaylistName string

// Adjusts the playlist string to fit
// within the terminal if it is too big.
func (pn PlaylistName) AdjustFit(term comp.Terminal) PlaylistName {
	c := comp.Content(pn)

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

func (pll PlaylistList) Update(msg tea.Msg) (PlaylistList, tea.Cmd) {
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
func (pl PlaylistList) Content() comp.Content {
	return comp.Content(pl.list.View())
}

func (pl PlaylistList) View() string {
	return pl.list.View()
}

func (pl PlaylistList) Init() tea.Cmd {
	return nil
}
